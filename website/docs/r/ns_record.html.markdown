---
layout: "constellix"
page_title: "CONSTELLIX: constellix_ns_record"
sidebar_current: "docs-constellix-resource-constellix_ns_record"
description: |-
  Manages records of type NS  for a specific domain.
---

# constellix_ns_record
 Manages records of type NS  for a specific domain.

## Example Usage ##

```hcl
resource "constellix_ns_record" "firstrecord" {
  domain_id   = "${constellix_domain.first_domain.id}"
  source_type = "domains"
  ttl         = 100
  name        = "firstrecord"
  roundrobin {
    value        = "prac."
    disable_flag = "false"
  }
  type       = "NS"
  gtd_region = 1
  note       = "First record"
  noanswer   = false

}


```

## Argument Reference ##
* `domain_id` - (Required) Record id of NS record.
* `source_type` - (Required) Type of the NS record. The values which can be applied are `domains` or `templates`.
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `name` - (Optional) Name of record. Name should be unique.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is `false` (Active).
* `note` - (Optional)Record note.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `type` - (Optional) Record type `NS`.
* `roundrobin` - (Required) Set.
* `roundrobin.value` - (Required) This will be the host name for the name server, for example ns0.nameserver.com. It is important to note, the domain name is automatically appended to the end of this field unless it ends with a dot (.).
* `roundrobin.disable_flag` - (Required) disable flag. Default is `false`.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the NS resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_ns_record.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.