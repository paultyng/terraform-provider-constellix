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

func TestAccCName_Basic(t *testing.T) {
	var cname models.CRecordAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixCNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConstellixCNameConfig_basic(1800),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixCNameExists("constellix_domain.domain1", "constellix_cname_record.cname1", &cname),
					testAccCheckConstellixCNameAttributes(1800, &cname),
				),
			},
		},
	})
}

func TestAccConstellixCName_Update(t *testing.T) {
	var cname models.CRecordAttributes

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixCNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConstellixCNameConfig_basic(1800),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixCNameExists("constellix_domain.domain1", "constellix_cname_record.cname1", &cname),
					testAccCheckConstellixCNameAttributes(1800, &cname),
				),
			},
			{
				Config: testAccCheckConstellixCNameConfig_basic(1900),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixCNameExists("constellix_domain.domain1", "constellix_cname_record.cname1", &cname),
					testAccCheckConstellixCNameAttributes("1900", &cname),
				),
			},
		},
	})
}

func testAccCheckConstellixCNameConfig_basic(ttl int) string {
	return fmt.Sprintf(`
	resource "constellix_domain" "domain1" {
		name = "checkcname.com"
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

	resource "constellix_cname_record" "cname1"{
		domain_id = "${constellix_domain.domain1.id}"
		source_type = "domains"
		name = "tempcnamerecord"
		ttl = "%d"
		note = "Practice record naptr"

	  geo_location {
		geo_ip_failover  = true
		drop             = false
	  }
		record_option = "failover"
	    record_failover_values  {
			     value = "a."
			     sort_order = 2
			     disable_flag = "false"
			   }
			   record_failover_values  {
				value = "c."
				sort_order = 3
				disable_flag = "false"
			  }
			   record_failover_failover_type = 1
			   record_failover_disable_flag = "false"
	}
	`, ttl)
}

func testAccCheckConstellixCNameExists(domainName string, cName string, cname *models.CRecordAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[domainName]
		rs2, err2 := s.RootModule().Resources[cName]

		if !err1 {
			return fmt.Errorf("Domain %s not found", domainName)
		}

		if !err2 {
			return fmt.Errorf("CName record %s not found", cName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Domain id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No CName record id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		resp, err := client.GetbyId("v1/domains/" + rs1.Primary.ID + "/records/cname/" + rs2.Primary.ID)

		if err != nil {
			return err
		}

		tp, _ := cnamefromcontainer(resp)

		*cname = *tp
		return nil
	}
}

func testAccCheckConstellixCNameAttributes(ttl interface{}, model *models.CRecordAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "tempcnamerecord" != model.Name {
			return fmt.Errorf("Bad CName record name %s", model.Name)
		}
		ttl, _ := strconv.Atoi(fmt.Sprintf("%v", ttl))
		if ttl != model.TTL {
			return fmt.Errorf("Bad CName record ttl %d", model.TTL)
		}

		return nil
	}
}

func testAccCheckConstellixCNameDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	rs1, err1 := s.RootModule().Resources["constellix_domain.domain1"]
	if !err1 {
		return fmt.Errorf("Domain %s not found", "constellix_domain.domain1")
	}
	domainid := rs1.Primary.ID
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "constellix_cname_record" {
			_, err := client.GetbyId("v1/domains/" + domainid + "/records/cname/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("CName record still exists")
			}
		} else {
			continue
		}
	}
	return nil
}

func cnamefromcontainer(resp *http.Response) (*models.CRecordAttributes, error) {

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	json.Unmarshal([]byte(bodyString), &data)
	model := models.CRecordAttributes{}

	model.Name = fmt.Sprintf("%v", data["name"])
	model.TTL, _ = strconv.Atoi(fmt.Sprintf("%v", data["ttl"]))
	model.Note = fmt.Sprintf("%v", data["note"])

	return &model, nil

}
