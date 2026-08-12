package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/pacfiles"
)

// pacVersionStatusToAction maps the declared pac_version_status to the
// version action the API expects. UNSTAGED is the declarative form of the
// UNSTAGE action: the version keeps existing but carries no status.
// REMOVE_LKG has no declarative equivalent here because it requires a
// replacement version to be named as last known good.
var pacVersionStatusToAction = map[string]string{
	"DEPLOYED": "DEPLOY",
	"STAGE":    "STAGE",
	"LKG":      "LKG",
	"UNSTAGED": "UNSTAGE",
}

// normalizePacVersionStatus maps the API's representation of "no status"
// (empty field) to the declarative UNSTAGED value.
func normalizePacVersionStatus(apiStatus string) string {
	if apiStatus == "" {
		return "UNSTAGED"
	}
	return apiStatus
}

// isPacResourceNotFound reports whether err means the PAC file or version
// does not exist. The PAC endpoints report a missing version as HTTP 400
// with code RESOURCE_NOT_FOUND rather than a plain 404.
func isPacResourceNotFound(err error) bool {
	respErr, ok := errorx.AsErrorResponse(err)
	if !ok {
		return false
	}
	if respErr.IsObjectNotFound() {
		return true
	}
	return respErr.Parsed != nil && fmt.Sprintf("%v", respErr.Parsed.Code) == "RESOURCE_NOT_FOUND"
}

func resourcePacFiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePacFilesCreate,
		ReadContext:   resourcePacFilesRead,
		UpdateContext: resourcePacFilesUpdate,
		DeleteContext: resourcePacFilesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				// Parse with the platform's int size so the conversion to
				// int below can never truncate; an oversized value fails
				// parsing and is treated as a name lookup instead.
				idInt, parseIDErr := strconv.ParseInt(id, 10, strconv.IntSize)
				var foundFile *pacfiles.PACFileConfig
				if parseIDErr == nil {
					allFiles, err := pacfiles.GetPacFiles(ctx, service, "")
					if err != nil {
						return nil, err
					}
					for _, pacFile := range allFiles {
						if pacFile.ID == int(idInt) {
							f := pacFile
							foundFile = &f
							break
						}
					}
					if foundFile == nil {
						return nil, fmt.Errorf("no PAC file found with ID: %d", idInt)
					}
				} else {
					resp, err := pacfiles.GetPacFileByName(ctx, service, id)
					if err != nil {
						return nil, err
					}
					foundFile = resp
				}

				// Adopt the deployed version as the tracked version.
				d.SetId(strconv.Itoa(foundFile.ID))
				_ = d.Set("pac_id", foundFile.ID)
				_ = d.Set("pac_version", foundFile.PACVersion)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the PAC file",
			},
			"pac_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the PAC file",
			},
			"pac_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The PAC file version managed by this resource. When omitted, the resource tracks the latest version it created (a new version is created whenever the content changes). Declare it explicitly to target a specific existing version with `pac_version_status`.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the PAC file",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the PAC file",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain of your organization to which the PAC file applies",
			},
			"pac_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content of the PAC file. The content is validated before the PAC file is created or a new version is saved.",
			},
			"pac_commit_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The commit message entered when saving the PAC file version",
			},
			"pac_url_obfuscated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the PAC file URL is obfuscated. When true, the obfuscated URL is returned in the `pac_sub_url` attribute.",
			},
			"pac_version_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEPLOYED",
				ValidateFunc: validation.StringInSlice([]string{"DEPLOYED", "STAGE", "LKG", "UNSTAGED"}, false),
				Description:  "The desired status of the PAC file version managed by this resource. The initial version is always deployed; `STAGE`, `LKG`, and `UNSTAGED` can be applied to subsequent versions. `UNSTAGED` removes the staged status from a staged version. Supported values: `DEPLOYED`, `STAGE`, `LKG`, `UNSTAGED`",
			},
			"delete_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The version to remove when a content change would exceed the maximum number of PAC file versions. Only used when saving a new version.",
			},
			"pac_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL location of the PAC file, auto-generated when the PAC file is first added",
			},
			"pac_sub_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The obfuscated URL of the PAC file. Returned when `pac_url_obfuscated` is true.",
			},
			"pac_verification_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The verification status of the PAC file, indicating whether any syntax errors were identified. Supported values: `VERIFY_NOERR`, `VERIFY_ERR`, `NOVERIFY`",
			},
		},
	}
}

