package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceConstellixAnamerecord() *schema.Resource {
	return &schema.Resource{
		Read: datasourceConstellixAnamerecordRead,
		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"geo_location": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geo_ip_user_region": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"drop": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"geo_ip_proximity": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"geo_ip_failover": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},

			"record_option": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"noanswer": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"skip_lookup": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gtd_region": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"contact_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"roundrobin": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disable_flag": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},

			"record_failover_values": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disable_flag": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"record_failover_failover_type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"record_failover_disable_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func datasourceConstellixAnamerecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)

	resp, err := client.GetbyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/aname")
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)

	var data []interface{}
	json.Unmarshal([]byte(bodyString), &data)

	var tp map[string]interface{}
	var flag bool
	for _, val := range data {
		tp = val.(map[string]interface{})
		if name == tp["name"].(string) {
			flag = true

			geoloc1 := tp["geolocation"]

			geoLocMap := make(map[string]interface{})
			if geoloc1 != nil {
				geoloc := geoloc1.(map[string]interface{})
				if geoloc["geoipFilter"] != nil {
					geoLocMap["geo_ip_user_region"] = fmt.Sprintf("%v", geoloc["geoipFilter"])
				}
				if geoloc["drop"] != nil {
					geoLocMap["drop"] = fmt.Sprintf("%v", geoloc["drop"])
				}
				if geoloc["geoipFailover"] != nil {
					geoLocMap["geo_ip_failover"] = fmt.Sprintf("%v", geoloc["geoipFailover"])
				}
				if geoloc["geoipProximity"] != nil {
					geoLocMap["geo_ip_proximity"] = fmt.Sprintf("%v", geoloc["geoipProximity"])
				}
				d.Set("geo_location", geoLocMap)
			} else {
				d.Set("geo_location", geoLocMap)
			}

			rrlist := make([]interface{}, 0, 1)
			if tp["roundRobin"] != nil {
				recroundrobin := tp["roundRobin"].([]interface{})
				rrlist := make([]interface{}, 0, 1)
				for _, valrrf := range recroundrobin {
					map1 := make(map[string]interface{})
					val1 := valrrf.(map[string]interface{})
					map1["value"] = fmt.Sprintf("%v", val1["value"])
					map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
					rrlist = append(rrlist, map1)
				}
			}

			rcdf := tp["recordFailover"]
			rcdflist := make([]interface{}, 0, 1)
			if rcdf != nil {
				rcdf1 := rcdf.(map[string]interface{})

				d.Set("record_failover_failover_type", fmt.Sprintf("%v", rcdf1["failoverType"]))
				d.Set("record_failover_disable_flag", fmt.Sprintf("%v", rcdf1["disabled"]))

				rcdfvalues := rcdf1["values"].([]interface{})

				for _, valrcdf := range rcdfvalues {
					map1 := make(map[string]interface{})
					val1 := valrcdf.(map[string]interface{})
					map1["value"] = fmt.Sprintf("%v", val1["value"])
					map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
					map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
					rcdflist = append(rcdflist, map1)
				}
			}

			d.SetId(fmt.Sprintf("%v", tp["id"]))
			d.Set("name", tp["name"])
			d.Set("ttl", tp["ttl"])
			d.Set("record_option", tp["recordOption"])
			d.Set("noanswer", tp["noAnswer"])
			d.Set("skip_lookup", tp["skipLookup"])
			d.Set("note", tp["note"])
			d.Set("gtd_region", tp["gtdRegion"])
			d.Set("type", tp["type"])
			d.Set("contact_ids", tp["contactids"])
			d.Set("roundrobin", rrlist)
			d.Set("record_failover_values", rcdflist)
		}
	}
	if !flag {
		return fmt.Errorf("ANAME record of specified name is not available")
	}
	return nil
}
