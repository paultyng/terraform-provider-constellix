---
layout: "constellix"
page_title: "Constellix: constellix_caa_record"
sidebar_current: "docs-constellix-resource-constellix_caa_record"
description: |-
    Manages records of type CAA  for a specific domain.
---
# constellix_caa_record #
Manages records of type CAA  for a specific domain.

# Example Usage #
```hcl

resource "constellix_caa_record" "caacheck" {
  domain_id   = "${constellix_domain.domain1.id}"
  source_type = "domains"
  name        = "caarecord"
  ttl         = 1900
  noanswer    = false
  note        = ""
  gtd_region  = 1
  type        = "CAA"
  roundrobin {
    caa_provider_id = 3
    tag             = "issue"
    data            = "como.com"
    flag            = "0"
    disable_flag    = "false"
  }
  roundrobin {
    caa_provider_id = 4
    tag             = "issue"
    data            = "como01.com"
    flag            = "1"
    disable_flag    = "true"
  }
}

```

## Argument Reference ##
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `source_type` - (Required) `domains` for Domain records and `template` for Template records.
* `roundrobin` - (Required) Set.
* `roundrobin.caa_provider_id` - (Required) 
  * `1` for [ Custom ], 
  * `2` for [ No Provider ], 
  * `3` for Comodo, 
  * `4` for Digicert, 
  * `5` for Entrust, 
  * `6` for GeoTrust, 
  * `7` for Izenpe, 
  * `8` for Lets Encrypt, 
  * `9` for Symantec, 
  * `10` for Thawte

* `roundrobin.tag` - (Required) "issue" for Issue, "IssueWild" for IssueWild, "iodef" for iodef. Type allows you to choose how you want certificates to be issued by the CA. Each CAA record can contain only one tag-value pair. Options:
issue: Explicitly authorizes a single certificate authority to issue a certificate (any type) for the hostname.

issuewild: Authorization to issue certificates that specify a wildcard domain. Please note: issuewild properties take precedence over issue properties when specified.

iodef: (Incident Description Exchange Format) Specifies a means of reporting certificate issue requests or cases of certificate issue for the corresponding domain that violate the security policy of the issuer or the domain name holder.

* `roundrobin.data` - (Required) 
  * `""` for [ Custom ] if CAA provider Id is 1, 
  * `;` for [ No Provider ], 
  * `comodoca.com` for Comodo, 
  * `digicert.com` for Digicert, 
  * `entrust.net` for Entrust, 
  * `geotrust.com` for GeoTrust, 
  * `izenpe.com` for Izenpe, 
  * `letsencrypt.org` for Lets Encrypt, 
  * `symantec.com` for Symantec, 
  * `thawte.com` for Thawte.

* `roundrobin.flag` - (Required) roundRobin. Issuer Critical.

There is currently only one flag defined, “issuer critical” at a value of 1. If a CA does not understand the flag value for an issuer critical record, then the CA will return with “no issue” for the certification.

All records will have the default issuer critical value of 0, which means they are “not critical”. Issuer Critical Value should be between `0` to `255`.

* `roundrobin.disable_flag` - (Required) Disable flag. Default is `false`.

* `name` - (Optional) Name of record. Name should be unique.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is `false` (Active).
* `note` - (Optional) Record note.
* `type` - (Optional) Record type `CAA`.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of caa resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_caa_record.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.