package constellix

import (
	"fmt"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "API key for HTTP call",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"apikey", "CONSTELLIX_API_KEY"}, nil),
			},

			"secretkey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret Key for HMAC",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"secretkey", "CONSTELLIX_SECRET_KEY"}, nil),
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allows insecure HTTPS client",
			},

			"proxyurl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Proxy server URL",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"constellix_domain":                  resourceConstellixDomain(),
			"constellix_a_record":                resourceConstellixARecord(),
			"constellix_aaaa_record":             resourceConstellixAAAARecord(),
			"constellix_aname_record":            resourceConstellixANAMERecord(),
			"constellix_cname_record":            resourceConstellixCNameRecord(),
			"constellix_hinfo_record":            resourceConstellixHinfo(),
			"constellix_http_redirection_record": resourceConstellixHTTPRedirection(),
			"constellix_mx_record":               resourceConstellixMX(),
			"constellix_naptr_record":            resourceConstellixNAPTR(),
			"constellix_caa_record":              resourceConstellixCaa(),
			"constellix_cert_record":             resourceConstellixCert(),
			"constellix_ns_record":               resourceConstellixNS(),
			"constellix_ptr_record":              resourceConstellixPtr(),
			"constellix_rp_record":               resourceConstellixRP(),
			"constellix_spf_record":              resourceConstellixSpf(),
			"constellix_srv_record":              resourceConstellixSRVRecord(),
			"constellix_txt_record":              resourceConstellixTxt(),
			"constellix_template":                resourceConstellixTemplate(),
			"constellix_a_record_pool":           resourceConstellixARecordPool(),
			"constellix_aaaa_record_pool":        resourceConstellixAAAArecordPool(),
			"constellix_cname_record_pool":       resourceConstellixCnameRecordPool(),
			"constellix_geo_filter":              resourceConstellixIPFilter(),
			"constellix_geo_proximity":           resourceConstellixGeoProximity(),
			"constellix_vanity_nameserver":       resourceConstellixVanityNameserver(),
			"constellix_contact_lists":           resourceConstellixContactList(),
			"constellix_tags":                    resourceConstellixTags(),
			"constellix_http_check":              resourceConstellixHTTPCheck(),
			"constellix_tcp_check":               resourceConstellixTCPCheck(),
			"constellix_dns_check":               resourceConstellixDNSCheck(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"constellix_a_record":                datasourceConstellixARecord(),
			"constellix_aaaa_record":             datasourceConstellixAAAARecord(),
			"constellix_aname_record":            datasourceConstellixAnamerecord(),
			"constellix_cname_record":            datasourceConstellixCNameRecord(),
			"constellix_a_record_pool":           datasourceConstellixARecordPool(),
			"constellix_cert_record":             datasourceConstellixCert(),
			"constellix_domain":                  datasourceConstellixDomain(),
			"constellix_caa_record":              datasourceConstellixCaa(),
			"constellix_contact_lists":           datasourceConstellixContactList(),
			"constellix_geo_proximity":           datasourceConstellixGeoProximity(),
			"constellix_http_redirection_record": datasourceConstellixHTTPRedirection(),
			"constellix_ptr_record":              datasourceConstellixPtr(),
			"constellix_rp_record":               datasourceConstellixRP(),
			"constellix_hinfo_record":            datasourceConstellixHinfo(),
			"constellix_mx_record":               datasourceConstellixMX(),
			"constellix_naptr_record":            datasourceConstellixNAPTR(),
			"constellix_ns_record":               datasourceConstellixNS(),
			"constellix_txt_record":              datasourceConstellixTxt(),
			"constellix_spf_record":              datasourceConstellixSPF(),
			"constellix_tags":                    datasourceConstellixTags(),
			"constellix_vanity_nameserver":       datasourceConstellixVanityNameserver(),
			"constellix_cname_record_pool":       datasourceConstellixCnamerecordPool(),
			"constellix_template":                datasourceConstellixTemplate(),
			"constellix_aaaa_record_pool":        datasourceConstellixAAAArecordpool(),
			"constellix_srv_record":              datasourceConstellixSRV(),
			"constellix_geo_filter":              datasourceConstellixIPFilter(),
			"constellix_http_check":              datasourceConstellixHTTPCheck(),
			"constellix_tcp_check":               datasourceConstellixTCPCheck(),
			"constellix_dns_check":               datasourceConstellixDNSCheck(),
		},

		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {
	config := config{
		apikey:    d.Get("apikey").(string),
		secretkey: d.Get("secretkey").(string),
		insecure:  d.Get("insecure").(bool),
		proxyurl:  d.Get("proxyurl").(string),
	}

	if err := config.Valid(); err != nil {
		return nil, err
	}
	cli := config.getClient()
	return cli, nil
}

func (c config) Valid() error {

	if c.apikey == "" {
		return fmt.Errorf("API Key is required")
	}

	if c.secretkey == "" {
		return fmt.Errorf("secret key is required")
	}
	return nil
}

func (c config) getClient() interface{} {

	return client.GetClient(c.apikey, c.secretkey)
}

type config struct {
	apikey    string
	secretkey string
	insecure  bool
	proxyurl  string
}
