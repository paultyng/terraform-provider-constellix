---
layout: "constellix"
page_title: "Constellix: constellix_txt_record"
sidebar_current: "docs-constellix-resource-constellix_txt"
description: |-
    Manages records of type TXT for a specific domain.
---
# constellix_txt_record #
Manages records of type TXT for a specific domain.

# Example Usage #
```hcl
resource "constellix_txt_record" "txtrecord1" {
  domain_id   = "${constellix_domain.domain1.id}"
  ttl         = 1800
  name        = "txtrecord"
  noanswer    = false
  note        = ""
  gtd_region  = 1
  type        = "TXT"
  source_type = "domains"
  roundrobin {
    value        = "\"{\\\"cfg\\\":[{\\\"useAS\\\":0}]}\""
    disable_flag = false
  }
}


```

## Argument Reference ##
* `domain_id` - (Required) Record id of TXT record.
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `source_type` - (Required) `domains` for Domain records and `template` for Template records.
* `roundrobin` - (Required) Set.
* `roundrobin.value` - (Required) Free form text data of any type which may be no longer than 255 characters unless divided into multiple strings with sets of quotation marks.
* `roundrobin.disable_flag` - (Optional) Disable flag. Default is false.
* `name` - (Optional) Name of record. Name should be unique.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is false (Active).
* `note` - (Optional) Record note.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `type` - (Optional) Record type `TXT`.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of txt resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_txt_record.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.