package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/Constellix/constellix-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConstellixDomainCreation(t *testing.T) {
	testName1 := "terraform-domain-create-test-1"
	domainName1 := testName1 + ".test"
	resourceName1 := "constellix_domain." + testName1

	testName2 := "terraform-domain-create-test-2"
	domainName2 := testName2 + ".test"
	resourceName2 := "constellix_domain." + testName2

	var domain1, domain2 DomainAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Println("should successfully create a domain with note and disabled (set to false) attributes")
				},
				Config: testAccCheckConstellixDomainConfig(
					testName1,
					domainName1,
					"note-1",
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain1, resourceName1),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain1, domainName1, "note-1", false),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName1, "name", domainName1),
					resource.TestCheckResourceAttr(resourceName1, "note", "note-1"),
					resource.TestCheckResourceAttr(resourceName1, "disabled", "false"),
				),
			},
			{
				PreConfig: func() {
					log.Println("should successfully create a domain with note and disabled (set to true) attributes")
				},
				Config: testAccCheckConstellixDomainConfig(
					testName2,
					domainName2,
					"note-2",
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain2, resourceName2),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain2, domainName2, "note-2", true),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName2, "name", domainName2),
					resource.TestCheckResourceAttr(resourceName2, "note", "note-2"),
					resource.TestCheckResourceAttr(resourceName2, "disabled", "true"),
				),
			},
		},
	})
}

func TestAccConstellixDomainCreationIdempotency(t *testing.T) {
	testName := "terraform-domain-create-idempotent-reapply"
	domainName := testName + ".test"
	resourceName := "constellix_domain." + testName
	domainConfig := testAccCheckConstellixDomainConfig(
		testName,
		domainName,
		"note-1",
		false,
	)

	var domain DomainAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Println("should successfully create for a follow up test step to recreate (idempotent) domain creation")
				},
				Config: domainConfig,
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain, domainName, "note-1", false),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "note-1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
			{
				PreConfig: func() {
					log.Println("should be idempotent when reapplying the same domain creation configuration")
				},
				Config: domainConfig,
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain, domainName, "note-1", false),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "note-1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
		},
	})
}

func TestAccConstellixDomainCreationExisting(t *testing.T) {
	// Should be able to import existing domain via create operation"
	// when domain metadata in terraform config matches domain's metadata on server"
	testName := "terraform-domain-create-import-existing-same-metadata"
	domainName := testName + ".test"
	resourceName := "constellix_domain." + testName

	domainID, err := givenDomainOnServer(domainName, "created outside terraform", true)
	if err != nil {
		log.Println("error creating test domain", err)
		t.FailNow()
	}
	// Delete domain in case missed by the follow-up destroy step.
	t.Cleanup(cleanupDomain(domainID))

	var domain DomainAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Println("should be able to import an existing domain (created outside terraform)")
				},
				Config: testAccCheckConstellixDomainConfig(
					testName,
					domainName,
					"created outside terraform",
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain, domainName, "created outside terraform", true),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "created outside terraform"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					// Ensure has expected domain ID
					testAccCheckConstellixDomainHasDomainID(domainID, resourceName),
				),
			},
		},
	})
}

func testAccCheckConstellixDomainHasDomainID(expectedDomainID string, resourceName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("domain resource %s not found", resourceName)
		}
		if expectedDomainID != rs.Primary.ID {
			return fmt.Errorf("domain ids do not match, %s != %s", expectedDomainID, rs.Primary.ID)
		}
		return nil
	}
}

func TestAccConstellixDomainCreationFailure(t *testing.T) {
	testName := "terraform-domain-create-failure-1"
	invalidDomainName := "terraform_test_invalid_domain_name"
	resourceName := "constellix_domain." + testName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Println("should error when attempt to create a domain with an invalid name")
				},
				Config: testAccCheckConstellixDomainConfig(
					testName,
					invalidDomainName,
					"",
					false,
				),
				ExpectError: regexp.MustCompile(`"terraform_test_invalid_domain_name" is not a valid domain name`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstellixDomainDoesNotExist(resourceName),
				),
			},
		},
	})
}

func TestAccConstellixDomainUpdate(t *testing.T) {
	testName := "terraform-domain-update"
	domainName := testName + ".test"
	resourceName := "constellix_domain." + testName
	initialConfig := testAccCheckConstellixDomainConfig(
		testName,
		domainName,
		"note-1",
		false,
	)
	updatedConfig1 := testAccCheckConstellixDomainConfig(
		testName,
		domainName,
		"note-2",
		true,
	)
	updatedConfig2 := testAccCheckConstellixDomainConfig(
		testName,
		domainName,
		"note-3",
		false,
	)

	var domain1, domain2, domain3 DomainAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Println("should create a domain for follow up steps testing update operation on a domain")
				},
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain1, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain1, domainName, "note-1", false),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "note-1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
			{
				PreConfig: func() {
					log.Println("should update values of note and disabled attributes (set to true) of a domain")
				},
				Config: updatedConfig1,
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain2, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain2, domainName, "note-2", true),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "note-2"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
			{
				PreConfig: func() {
					log.Println("should update values of note and disabled attributes (set to false) of a domain")
				},
				Config: updatedConfig2,
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain3, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain3, domainName, "note-3", false),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "note-3"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
		},
	})
}

