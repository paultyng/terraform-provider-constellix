package constellix

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"constellix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	// We will use this function later on to make sure our test environment is valid.
	// For example, you can make sure here that some environment variables are set.
	if os.Getenv("apikey") == "" && os.Getenv("CONSTELLIX_API_KEY") == "" {
		t.Fatal("API KEY env variable must be set for acceptance tests")
	}

	if os.Getenv("secretkey") == "" && os.Getenv("CONSTELLIX_SECRET_KEY") == "" {
		t.Fatal("SECRET KEY env variable must be set for acceptance tests")
	}
}
