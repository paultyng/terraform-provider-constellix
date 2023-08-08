package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceConstellixHTTPCheck() *schema.Resource {
	return &schema.Resource{
		Read: datasourceConstellixHTTPCheckRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"protocol_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"check_sites": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"notification_groups": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"THIRTYSECONDS",
					"ONEMINUTE",
					"TWOMINUTES",
					"THREEMINUTES",
					"FOURMINUTES",
					"FIVEMINUTES",
					"TENMINUTES",
					"THIRTYMINUTES",
					"HALFDAY",
					"DAY",
				}, false),
			},
			"interval_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PARALLEL",
					"ONCEPERSITE",
					"ONCEPERREGION",
				}, false),
			},
			"verification_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SIMPLE",
					"MAJORITY",
				}, false),
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"search_string": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expected_status_code": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"notification_report_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func datasourceConstellixHTTPCheckRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)

	resp, err := client.GetbyId("https://api.sonar.constellix.com/rest/api/http")
	if err != nil {
		return err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodystring := string(bodybytes)

	var data []interface{}
	json.Unmarshal([]byte(bodystring), &data)

	var tp map[string]interface{}
	var flag bool
	for _, val := range data {
		tp = val.(map[string]interface{})
		if name == tp["name"].(string) {
			flag = true
			d.SetId(fmt.Sprintf("%v", tp["id"]))
			d.Set("name", tp["name"])
			d.Set("host", tp["host"])
			d.Set("protocol_type", tp["protocolType"])
			d.Set("ip_version", tp["ipVersion"])
			d.Set("port", tp["port"])
			d.Set("check_sites", tp["checkSites"])
			d.Set("notification_groups", tp["notificationGroups"])
			d.Set("interval", tp["interval"])
			d.Set("interval_policy", tp["monitorIntervalPolicy"])
			d.Set("verification_policy", tp["verificationPolicy"])
			d.Set("fqdn", tp["fqdn"])
			d.Set("path", tp["path"])
			d.Set("search_string", tp["searchString"])
			d.Set("expected_status_code", tp["expectedStatusCode"])
			d.Set("notification_report_timeout", tp["notificationReportTimeout"])
		}
	}
	if flag != true {
		return fmt.Errorf("HTTP check of specified name is not found")
	}
	return nil
}