func TestAccConstellixDomainImport(t *testing.T) {
	testName := "terraform-domain-import"
	domainName := testName + ".test"
	resourceName := "constellix_domain." + testName

	var domain DomainAttributes
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConstellixDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Println("should create a domain for follow up steps testing import operation on domain")
				},
				Config: testAccCheckConstellixDomainConfig(
					testName,
					domainName,
					"note-1",
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					// Load domain from API.
					testAccCheckConstellixDomainExists(&domain, resourceName),
					// Check if load values are correct.
					testAccCheckConstellixDomainAttributes(&domain, domainName, "note-1", true),
					// Check if the values inside terraform state are correct.
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "note", "note-1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
			{
				PreConfig: func() {
					log.Println("should validate import state of a domain")
				},
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckConstellixDomainConfig(testName, domainName string, note string, disabled bool) string {
	return fmt.Sprintf(`
	resource "constellix_domain" "%s" {
		name = "%s"
		soa = {
			ttl = 1800
			primary_nameserver = "ns41.constellix.com."
			email = "dns.constellix.com."
			refresh = 48100
			retry = 7200
			expire = 1209
			negcache = 8000
		}
		note = "%s"
		disabled = "%t"
	}
	`, testName, domainName, note, disabled)
}

func testAccCheckConstellixDomainExists(domain *DomainAttributes, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("domain resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no domain id was set")
		}

		resp, err := loadDomainFromServer(rs.Primary.ID)
		if err != nil {
			return err
		}

		tp, _ := domainFromResponse(resp)
		*domain = *tp
		return nil

	}
}

func loadDomainFromServer(domainID string) (*http.Response, error) {
	cl := testAccProvider.Meta().(*client.Client)
	resp, err := cl.GetbyId("v1/domains/" + domainID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func cleanupDomain(domainID string) func() {
	return func() {
		cl := givenClient()
		err := cl.DeletebyId("v1/domains/" + domainID)
		if err != nil && !strings.Contains(err.Error(), "Domain not found with id "+domainID) {
			log.Println("error deleting domain with ID", domainID)
		}
	}
}

func givenDomainOnServer(domainName, note string, disabled bool) (string, error) {
	domainAttr := DomainAttributes{
		Disabled: disabled,
		DomainAttributes: models.DomainAttributes{
			Name:          []string{domainName},
			HasGtdRegions: false,
			HasGeoIP:      false,
			Note:          note,
			Soa: &models.Soa{
				PrimaryNameServer: "ns41.constellix.com.",
				Email:             "dns.constellix.com.",
				TTL:               "1800",
				Refresh:           "48100",
				Retry:             "7200",
				Expire:            "1209",
				NegCache:          "8000",
			},
		},
	}
	cl := givenClient()
	resp, err := cl.Save(domainAttr, "v1/domains")
	if err != nil {
		return "", err
	}
	return extractDomainIDFromDomainCreationResponse(resp.Body)
}

func givenClient() *client.Client {
	return client.GetClient(os.Getenv("apikey"), os.Getenv("secretkey"))
}

func domainFromResponse(resp *http.Response) (*DomainAttributes, error) {
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
	domain := DomainAttributes{}

	// FIXME avoid using converters (like toStringList) both in production and test code.
	// TODO Revamp codebase to use correct data model, allowing direct unmarshalling instead of manual mapping of fields.
	nameList := toStringList(data["name"])
	domain.Name = nameList
	domain.Note = data["note"].(string)
	domain.Disabled = data["disabled"].(bool)

	return &domain, nil

}

func testAccCheckConstellixDomainDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("expected domain resource %s to not exist, but it was found in the state", resourceName)
		}
		return nil
	}
}

func testAccCheckConstellixDomainDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "constellix_domain" {
			_, err := loadDomainFromServer(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("domain is still exists, id: %s", rs.Primary.ID)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckConstellixDomainAttributes(domain *DomainAttributes, expectedName, expectedNote string, expectedDisabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if expectedName != domain.Name[0] {
			return fmt.Errorf("bad domain name %s", domain.Name)
		}
		if expectedNote != domain.Note {
			return fmt.Errorf("bad domain note %s", domain.Note)
		}
		if expectedDisabled != domain.Disabled {
			return fmt.Errorf("bad domain's disabled value %t", domain.Disabled)
		}
		return nil
	}
}
