---
layout: "constellix"
page_title: "CONSTELLIX: constellix_spf_record"
sidebar_current: "docs-constellix-resource-constellix_spf_record"
description: |-
  Manages records of type SPF  for a specific domain.
---

# constellix_spf_record
Manages records of type SPF  for a specific domain.

## Example Usage ##

```hcl
resource "constellix_spf_record" "spf1" {
  domain_id   = "${constellix_domain.domain1.id}"
  source_type = "domains"
  name        = "temp"
  ttl         = 10
  noanswer    = false
  gtd_region  = 1
  type        = "SPF"
  note        = "Practice record"
  roundrobin {
    value        = "124.56.8.1"
    disable_flag = "false"
  }

}


```

## Argument Reference ##
* `domain_id` - (Required) Record id of SPF record.
* `source_type` - (Required) Type of the PTR record. The values which can be applied are `domains` or `templates`.
* `name` - (Optional) Name of record. Name should be unique.
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is `false` (Active).
* `note` - (Optional) Record note.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `type` - (Optional) Record type `SPF`.
* `roundrobin` - (Required) Object.
* `roundrobin.value` - (Required) Value may contain multiple strings (each string enclosed in double quotes). Individual string length should not exceed 255 characters.
* `roundrobin.disable_flag` - (Optional) enable or disable the roundrobin object. Default is `false`. At least one roundrobin object should be false.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the SPF resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_spf_record.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.