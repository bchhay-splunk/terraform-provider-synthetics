// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package synthetics

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"time"

	sc2 "github.com/splunk/syntheticsclient/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceHttpCheckV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHttpCheckV2Create,
		ReadContext:   resourceHttpCheckV2Read,
		UpdateContext: resourceHttpCheckV2Update,
		DeleteContext: resourceHttpCheckV2Delete,

		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "http",
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"frequency": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"scheduling_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "round_robin",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^concurrent$|^round_robin$)`), "Setting must match concurrent or round_robin"),
						},
						"request_method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceHttpCheckV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkData := processHttpCheckV2Items(d)

	o, _, err := c.CreateHttpCheckV2(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Test.ID))

	resourceHttpCheckV2Read(ctx, d, meta)

	return diags
	// return nil
}

func resourceHttpCheckV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.GetHttpCheckV2(checkID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] GET HTTP BODY: ", o)

	return diags
}

func resourceHttpCheckV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkID := d.Id()

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeleteHttpCheckV2(checkIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Delete check response data: ", resp)
	d.SetId("")

	return diags
}

func resourceHttpCheckV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	checkID := d.Id()

	checkData := processHttpCheckV2Items(d)

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateHttpCheckV2(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] UPDATE BODY: ", o)

	d.Set("test.updated_at", time.Now().Format(time.RFC850))

	return resourceHttpCheckV2Read(ctx, d, meta)
}

func processHttpCheckV2Items(d *schema.ResourceData) sc2.HttpCheckV2Input {

	var check = buildHttpV2Data(d)
	log.Println("[DEBUG] HTTP V2 CHECK OUTPUT: ", check)
	return check
}