package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/Constellix/constellix-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConstellixCNameRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceConstellixCNameRecordCreate,
		Update: resourceConstellixCNameRecordUpdate,
		Read:   resourceConstellixCNameRecordRead,
		Delete: resourceConstellixCNameRecordDelete,

		Importer: &schema.ResourceImporter{
			State: resourceConstellixCNameRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"source_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"geo_location": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
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
			"record_failover_values": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"check_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disable_flag": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
			},

			"record_failover_failover_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"record_failover_disable_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pools": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceConstellixCNameRecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	constellixClient := m.(*client.Client)
	params := strings.Split(d.Id(), ":")
	resp, err := constellixClient.GetbyId("v1/" + params[0] + "/" + params[1] + "/records/cname/" + params[2])
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil, err
		}
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		return nil, err
	}

	geoset := parseCGeoResponse(data["geolocation"].(map[string]interface{}))

	rcdf := data["recordFailover"]
	rcdfset := make(map[string]interface{})
	rcdflist := make([]interface{}, 0, 1)
	if rcdf != nil {
		rcdf1 := rcdf.(map[string]interface{})
		rcdfset["record_failover_failover_type"] = fmt.Sprintf("%v", rcdf1["failoverType"])
		rcdfset["record_failover_disable_flag"] = fmt.Sprintf("%v", rcdf1["disabled"])

		rcdfvalues := rcdf1["values"].([]interface{})
		log.Println("LIST VALUES: ", rcdfvalues)

		for _, valrcdf := range rcdfvalues {
			map1 := make(map[string]interface{})
			val1 := valrcdf.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", val1["value"])
			map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
			map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
			map1["check_id"] = fmt.Sprintf("%v", val1["checkId"])
			rcdflist = append(rcdflist, map1)
		}
	}

	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	d.Set("host", data["host"])
	d.Set("name", data["name"])
	d.Set("domain_id", params[1])
	d.Set("source_type", params[0])
	d.Set("ttl", data["ttl"])
	d.Set("geo_location", geoset)
	d.Set("record_option", data["recordOption"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("note", data["note"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("pools", data["pools"])
	d.Set("contact_ids", data["contactId"])
	d.Set("record_failover_values", rcdflist)
	d.Set("record_failover_failover_type", rcdfset["record_failover_failover_type"])
	d.Set("record_failover_disable_flag", rcdfset["record_failover_disable_flag"])
	log.Printf("[DEBUG] %s finished import", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceConstellixCNameRecordCreate(d *schema.ResourceData, m interface{}) error {

	constellixConnect := m.(*client.Client)

	aAttr := models.CRecordAttributes{}

	if name, ok := d.GetOk("name"); ok {
		aAttr.Name = name.(string)
	}
	if host, ok := d.GetOk("host"); ok {
		aAttr.Host = host.(string)
	}
	if TTL, ok := d.GetOk("ttl"); ok {
		aAttr.TTL = TTL.(int)
	}
	if RecordOption, ok := d.GetOk("record_option"); ok {
		aAttr.RecordOption = RecordOption.(string)
	}
	if NoAnswer, ok := d.GetOk("noanswer"); ok {
		aAttr.NoAnswer = NoAnswer.(bool)
	}
	if Note, ok := d.GetOk("note"); ok {
		aAttr.Note = Note.(string)
	}
	if GtdRegion, ok := d.GetOk("gtd_region"); ok {
		aAttr.GtdRegion = GtdRegion.(int)
	}
	if Type, ok := d.GetOk("type"); ok {
		aAttr.Type = Type.(string)
	}
	if contactid, ok := d.GetOk("contact_ids"); ok {
		aAttr.ContactId = toListOfInt(contactid)
	}
	if pools, ok := d.GetOk("pools"); ok {
		aAttr.Pools = toListOfInt(pools)
	}

	aAttr.GeoLocation = buildCGeoPayload(d)

	var valuesrcdf *models.ValuesRCDFCrecord
	var rcdfa *models.RCDFACRecord //added
	valueslist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("record_failover_values"); ok {
		rcdfa = &models.RCDFACRecord{} //added
		valuesrcdf = &models.ValuesRCDFCrecord{}
		tp := value.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["checkId"] = fmt.Sprintf("%v", inner["check_id"])
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["sortOrder"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["sort_order"]))
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			valueslist = append(valueslist, map1)
		}
		rcdfa.Values = valueslist
	}

	if failovertype, ok := d.GetOk("record_failover_failover_type"); ok {
		rcdfa.FailoverTypeRCDFA, _ = strconv.Atoi(fmt.Sprintf("%v", failovertype)) //added
	}

	if disableflag, ok := d.GetOk("record_failover_disable_flag"); ok {
		rcdfa.DisableFlagRCDFA, _ = strconv.ParseBool(fmt.Sprintf("%v", disableflag)) //added
	}

	if valuesrcdf != nil {
		rcdfa.Values = valueslist     //added
		aAttr.RecordFailoverA = rcdfa //added
	} else {
		aAttr.RecordFailoverA = nil
	}
	resp, err := constellixConnect.Save(aAttr, "v1/"+d.Get("source_type").(string)+"/"+d.Get("domain_id").(string)+"/records/cname")

	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	json.Unmarshal([]byte(bodyString[1:len(bodyString)-1]), &data)

	d.SetId(fmt.Sprintf("%.0f", data["id"]))

	return resourceConstellixCNameRecordRead(d, m)
}

func resourceConstellixCNameRecordRead(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	id := d.Id()

	resp, err := constellixClient.GetbyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/cname/" + id)
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		return err
	}

	geoset := parseCGeoResponse(data["geolocation"].(map[string]interface{}))

	rcdf := data["recordFailover"]
	rcdfset := make(map[string]interface{})
	rcdflist := make([]interface{}, 0, 1)
	if rcdf != nil {
		rcdf1 := rcdf.(map[string]interface{})
		rcdfset["record_failover_failover_type"] = fmt.Sprintf("%v", rcdf1["failoverType"])
		rcdfset["record_failover_disable_flag"] = fmt.Sprintf("%v", rcdf1["disabled"])

		rcdfvalues := rcdf1["values"].([]interface{})
		log.Println("LIST VALUES: ", rcdfvalues)

		for _, valrcdf := range rcdfvalues {
			map1 := make(map[string]interface{})
			val1 := valrcdf.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", val1["value"])
			map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
			map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
			map1["check_id"] = fmt.Sprintf("%v", val1["checkId"])
			rcdflist = append(rcdflist, map1)
		}
	}

	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	d.Set("host", data["host"])
	d.Set("name", data["name"])
	d.Set("ttl", data["ttl"])
	d.Set("geo_location", geoset)
	d.Set("record_option", data["recordOption"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("note", data["note"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("pools", data["pools"])
	d.Set("contact_ids", data["contactId"])
	d.Set("record_failover_values", rcdflist)
	d.Set("record_failover_failover_type", rcdfset["record_failover_failover_type"])
	d.Set("record_failover_disable_flag", rcdfset["record_failover_disable_flag"])
	return nil
}

func resourceConstellixCNameRecordUpdate(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	aAttr := models.CRecordAttributes{}

	if ttl, ok := d.GetOk("ttl"); ok {
		aAttr.TTL = ttl.(int)
	}
	if name, ok := d.GetOk("name"); ok {
		aAttr.Name = name.(string)
	}
	if host, ok := d.GetOk("host"); ok {
		aAttr.Host = host.(string)
	}
	if recordoption, ok := d.GetOk("record_option"); ok {
		aAttr.RecordOption = recordoption.(string)
	}

	if _, ok := d.GetOk("noanswer"); ok {
		aAttr.NoAnswer = d.Get("noanswer").(bool)
	}

	if note, ok := d.GetOk("note"); ok {
		aAttr.Note = note.(string)
	}

	if _, ok := d.GetOk("gtd_region"); ok {
		aAttr.GtdRegion = d.Get("gtd_region").(int)
	}

	if _, ok := d.GetOk("type"); ok {
		aAttr.Type = d.Get("type").(string)
	}
	if contactid, ok := d.GetOk("contact_ids"); ok {
		aAttr.ContactId = toListOfInt(contactid)
	}
	if pools, ok := d.GetOk("pools"); ok {
		aAttr.Pools = toListOfInt(pools)
	}

	aAttr.GeoLocation = buildCGeoPayload(d)

	var valuesrcdf *models.ValuesRCDFCrecord
	var rcdfa *models.RCDFACRecord //added
	valueslist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("record_failover_values"); ok {
		rcdfa = &models.RCDFACRecord{} //added
		valuesrcdf = &models.ValuesRCDFCrecord{}
		tp := value.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["checkId"] = fmt.Sprintf("%v", inner["check_id"])
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["sortOrder"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["sort_order"]))
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			valueslist = append(valueslist, map1)
		}
		rcdfa.Values = valueslist
	}

	if failovertype, ok := d.GetOk("record_failover_failover_type"); ok {
		rcdfa.FailoverTypeRCDFA, _ = strconv.Atoi(fmt.Sprintf("%v", failovertype)) //added
	}

	if disableflag, ok := d.GetOk("record_failover_disable_flag"); ok {
		rcdfa.DisableFlagRCDFA, _ = strconv.ParseBool(fmt.Sprintf("%v", disableflag)) //added
	}

	if valuesrcdf != nil {
		rcdfa.Values = valueslist     //added
		aAttr.RecordFailoverA = rcdfa //added
	} else {
		aAttr.RecordFailoverA = nil
	}

	arecordid := d.Id()

	_, err := constellixClient.UpdatebyID(aAttr, "v1/"+d.Get("source_type").(string)+"/"+d.Get("domain_id").(string)+"/records/cname/"+arecordid)
	if err != nil {
		return err
	}
	return resourceConstellixCNameRecordRead(d, m)

}

func resourceConstellixCNameRecordDelete(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	arecordid := d.Id()

	err := constellixClient.DeletebyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/cname/" + arecordid)
	if err != nil {
		return err
	}
	d.SetId("")
	return err
}

// parseCGeoResponse reads the description of a geolocation returned
// by the API as part of a CNAME record resource, and returns
// a map with string values that conforms to the schema
func parseCGeoResponse(g map[string]interface{}) []interface{} {
	log.Println("GEOLOC VALUE: ", g)
	geo := make([]interface{}, 0, 1)
	if g != nil {
		m := make(map[string]interface{})
		if v, ok := g["geoipFilter"]; ok {
			m["geo_ip_user_region"], _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}
		if v, ok := g["drop"]; ok {
			m["drop"] = fmt.Sprintf("%v", v)
		}
		if v, ok := g["geoipProximity"]; ok {
			m["geo_ip_proximity"], _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}
		if v, ok := g["geoipFailover"]; ok {
			m["geo_ip_failover"] = fmt.Sprintf("%v", v)
		}
		geo = append(geo, m)
	}
	return geo
}

func buildCGeoPayload(d *schema.ResourceData) *models.GeolocationCrecord {
	var geoloc models.GeolocationCrecord
	userRegion := make([]int, 0, 1)
	if v, ok := d.GetOk("geo_location"); ok {
		vs := v.(*schema.Set).List()

		for _, vs := range vs {

			inner := vs.(map[string]interface{})
			if val, ok := inner["geo_ip_user_region"]; ok {
				i, _ := strconv.Atoi(fmt.Sprintf("%v", val))
				if i != 0 {
					userRegion = append(userRegion, i)
				}
			}
			if val, ok := inner["drop"]; ok {
				geoloc.Drop, _ = strconv.ParseBool(fmt.Sprintf("%v", val))
			}
			if val, ok := inner["geo_ip_proximity"]; ok {
				geoloc.GeoIpProximity, _ = strconv.Atoi(fmt.Sprintf("%v", val))
			}
			if val, ok := inner["geo_ip_failover"]; ok {
				geoloc.GeoIpFailOver, _ = strconv.ParseBool(fmt.Sprintf("%v", val))
			}
		}
	}
	if len(userRegion) > 0 {
		geoloc.GeoIpUserRegion = userRegion
	}
	return &geoloc
}
