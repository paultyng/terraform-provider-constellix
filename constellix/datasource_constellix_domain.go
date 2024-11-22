package constellix

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/Jeffail/gabs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceConstellixDomain() *schema.Resource {
	return &schema.Resource{
		Read: datasourceConstellixDomainRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"has_gtd_regions": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"has_geoip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"vanity_nameserver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nameserver_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"soa": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_nameserver": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"email": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"ttl": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"refresh": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"serial": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"retry": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expire": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"negcache": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceConstellixDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	name := d.Get("name").(string)

	resp, err := client.GetbyId("v1/domains")
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	obj, err := gabs.ParseJSON(bodyBytes)
	if err != nil {
		return err
	}

	flag := false
	for i := 0; i < len(obj.Data().([]interface{})); i++ {
		if stripQuotes(obj.Index(i).S("name").String()) == name {
			flag = true

			soaset := make(map[string]interface{})
			soaset["primary_nameserver"] = stripQuotes(obj.Index(i).S("soa", "primaryNameserver").String())
			soaset["ttl"] = stripQuotes(obj.Index(i).S("soa", "ttl").String())
			if value, ok := d.GetOk("soa"); ok {
				tp := value.(map[string]interface{})
				if tp["email"] != nil {
					soaset["email"] = stripQuotes(obj.Index(i).S("soa", "email").String())
				}
			}
			soaset["refresh"] = stripQuotes(obj.Index(i).S("soa", "refresh").String())
			soaset["expire"] = stripQuotes(obj.Index(i).S("soa", "expire").String())
			soaset["retry"] = stripQuotes(obj.Index(i).S("soa", "retry").String())
			soaset["negcache"] = stripQuotes(obj.Index(i).S("soa", "negCache").String())

			d.Set("id", stripQuotes(obj.Index(i).S("id").String()))
			d.SetId(stripQuotes(obj.Index(i).S("id").String()))
			d.Set("name", stripQuotes(obj.Index(i).S("name").String()))
			d.Set("soa", soaset)
			if disabled, err := strconv.ParseBool(stripQuotes(obj.Index(i).S("disabled").String())); err == nil {
				d.Set("disabled", disabled)
			}
			if hasGeoIP, err := strconv.ParseBool(stripQuotes(obj.Index(i).S("hasGeoIP").String())); err == nil {
				d.Set("has_geoip", hasGeoIP)
			}
			if hasGTDRegion, err := strconv.ParseBool(stripQuotes(obj.Index(i).S("hasGtdRegions").String())); err == nil {
				d.Set("has_gtd_regions", hasGTDRegion)
			}
			if obj.Index(i).Exists("vanityNameServer") {
				d.Set("vanity_nameserver", stripQuotes(obj.Index(i).S("vanityNameServer").String()))
			}
			if obj.Index(i).Exists("nameserverGroup") {
				d.Set("nameserver_group", stripQuotes(obj.Index(i).S("nameserverGroup").String()))
			}
			if obj.Index(i).Exists("note") && obj.Index(i).S("note").String() != "{}" {
				d.Set("note", stripQuotes(obj.Index(i).S("note").String()))
			}

			if obj.Index(i).S("tags").Data() != nil {
				d.Set("tags", toListOfString(obj.Index(i).S("tags").Data()))
			} else {
				d.Set("tags", make([]string, 0, 1))
			}
		}
	}

	if flag != true {
		return fmt.Errorf("Domain of specified name is not available")
	}
	return nil
}
