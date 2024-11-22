---
layout: "constellix"
page_title: "Constellix: constellix_mx_record"
sidebar_current: "docs-constellix-resource-constellix_mx_record"
description: |-
    Manages records of type MX for a specific domain.
---
# constellix_mx_record #
Manages records of type MX for a specific domain.

# Example Usage #
```hcl

resource "constellix_mx_record" "mx1" {
  domain_id   = "${constellix_domain.domain1.id}"
  source_type = "domains"
  name        = "mxrecord"
  ttl         = "1900"
  noanswer    = false
  note        = ""
  gtd_region  = 1
  type        = "MX"
  roundrobin {
    value        = "abc"
    level        = "100"
    disable_flag = "false"
  }
  roundrobin {
    value        = "dce"
    level        = "200"
    disable_flag = "true"
  }
}

```

## Argument Reference ##
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `source_type` - (Required) `domains` for Domain records and `template` for Template records.
* `roundrobin` - (Required) Set.
* `roundrobin.value` - (Required) The mail server that will accept mail for the host that is specified in the name field. Your domain name is automatically appended to your value if it does not end it a dot.
* `roundrobin.level` - (Required) Level must be in between `0` and `65535`. The MX level determines the order (by priority) that remote mail servers will attempt to deliver email. The mail server with the lowest MX level will be the first priority.
* `roundrobin.disable_flag` - (Optional) Enable or Disable the roundrobin object. Default is `false`. At least one roundrobin object should be false.
* `name` - (Optional) Name of record. Name should be unique.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is `false` (Active).
* `note` - (Optional) Record note.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `type` - (Optional) Record type `MX`.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of mx resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_mx_record.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.