// runPacContentValidation runs the PAC content through the service's
// validation endpoint and returns an error carrying the validation messages
// when the content is rejected.
func runPacContentValidation(ctx context.Context, zClient *Client, pacContent string) error {
	result, err := pacfiles.ValidatePacFile(ctx, zClient.Service, pacContent)
	if err != nil {
		return fmt.Errorf("PAC content validation request failed: %w", err)
	}
	if result.Success && result.ErrorCount == 0 {
		if result.WarningCount > 0 {
			log.Printf("[WARN] PAC content validated with %d warning(s): %+v", result.WarningCount, result.Messages)
		}
		return nil
	}
	var msgs []string
	for _, m := range result.Messages {
		msgs = append(msgs, fmt.Sprintf("line %d, column %d: %s", m.Line, m.Column, m.Message))
	}
	return fmt.Errorf("PAC content failed validation with %d error(s): %s", result.ErrorCount, strings.Join(msgs, "; "))
}

// applyPacVersionAction transitions the given PAC file version to the desired
// status. It first reads the version and skips the call entirely when the
// version already carries the desired status, so the transition is idempotent
// regardless of what status the save/clone operation assigned. The action
// endpoint stores the request body as the version's commit message, so the
// declared commit message is passed through unchanged.
func applyPacVersionAction(ctx context.Context, zClient *Client, pacID, pacVersion int, desiredStatus, commitMessage string) error {
	service := zClient.Service

	current, err := pacfiles.GetPacVersionID(ctx, service, pacID, pacVersion, "")
	if err != nil {
		return fmt.Errorf("failed to read PAC file %d version %d before status transition: %w", pacID, pacVersion, err)
	}
	currentStatus := normalizePacVersionStatus(current.PACVersionStatus)
	// When no commit message is declared, preserve the version's existing one.
	if commitMessage == "" {
		commitMessage = current.PACCommitMessage
	}
	if currentStatus == desiredStatus && commitMessage == current.PACCommitMessage {
		log.Printf("[INFO] PAC file %d version %d already has status %s; no transition needed", pacID, pacVersion, desiredStatus)
		return nil
	}

	action, ok := pacVersionStatusToAction[desiredStatus]
	if !ok {
		return fmt.Errorf("unsupported pac_version_status %q", desiredStatus)
	}
	// The service refuses to stage or unstage the currently deployed
	// version; staging always applies to a non-deployed version.
	if currentStatus == "DEPLOYED" && (desiredStatus == "STAGE" || desiredStatus == "UNSTAGED") {
		return fmt.Errorf("version %d is the currently deployed version and cannot be transitioned to %s. Change the PAC file content to create a new version and stage that instead, or declare pac_version to target another existing version", pacVersion, desiredStatus)
	}
	// UNSTAGED means "remove the current status": from STAGE that is the
	// UNSTAGE action, from LKG it is REMOVE_LKG.
	if currentStatus == "LKG" && desiredStatus == "UNSTAGED" {
		action = "REMOVE_LKG"
	}

	log.Printf("[INFO] Applying action %s to PAC file %d version %d", action, pacID, pacVersion)
	if _, err := pacfiles.UpdatePacFile(ctx, service, pacID, pacVersion, action, commitMessage, nil); err != nil {
		return fmt.Errorf("failed to apply action %s to PAC file %d version %d: %w", action, pacID, pacVersion, err)
	}
	return nil
}

func resourcePacFilesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	// The service always deploys the initial version; surface a clear error
	// instead of silently deploying content the user asked to stage.
	if status := d.Get("pac_version_status").(string); status != "DEPLOYED" {
		return diag.Errorf("pac_version_status must be DEPLOYED when creating a PAC file: the initial version is always deployed. %q can be applied when a subsequent version is created by changing the PAC file content.", status)
	}
	if raw := d.GetRawConfig(); !raw.IsNull() {
		if v := raw.GetAttr("pac_version"); !v.IsNull() {
			return diag.Errorf("pac_version cannot be declared when creating a PAC file: the service assigns version 1 to the initial version. Declare pac_version only to target an existing version with a status transition.")
		}
	}

	if err := runPacContentValidation(ctx, zClient, d.Get("pac_content").(string)); err != nil {
		return diag.FromErr(err)
	}

	req := expandPacFile(d)
	req.PACVersionStatus = "DEPLOYED"
	log.Printf("[INFO] Creating ZIA PAC file\n%+v\n", req.Name)

	resp, err := pacfiles.CreatePacFile(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA PAC file. ID: %d, version: %d", resp.ID, resp.PACVersion)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("pac_id", resp.ID)
	_ = d.Set("pac_version", resp.PACVersion)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourcePacFilesRead(ctx, d, meta)
}

func resourcePacFilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "pac_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no PAC file id is set"))
	}

	// Read the tracked version, not the deployed list: when the tracked
	// version is staged, the deployed list would return another version's
	// content and produce permanent drift.
	version, versionSet := getIntFromResourceData(d, "pac_version")

	var resp *pacfiles.PACFileConfig
	if versionSet {
		var err error
		resp, err = pacfiles.GetPacVersionID(ctx, service, id, version, "")
		if err != nil && !isPacResourceNotFound(err) {
			return diag.FromErr(err)
		}
	}

	// No tracked version, or the tracked version no longer exists: re-sync
	// to the file's deployed version, and remove the resource from state
	// only when the file itself is gone.
	if resp == nil {
		allFiles, err := pacfiles.GetPacFiles(ctx, service, "")
		if err != nil {
			return diag.FromErr(err)
		}
		for _, pacFile := range allFiles {
			if pacFile.ID == id {
				f := pacFile
				resp = &f
				break
			}
		}
		if resp == nil {
			log.Printf("[WARN] Removing PAC file %d from state because it no longer exists in ZIA", id)
			d.SetId("")
			return nil
		}
		if versionSet {
			log.Printf("[WARN] PAC file %d version %d no longer exists; re-syncing state to the deployed version %d", id, version, resp.PACVersion)
		}
	}

	d.SetId(strconv.Itoa(id))
	_ = d.Set("pac_id", id)
	_ = d.Set("pac_version", resp.PACVersion)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("domain", resp.Domain)
	_ = d.Set("pac_content", resp.PACContent)
	_ = d.Set("pac_commit_message", resp.PACCommitMessage)
	_ = d.Set("pac_url_obfuscated", resp.PACUrlObfuscated)
	_ = d.Set("pac_version_status", normalizePacVersionStatus(resp.PACVersionStatus))
	_ = d.Set("pac_url", resp.PACUrl)
	_ = d.Set("pac_sub_url", resp.PACSubURL)
	_ = d.Set("pac_verification_status", resp.PACVerificationStatus)

	return nil
}

func resourcePacFilesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "pac_id")
	if !ok {
		return diag.Errorf("no PAC file ID is set")
	}
	currentVersion, ok := getIntFromResourceData(d, "pac_version")
	if !ok {
		return diag.Errorf("no PAC file version is set")
	}
	desiredStatus := d.Get("pac_version_status").(string)

	// A version declared in the configuration pins the status action to that
	// specific version instead of the latest one the resource created.
	versionDeclared := false
	if raw := d.GetRawConfig(); !raw.IsNull() {
		if v := raw.GetAttr("pac_version"); !v.IsNull() {
			versionDeclared = true
		}
	}

	// Commit message changes are intentionally NOT a clone trigger: the
	// status action stamps the declared commit message on the version, so a
	// message change never requires a new version.
	contentChanged := d.HasChanges("pac_content", "name", "description", "domain", "pac_url_obfuscated")

	trackedVersion := currentVersion
	if versionDeclared {
		declaredVersion := d.Get("pac_version").(int)
		if contentChanged {
			// Re-pinning to a different version legitimately shows a
			// content diff against the previously tracked version. It is
			// acceptable exactly when the declared content matches the
			// pinned version's content; anything else needs a new version
			// and therefore must not be combined with a pin.
			pinnedMatches := false
			if d.HasChange("pac_content") && !d.HasChanges("name", "description", "domain", "pac_url_obfuscated") {
				pinned, err := pacfiles.GetPacVersionID(ctx, service, id, declaredVersion, "")
				if err != nil {
					if isPacResourceNotFound(err) {
						return diag.Errorf("pac_version %d does not exist for this PAC file. Declare a version that already exists to apply a status to it, or remove pac_version and change the PAC file content to create a new version.", declaredVersion)
					}
					return diag.FromErr(err)
				}
				pinnedMatches = pinned.PACContent == d.Get("pac_content").(string)
			}
			if !pinnedMatches {
				return diag.Errorf("pac_version %d is declared, but the configuration does not match that version. When declaring pac_version, pac_content must match the declared version's content exactly. To save changed content as a new version, remove pac_version from the configuration.", declaredVersion)
			}
			contentChanged = false
		}
		trackedVersion = declaredVersion
	}

	if contentChanged {
		// Content is immutable per version: save the change as a new
		// version cloned from the currently tracked version.
		if err := runPacContentValidation(ctx, zClient, d.Get("pac_content").(string)); err != nil {
			return diag.FromErr(err)
		}

		var deleteVersion *int
		if v, ok := d.GetOk("delete_version"); ok {
			dv := v.(int)
			deleteVersion = &dv
		}

		req := expandPacFile(d)
		log.Printf("[INFO] Saving new version of PAC file %d (cloned from version %d)", id, currentVersion)
		cloned, err := pacfiles.CreateClonedPacFileVersion(ctx, service, id, currentVersion, deleteVersion, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to save a new PAC file version (if the maximum number of versions was reached, set delete_version to remove one): %w", err))
		}
		trackedVersion = cloned.PACVersion
		log.Printf("[INFO] PAC file %d: new version %d created", id, trackedVersion)
		// The version exists from this point on; record it even if the
		// status transition below fails, so state stays accurate.
		_ = d.Set("pac_version", trackedVersion)
	}

	// Converge the tracked version to the desired status. The helper skips
	// the call when the version already carries the desired status (e.g.
	// when the save operation itself assigned it).
	if err := applyPacVersionAction(ctx, zClient, id, trackedVersion, desiredStatus, d.Get("pac_commit_message").(string)); err != nil {
		if versionDeclared && isPacResourceNotFound(err) {
			return diag.Errorf("pac_version %d does not exist for this PAC file. Declare a version that already exists to apply a status to it, or remove pac_version and change the PAC file content to create a new version.", trackedVersion)
		}
		return diag.FromErr(err)
	}
	// Record the targeted version only after the transition succeeded.
	_ = d.Set("pac_version", trackedVersion)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourcePacFilesRead(ctx, d, meta)
}

func resourcePacFilesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "pac_id")
	if !ok {
		return diag.Errorf("no PAC file ID is set")
	}
	log.Printf("[INFO] Deleting PAC file ID: %d", id)

	if _, err := pacfiles.DeletePacFile(ctx, service, id); err != nil {
		if respErr, ok := errorx.AsErrorResponse(err); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] PAC file %d already deleted", id)
		} else {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	log.Printf("[INFO] PAC file deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

// expandPacFile builds the request payload from the configuration. The
// version status is intentionally omitted: on create it is forced to
// DEPLOYED, and on new versions the status transition is performed through
// the version action endpoint after the version is saved.
//
// PACVerificationStatus is a mandatory field on save requests: the client
// must declare the verification state of the content it submits. Both
// callers (Create and the new-version path in Update) run the content
// through the validation endpoint first and abort on errors, so the
// provider can always assert VERIFY_NOERR here.
func expandPacFile(d *schema.ResourceData) *pacfiles.PACFileConfig {
	return &pacfiles.PACFileConfig{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Domain:                d.Get("domain").(string),
		PACContent:            d.Get("pac_content").(string),
		PACCommitMessage:      d.Get("pac_commit_message").(string),
		PACUrlObfuscated:      d.Get("pac_url_obfuscated").(bool),
		PACVerificationStatus: "VERIFY_NOERR",
	}
}
