package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/Constellix/constellix-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConstellixTxt() *schema.Resource {
	return &schema.Resource{
		Create: resourceConstellixTxtCreate,
		Update: resourceConstellixTxtUpdate,
		Read:   resourceConstellixTxtRead,
		Delete: resourceConstellixTxtDelete,

		Importer: &schema.ResourceImporter{
			State: resourceConstellixTxtImport,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"noanswer": &schema.Schema{
				Type:     schema.TypeString,
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

			"parentid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"parent": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"roundrobin": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"disable_flag": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceConstellixTxtImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	constellixClient := m.(*client.Client)
	params := strings.Split(d.Id(), ":")
	resp, err := constellixClient.GetbyId("v1/" + params[0] + "/" + params[1] + "/records/txt/" + params[2])
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil, err
		}
		return nil, err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodystring := string(bodybytes)

	var data map[string]interface{}
	json.Unmarshal([]byte(bodystring), &data)
	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	d.Set("name", data["name"])
	d.Set("ttl", data["ttl"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("note", data["note"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("parentid", data["parentId"])
	d.Set("parent", data["parent"])
	d.Set("source", data["source"])

	resrr := (data["roundRobin"]).([]interface{})
	mapListRR := make([]interface{}, 0, 1)
	for _, val := range resrr {
		tpMap := make(map[string]interface{})
		inner := val.(map[string]interface{})
		tpMap["value"] = stripQuotes(inner["value"].(string)) // removing the quotes added by the server during the GET call
		tpMap["disable_flag"] = inner["disableFlag"].(bool)
		mapListRR = append(mapListRR, tpMap)
	}

	d.Set("roundrobin", mapListRR)
	d.Set("domain_id", params[1])
	d.Set("source_type", params[0])
	log.Printf("[DEBUG] %s finished import", d.Id())
	return []*schema.ResourceData{d}, nil
}
func resourceConstellixTxtCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	TxtAttr := models.TxtAttributes{}

	if nm, ok := d.GetOk("name"); ok {
		TxtAttr.Name = nm.(string)
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		TxtAttr.TTL = ttl.(int)
	}

	if noans, ok := d.GetOk("noanswer"); ok {
		TxtAttr.NoAnswer = noans.(string)
	}

	if note, ok := d.GetOk("note"); ok {
		TxtAttr.Note = note.(string)
	}

	if gtdr, ok := d.GetOk("gtd_region"); ok {
		TxtAttr.GtdRegion = gtdr.(int)
	}

	if tp, ok := d.GetOk("type"); ok {
		TxtAttr.Type = tp.(string)
	}

	if pid, ok := d.GetOk("parentid"); ok {
		TxtAttr.ParentID = pid.(int)
	}

	if p, ok := d.GetOk("parent"); ok {
		TxtAttr.Parent = p.(string)
	}

	if source, ok := d.GetOk("source"); ok {
		TxtAttr.Source = source.(string)
	}

	if rr, ok := d.GetOk("roundrobin"); ok {
		mapListRR := make([]interface{}, 0, 1)
		tp := rr.([]interface{})
		for _, val := range tp {
			tpMap := make(map[string]interface{})
			inner := val.(map[string]interface{})
			tpMap["value"] = fmt.Sprintf("%v", inner["value"])
			tpMap["disableFlag"] = fmt.Sprintf("%v", inner["disable_flag"])

			mapListRR = append(mapListRR, tpMap)
		}
		TxtAttr.RoundRobin = mapListRR

	}

	id := d.Get("domain_id").(string)
	stid := d.Get("source_type").(string)

	resp, err := client.Save(TxtAttr, "v1/"+stid+"/"+id+"/records/txt")
	if err != nil {
		return err
	}

	//Managing response and extracting id of resource
	bodybtes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodystring := string(bodybtes)
	var data map[string]interface{}
	json.Unmarshal([]byte(bodystring[1:len(bodystring)-1]), &data)

	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	return resourceConstellixTxtRead(d, m)
}

func resourceConstellixTxtUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	TxtAttr := models.TxtAttributes{}

	if nm, ok := d.GetOk("name"); ok {
		TxtAttr.Name = nm.(string)
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		TxtAttr.TTL = ttl.(int)
	}

	if noans, ok := d.GetOk("noanswer"); ok {
		TxtAttr.NoAnswer = noans.(string)
	}

	if note, ok := d.GetOk("note"); ok {
		TxtAttr.Note = note.(string)
	}

	if gtdr, ok := d.GetOk("gtd_region"); ok {
		TxtAttr.GtdRegion = gtdr.(int)
	}

	if tp, ok := d.GetOk("type"); ok {
		TxtAttr.Type = tp.(string)
	}

	if pid, ok := d.GetOk("parentid"); ok {
		TxtAttr.ParentID = pid.(int)
	}

	if p, ok := d.GetOk("parent"); ok {
		TxtAttr.Parent = p.(string)
	}

	if source, ok := d.GetOk("source"); ok {
		TxtAttr.Source = source.(string)
	}
	if rr, ok := d.GetOk("roundrobin"); ok {
		mapListRR := make([]interface{}, 0, 1)
		tp := rr.([]interface{})
		for _, val := range tp {
			tpMap := make(map[string]interface{})
			inner := val.(map[string]interface{})
			tpMap["value"] = fmt.Sprintf("%v", inner["value"])
			tpMap["disableFlag"] = fmt.Sprintf("%v", inner["disable_flag"])

			mapListRR = append(mapListRR, tpMap)
		}
		TxtAttr.RoundRobin = mapListRR
	}

	domainID := d.Get("domain_id").(string)
	stid := d.Get("source_type").(string)
	txtid := d.Id()
	_, err := client.UpdatebyID(TxtAttr, "v1/"+stid+"/"+domainID+"/records/txt/"+txtid)
	if err != nil {
		return err
	}
	return resourceConstellixTxtRead(d, m)
}

func resourceConstellixTxtRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	domainID := d.Get("domain_id").(string)
	stid := d.Get("source_type").(string)
	txtid := d.Id()

	resp, err := client.GetbyId("v1/" + stid + "/" + domainID + "/records/txt/" + txtid)
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodystring := string(bodybytes)

	var data map[string]interface{}
	json.Unmarshal([]byte(bodystring), &data)
	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	d.Set("name", data["name"])
	d.Set("ttl", data["ttl"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("note", data["note"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("parentid", data["parentId"])
	d.Set("parent", data["parent"])
	d.Set("source", data["source"])
	resrr := (data["roundRobin"]).([]interface{})
	mapListRR := make([]interface{}, 0, 1)
	for _, val := range resrr {
		tpMap := make(map[string]interface{})
		inner := val.(map[string]interface{})
		value := stripQuotes(inner["value"].(string))
		value = strings.ReplaceAll(value, "\" \"", "")
		tpMap["value"] = value
		tpMap["disable_flag"] = inner["disableFlag"].(bool)
		mapListRR = append(mapListRR, tpMap)
	}

	d.Set("roundrobin", mapListRR)
	return nil
}

func resourceConstellixTxtDelete(d *schema.ResourceData, m interface{}) error {
	constellixConnect := m.(*client.Client)
	domainID := d.Get("domain_id").(string)
	stid := d.Get("source_type").(string)
	dn := d.Id()

	err := constellixConnect.DeletebyId("v1/" + stid + "/" + domainID + "/records/txt/" + dn)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
