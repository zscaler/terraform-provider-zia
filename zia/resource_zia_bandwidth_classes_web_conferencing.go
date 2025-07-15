package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
)

func resourceBandwdithClassesWebConferencing() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceBandwdithClassesWebConferencingsRead,
		CreateContext: resourceBandwdithClassesWebConferencingCreate,
		UpdateContext: resourceBandwdithClassesWebConferencingUpdate,
		DeleteContext: resourceFuncNoOp,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			classType := d.Get("type").(string)
			apps := d.Get("applications").(*schema.Set).List()

			var validApps map[string]struct{}

			switch classType {
			case "BANDWIDTH_CAT_WEBCONF":
				validApps = map[string]struct{}{
					"WEBEX":       {},
					"GOTOMEETING": {},
					"LIVEMEETING": {},
					"INTERCALL":   {},
					"CONNECT":     {},
				}
			case "BANDWIDTH_CAT_VOIP":
				validApps = map[string]struct{}{
					"SKYPE": {},
					"":      {},
				}
			default:
				// If type isn't one of those two, no validation needed
				return nil
			}

			for _, a := range apps {
				app := a.(string)
				if _, ok := validApps[app]; !ok {
					return fmt.Errorf("application %q is not valid for type %q", app, classType)
				}
			}

			return nil
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("class_id", idInt)
				} else {
					resp, err := bandwidth_classes.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("class_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"class_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bandwidth classname.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bandwidth classname.",
				ValidateFunc: validation.StringInSlice([]string{
					"BANDWIDTH_CAT_WEBCONF",
					"BANDWIDTH_CAT_VOIP",
				}, false),
			},
			"applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBandwdithClassesWebConferencingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	name := d.Get("name").(string)
	if name == "" {
		return diag.Errorf("'name' must be provided to locate the existing bandwidth class")
	}

	existing, err := bandwidth_classes.GetByName(ctx, service, name)
	if err != nil {
		return diag.Errorf("failed to find existing bandwidth class by name %q: %v", name, err)
	}

	log.Printf("[INFO] Found existing bandwidth class %q with ID %d", existing.Name, existing.ID)
	d.SetId(strconv.Itoa(existing.ID))
	_ = d.Set("class_id", existing.ID)

	req := expandBandwidthClassesWebConferencing(d)
	req.ID = existing.ID

	log.Printf("[INFO] Updating bandwidth class %q via PUT:\n%+v\n", name, req)

	if _, _, err := bandwidth_classes.Update(ctx, service, req.ID, &req); err != nil {
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

	return resourceBandwdithClassesWebConferencingsRead(ctx, d, meta)
}

func resourceBandwdithClassesWebConferencingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "class_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no bandwidth classs id is set"))
	}
	resp, err := bandwidth_classes.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia bandwidth classs %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia bandwidth classs:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("class_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("type", resp.Type)
	_ = d.Set("applications", resp.Applications)

	return nil
}

func resourceBandwdithClassesWebConferencingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var id int
	var ok bool

	id, ok = getIntFromResourceData(d, "class_id")
	if !ok || id == 0 {
		name := d.Get("name").(string)
		if name == "" {
			return diag.Errorf("either 'class_id' or 'name' must be set for update")
		}

		existing, err := bandwidth_classes.GetByName(ctx, service, name)
		if err != nil {
			return diag.Errorf("failed to find bandwidth class with name %q: %v", name, err)
		}

		id = existing.ID
		d.SetId(strconv.Itoa(id))
		_ = d.Set("class_id", id)
		log.Printf("[INFO] Retrieved class ID %d for update by name %q", id, name)
	}

	log.Printf("[INFO] Updating ZIA bandwidth class ID: %d\n", id)

	req := expandBandwidthClassesWebConferencing(d)
	req.ID = id

	if _, err := bandwidth_classes.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, _, err := bandwidth_classes.Update(ctx, service, id, &req); err != nil {
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

	return resourceBandwdithClassesWebConferencingsRead(ctx, d, meta)
}

func expandBandwidthClassesWebConferencing(d *schema.ResourceData) bandwidth_classes.BandwidthClasses {
	return bandwidth_classes.BandwidthClasses{
		Name:         d.Get("name").(string),
		Type:         d.Get("type").(string),
		Applications: SetToStringList(d, "applications"),
	}
}
