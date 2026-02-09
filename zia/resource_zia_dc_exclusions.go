package zia

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dc_exclusions"
)

func resourceDCExclusions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDCExclusionsCreate,
		ReadContext:   resourceDCExclusionsRead,
		UpdateContext: resourceDCExclusionsUpdate,
		DeleteContext: resourceDCExclusionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt64, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					if idInt64 < int64(math.MinInt) || idInt64 > int64(math.MaxInt) {
						return nil, fmt.Errorf("invalid id %q: out of range for int type", id)
					}
					_ = d.Set("datacenter_id", int(idInt64))
				} else {
					resp, err := dc_exclusions.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.DcID))
						_ = d.Set("datacenter_id", resp.DcID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Terraform state ID (string). Matches the API datacenter ID.",
			},
			"datacenter_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Datacenter ID (dcid) to exclude. Required on create; can be omitted when importing by ID. Use the numeric ID from zia_datacenters (e.g. datacenter_id) or this resource's id.",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Unix timestamp when the exclusion window starts. Either start_time or start_time_utc must be set.",
			},
			"start_time_utc": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Exclusion window start (UTC). Format: MM/DD/YYYY HH:MM am/pm. If set, overrides start_time.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Unix timestamp when the exclusion window ends. Either end_time or end_time_utc must be set.",
			},
			"end_time_utc": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Exclusion window end (UTC). Format: MM/DD/YYYY HH:MM am/pm. If set, overrides end_time.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
				Description:      "Description of the DC exclusion.",
			},
			"expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the exclusion has expired (read-only from API).",
			},
		},
	}
}

func resourceDCExclusionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req, err := expandDCExclusions(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating ZIA dc exclusions\n%+v\n", req)

	resp, _, err := dc_exclusions.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA dc exclusions. DcID: %v\n", resp.DcID)
	d.SetId(strconv.Itoa(resp.DcID))
	_ = d.Set("datacenter_id", resp.DcID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDCExclusionsRead(ctx, d, meta)
}

func resourceDCExclusionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	dcID, err := resourceDCExclusionsID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	all, err := dc_exclusions.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting DC exclusions: %w", err))
	}

	var found *dc_exclusions.DCExclusions
	for i := range all {
		if all[i].DcID == dcID {
			found = &all[i]
			break
		}
	}
	if found == nil {
		log.Printf("[WARN] Removing zia_dc_exclusions %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(found.DcID))
	_ = d.Set("datacenter_id", found.DcID)
	_ = d.Set("start_time", found.StartTime)
	_ = d.Set("end_time", found.EndTime)
	_ = d.Set("start_time_utc", FormatExclusionTimeUTC(found.StartTime))
	_ = d.Set("end_time_utc", FormatExclusionTimeUTC(found.EndTime))
	_ = d.Set("description", found.Description)
	_ = d.Set("expired", found.Expired)

	return nil
}

func resourceDCExclusionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	dcID, err := resourceDCExclusionsID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Updating zia dc exclusion datacenter_id: %v\n", dcID)

	all, err := dc_exclusions.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}
	var existing *dc_exclusions.DCExclusions
	for i := range all {
		if all[i].DcID == dcID {
			existing = &all[i]
			break
		}
	}
	if existing == nil {
		d.SetId("")
		return nil
	}

	req, err := expandDCExclusions(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, _, err := dc_exclusions.Update(ctx, service, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDCExclusionsRead(ctx, d, meta)
}

func resourceDCExclusionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	dcID, err := resourceDCExclusionsID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Deleting zia dc exclusion datacenter_id: %v\n", dcID)

	if _, err := dc_exclusions.Delete(ctx, service, dcID); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia_dc_exclusions %d deleted", dcID)

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

// resourceDCExclusionsID returns the datacenter ID for API calls: from datacenter_id attribute or from resource id (string).
// Uses the same technique as rule_labels: getIntFromResourceData for the int attribute, with fallback to parsing d.Id().
func resourceDCExclusionsID(d *schema.ResourceData) (int, error) {
	if id, ok := getIntFromResourceData(d, "datacenter_id"); ok {
		return id, nil
	}
	idStr := d.Id()
	if idStr == "" {
		return 0, fmt.Errorf("datacenter_id is required: set it (e.g. from zia_datacenters.datacenter_id) or import this resource by id")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: %w", idStr, err)
	}
	if id < int64(math.MinInt) || id > int64(math.MaxInt) {
		return 0, fmt.Errorf("invalid id %q: out of range for int type", idStr)
	}
	return int(id), nil
}

func expandDCExclusions(d *schema.ResourceData) (dc_exclusions.DCExclusions, error) {
	dcID, err := resourceDCExclusionsID(d)
	if err != nil {
		return dc_exclusions.DCExclusions{}, err
	}

	startTime, err := resolveExclusionTimeFromResourceData(d, "start_time", "start_time_utc")
	if err != nil {
		return dc_exclusions.DCExclusions{}, fmt.Errorf("start time: %w", err)
	}
	endTime, err := resolveExclusionTimeFromResourceData(d, "end_time", "end_time_utc")
	if err != nil {
		return dc_exclusions.DCExclusions{}, fmt.Errorf("end time: %w", err)
	}

	return dc_exclusions.DCExclusions{
		DcID:        dcID,
		StartTime:   startTime,
		EndTime:     endTime,
		Description: getString(d.Get("description")),
	}, nil
}

// resolveExclusionTimeFromResourceData uses the shared ResolveExclusionTimeEpoch util for consistency with zia_sub_cloud.
func resolveExclusionTimeFromResourceData(d *schema.ResourceData, epochKey, utcKey string) (int, error) {
	epochVal, hasEpoch := d.GetOk(epochKey)
	epoch := 0
	if hasEpoch && epochVal != nil {
		switch v := epochVal.(type) {
		case int:
			epoch = v
		case int64:
			epoch = int(v)
		}
	}
	utcStr := ""
	if v, ok := d.GetOk(utcKey); ok && v != nil {
		utcStr, _ = v.(string)
	}
	return ResolveExclusionTimeEpoch(hasEpoch, epoch, utcStr)
}
