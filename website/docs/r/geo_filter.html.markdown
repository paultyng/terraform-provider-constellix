---
layout: "constellix"
page_title: "CONSTELLIX: constellix_geo_filter"
sidebar_current: "docs-constellix-resource-constellix_geo_filter"
description: |-
  Manage Geofilters for A, AAAA, CNAME or ANAME records. 
---

# constellix_geo_filter
 Manage Geofilters for A, AAAA, CNAME or ANAME records. 

## Example Usage ##

```hcl
resource "constellix_geo_filter" "ipfilter1" {
  name               = "first135"
  geoip_continents   = ["AS"]
  geoip_countries    = ["IN", "PK"]
  geoip_regions      = ["IN/BR", "IN/MP"]
  asn                = [1, 2]
  ipv4               = ["1.1.1.0/32", "1.1.2.2/32"]
  ipv6               = ["2:0:0:2:0:0:1:abc/128"]
  filter_rules_limit = 100
}


```

## Argument Reference ##
* `name` - (Required) Geo Filter name should be unique.
* `geoip_continents` - (Optional) Two digit Continents Code.
* `geoip_countries` - (Optional) Two digit Countries Code.
* `geoip_regions` - (Optional) Two digit country code followed by "/" followed by two digit region code. 
* `asn` - (Optional) Autonomous System Number (ASN). ASN code should be a number between `0` and `4294967295`.
* `ipv4` - (Optional) IPV4 Address.
* `ipv6` - (Optional) IPV6 Address.
* `filter_rules_limit` - (Optional) Default is `100`. For more than 100 rules, parameter should be set explicitly for ADD and Update API calls. Value should be in mulitple of 100 like 200, 300 ...upto the quota limit assigned to the account. Check quota details for IP Filter Rule Limit.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the Geo Filter.

## Importing ##

An existing Geo Filter can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_geo_filter.example <filter-id>
```

Where filter-id is the Id of filter calculated via Constellix API.