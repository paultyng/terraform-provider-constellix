package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/Constellix/constellix-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccA_Basic(t *testing.T) {
	var a models.ARecordAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConstellixARecordConfig_basic(1800),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixARecordExists("constellix_domain.domain1", "constellix_a_record.a1", &a),
					testAccCheckConstellixARecordAttributes(1800, &a),
				),
			},
		},
	})
}

func TestAccConstellixARecord_Update(t *testing.T) {
	var a models.ARecordAttributes

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConstellixARecordConfig_basic(1800),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixARecordExists("constellix_domain.domain1", "constellix_a_record.a1", &a),
					testAccCheckConstellixARecordAttributes(1800, &a),
				),
			},
			{
				Config: testAccCheckConstellixARecordConfig_basic(1900),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixARecordExists("constellix_domain.domain1", "constellix_a_record.a1", &a),
					testAccCheckConstellixARecordAttributes(1900, &a),
				),
			},
		},
	})
}

func testAccCheckConstellixARecordConfig_basic(ttl int) string {
	return fmt.Sprintf(`
	resource "constellix_domain" "domain1" {
		name = "checkarecord.com"
		soa = {
			email = "com.com."
			primary_nameserver = "ns41.constellix.com."
			ttl = 1900
			refresh = 48100
			retry = 7200
			expire = 1209
			negcache = 8000
		}
	}

	resource "constellix_a_record" "a1"{
		domain_id = "${constellix_domain.domain1.id}"
		source_type = "domains"
		name = "temparecord"
		ttl = "%d"
		note = "Practice record"

	  geo_location {
		geo_ip_failover  = true
		drop             = false
	  }

		roundrobin  {
			     value       = "16.45.25.35"
			     disable_flag = "false"
			     }
		roundrobin {
			       value = "15.45.25.30"
			       disable_flag = "true"
			   }
			  
	}
	`, ttl)
}

func testAccCheckConstellixARecordExists(domainName string, aName string, model *models.ARecordAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[domainName]
		rs2, err2 := s.RootModule().Resources[aName]

		if !err1 {
			return fmt.Errorf("Domain %s not found", domainName)
		}

		if !err2 {
			return fmt.Errorf("A record %s not found", aName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Domain id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No A record id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		resp, err := client.GetbyId("v1/domains/" + rs1.Primary.ID + "/records/a/" + rs2.Primary.ID)

		if err != nil {
			return err
		}

		tp, _ := arecordfromcontainer(resp)

		*model = *tp
		return nil
	}
}

func testAccCheckConstellixARecordAttributes(ttl interface{}, model *models.ARecordAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "temparecord" != model.Name {
			return fmt.Errorf("Bad A record name %s", model.Name)
		}
		ttl, _ := strconv.Atoi(fmt.Sprintf("%v", ttl))
		if ttl != model.TTL {
			return fmt.Errorf("Bad A record ttl %d", model.TTL)
		}

		return nil
	}
}

func testAccCheckConstellixARecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	rs1, err1 := s.RootModule().Resources["constellix_domain.domain1"]
	if !err1 {
		return fmt.Errorf("Domain %s not found", "constellix_domain.domain1")
	}
	domainid := rs1.Primary.ID
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "constellix_a_record" {
			_, err := client.GetbyId("v1/domains/" + domainid + "/records/a/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("A record is still exists")
			}
		} else {
			continue
		}
	}
	return nil
}

func arecordfromcontainer(resp *http.Response) (*models.ARecordAttributes, error) {

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	json.Unmarshal([]byte(bodyString), &data)
	model := models.ARecordAttributes{}
	model.Name = fmt.Sprintf("%v", data["name"])
	model.TTL, _ = strconv.Atoi(fmt.Sprintf("%v", data["ttl"]))
	model.Note = fmt.Sprintf("%v", data["note"])

	return &model, nil

}
