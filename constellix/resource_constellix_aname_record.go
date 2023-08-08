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

func resourceConstellixANAMERecord() *schema.Resource {
	return &schema.Resource{
		Create:        resourceConstellixANAMERecordCreate,
		Read:          resourceConstellixANAMERecordRead,
		Update:        resourceConstellixANAMERecordUpdate,
		Delete:        resourceConstellixANAMERecordDelete,
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: resourceConstellixANAMERecordImport,
		},
		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
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
							Required: true,
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
				Type:     schema.TypeInt,
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

func resourceConstellixANAMERecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	constellixClient := m.(*client.Client)
	params := strings.Split(d.Id(), ":")
	resp, err := constellixClient.GetbyId("v1/" + params[0] + "/" + params[1] + "/records/aname/" + params[2])
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

	geoloc1 := data["geolocation"]
	log.Println("GEOLOC VALUE INSIDE READ :", geoloc1)

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

	arecroundrobin := data["roundRobin"].([]interface{})
	rrlist := make([]interface{}, 0, 1)
	for _, valrrf := range arecroundrobin {
		map1 := make(map[string]interface{})
		val1 := valrrf.(map[string]interface{})
		map1["value"] = fmt.Sprintf("%v", val1["value"])
		map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
		rrlist = append(rrlist, map1)
	}
	log.Printf("tttttt %v", rrlist)

	rcdf := data["recordFailover"]
	rcdfSet := make(map[string]interface{})
	rcdflist := make([]interface{}, 0, 1)
	if rcdf != nil {
		rcdf1 := rcdf.(map[string]interface{})
		rcdfSet["record_failover_failover_type"] = fmt.Sprintf("%v", rcdf1["failoverType"])
		rcdfSet["record_failover_disable_flag"] = fmt.Sprintf("%v", rcdf1["disabled"])

		rcdfValues := rcdf1["values"].([]interface{})

		for _, valrcdf := range rcdfValues {
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
	d.Set("name", data["name"])
	d.Set("ttl", data["ttl"])
	d.Set("domain_id", params[1])
	d.Set("source_type", params[0])
	d.Set("record_option", data["recordOption"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("skip_lookup", data["skipLookup"])
	d.Set("note", data["note"])
	d.Set("pools", data["pools"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("contact_ids", data["contactIds"])
	d.Set("roundrobin", rrlist)
	d.Set("record_failover_values", rcdflist)
	d.Set("record_failover_failover_type", rcdfSet["record_failover_failover_type"])
	d.Set("record_failover_disable_flag", rcdfSet["record_failover_disable_flag"])
	log.Printf("[DEBUG] %s finished import", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceConstellixANAMERecordCreate(d *schema.ResourceData, m interface{}) error {
	constellixConnect := m.(*client.Client)

	anameAttr := models.AnameAttributes{}

	if name, ok := d.GetOk("name"); ok {
		anameAttr.Name = name.(string)
	}

	if value, ok := d.GetOk("ttl"); ok {
		anameAttr.TTL = value.(int)
	}

	if recordoption, ok := d.GetOk("record_option"); ok {
		anameAttr.RecordOption = recordoption.(string)
	}

	if noanswer, ok := d.GetOk("noanswer"); ok {
		anameAttr.NoAnswer = noanswer.(bool)
	}

	if skip_lookup, ok := d.GetOk("skip_lookup"); ok {
		anameAttr.SkipLookup = skip_lookup.(bool)
	}

	if note, ok := d.GetOk("note"); ok {
		anameAttr.Note = note.(string)
	}

	if gtdregion, ok := d.GetOk("gtd_region"); ok {
		anameAttr.GtdRegion = gtdregion.(int)
	}

	if types, ok := d.GetOk("type"); ok {
		anameAttr.Type = types.(string)
	}
	if pools, ok := d.GetOk("pools"); ok {
		anameAttr.Pools = toListOfInt(pools)
	}

	if contactid, ok := d.GetOk("contact_ids"); ok {
		anameAttr.ContactIDs = toListOfInt(contactid)
	}

	geoloc := &models.GeolocationANAME{}
	if geoipuserregion, ok := d.GetOk("geo_location"); ok {
		geouserlist := make([]int, 0, 1)
		tp := geoipuserregion.(map[string]interface{})
		if tp["geo_ip_user_region"] != nil {
			var1, _ := strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_user_region"]))
			geouserlist = append(geouserlist, var1)
			geoloc.GeoIpUserRegion = geouserlist
		}
		if tp["drop"] != nil {
			geoloc.Drop, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["drop"]))
		}
		if tp["geo_ip_failover"] != nil {
			geoloc.GeoIpFailOver, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["geo_ip_failover"]))
		}
		if tp["geo_ip_proximity"] != nil {
			geoloc.GeoIpProximity, _ = strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_proximity"]))
		}
	}
	anameAttr.GeoLocation = geoloc

	maplistrr := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roundrobin"); ok {
		tp := val.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			maplistrr = append(maplistrr, map1)
		}
		anameAttr.RoundRobin = maplistrr
	}

	valueslist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("record_failover_values"); ok {
		rcdfaname := &models.RCDFAname{}
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
		rcdfaname.Values = valueslist

		if failtype, ok := d.GetOk("record_failover_failover_type"); ok {
			rcdfaname.FailoverTypeAname, _ = strconv.Atoi(fmt.Sprintf("%v", failtype))
		}

		if disflag, ok := d.GetOk("record_failover_disable_flag"); ok {
			rcdfaname.DisableFlagAname, _ = strconv.ParseBool(fmt.Sprintf("%v", disflag))
		}

		rcdfaname.Values = valueslist
		anameAttr.RecordFailoverAname = rcdfaname
	}

	resp, err := constellixConnect.Save(anameAttr, "v1/"+d.Get("source_type").(string)+"/"+d.Get("domain_id").(string)+"/records/aname")

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

	return resourceConstellixANAMERecordRead(d, m)
}

func resourceConstellixANAMERecordRead(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	anameid := d.Id()

	resp, err := constellixClient.GetbyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/aname/" + anameid)
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

	geoloc1 := data["geolocation"]
	log.Println("GEOLOC VALUE INSIDE READ :", geoloc1)

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

	arecroundrobin := data["roundRobin"].([]interface{})
	rrlist := make([]interface{}, 0, 1)
	for _, valrrf := range arecroundrobin {
		map1 := make(map[string]interface{})
		val1 := valrrf.(map[string]interface{})
		map1["value"] = fmt.Sprintf("%v", val1["value"])
		map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
		rrlist = append(rrlist, map1)
	}

	rcdf := data["recordFailover"]
	rcdfSet := make(map[string]interface{})
	rcdflist := make([]interface{}, 0, 1)
	if rcdf != nil {
		rcdf1 := rcdf.(map[string]interface{})
		rcdfSet["record_failover_failover_type"] = fmt.Sprintf("%v", rcdf1["failoverType"])
		rcdfSet["record_failover_disable_flag"] = fmt.Sprintf("%v", rcdf1["disabled"])

		rcdfValues := rcdf1["values"].([]interface{})

		for _, valrcdf := range rcdfValues {
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
	d.Set("name", data["name"])
	d.Set("ttl", data["ttl"])
	d.Set("record_option", data["recordOption"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("skip_lookup", data["skipLookup"])
	d.Set("note", data["note"])
	d.Set("pools", data["pools"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("contact_ids", data["contactIds"])
	d.Set("roundrobin", rrlist)
	d.Set("record_failover_values", rcdflist)
	d.Set("record_failover_failover_type", rcdfSet["record_failover_failover_type"])
	d.Set("record_failover_disable_flag", rcdfSet["record_failover_disable_flag"])

	return nil

}

func resourceConstellixANAMERecordUpdate(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)

	anameAttr := models.AnameAttributes{}

	if value, ok := d.GetOk("ttl"); ok {
		anameAttr.TTL = value.(int)
	}

	if _, ok := d.GetOk("record_option"); ok {
		anameAttr.RecordOption = d.Get("record_option").(string)
	}

	if name, ok := d.GetOk("name"); ok {
		anameAttr.Name = name.(string)
	}

	if _, ok := d.GetOk("noanswer"); ok {
		anameAttr.NoAnswer = d.Get("noanswer").(bool)
	}

	if _, ok := d.GetOk("skip_lookup"); ok {
		anameAttr.SkipLookup = d.Get("skip_lookup").(bool)
	}

	if _, ok := d.GetOk("note"); ok {
		anameAttr.Note = d.Get("note").(string)
	}

	if _, ok := d.GetOk("gtd_region"); ok {
		anameAttr.GtdRegion = d.Get("gtd_region").(int)
	}

	if _, ok := d.GetOk("type"); ok {
		anameAttr.Type = d.Get("type").(string)
	}
	if pools, ok := d.GetOk("pools"); ok {
		anameAttr.Pools = toListOfInt(pools)
	}

	if contactid, ok := d.GetOk("contact_ids"); ok {
		anameAttr.ContactIDs = toListOfInt(contactid)
	}

	geoloc := &models.GeolocationANAME{}
	if geoipuserregion, ok := d.GetOk("geo_location"); ok {
		geouserlist := make([]int, 0, 1)
		tp := geoipuserregion.(map[string]interface{})
		if tp["geo_ip_user_region"] != nil {
			var1, _ := strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_user_region"]))
			geouserlist = append(geouserlist, var1)
			geoloc.GeoIpUserRegion = geouserlist
		}
		if tp["drop"] != nil {
			geoloc.Drop, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["drop"]))
		}
		if tp["geo_ip_failover"] != nil {
			geoloc.GeoIpFailOver, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["geo_ip_failover"]))
		}
		if tp["geo_ip_proximity"] != nil {
			geoloc.GeoIpProximity, _ = strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_proximity"]))
		}
	}
	anameAttr.GeoLocation = geoloc

	maplistrr := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roundrobin"); ok {
		tp := val.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			maplistrr = append(maplistrr, map1)
		}
		anameAttr.RoundRobin = maplistrr
	}

	valueslist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("record_failover_values"); ok {
		rcdfaname := &models.RCDFAname{}
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
		rcdfaname.Values = valueslist

		if failtype, ok := d.GetOk("record_failover_failover_type"); ok {
			rcdfaname.FailoverTypeAname, _ = strconv.Atoi(fmt.Sprintf("%v", failtype))
		}

		if disflag, ok := d.GetOk("record_failover_disable_flag"); ok {
			rcdfaname.DisableFlagAname, _ = strconv.ParseBool(fmt.Sprintf("%v", disflag))
		}

		rcdfaname.Values = valueslist
		anameAttr.RecordFailoverAname = rcdfaname
	}

	anamerecordid := d.Id()

	_, err := constellixClient.UpdatebyID(anameAttr, "v1/"+d.Get("source_type").(string)+"/"+d.Get("domain_id").(string)+"/records/aname/"+anamerecordid)
	if err != nil {
		return err
	}
	return resourceConstellixANAMERecordRead(d, m)
}

func resourceConstellixANAMERecordDelete(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	anamerecordid := d.Id()

	err := constellixClient.DeletebyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/aname/" + anamerecordid)
	if err != nil {
		return err
	}
	d.SetId("")
	return err
}
