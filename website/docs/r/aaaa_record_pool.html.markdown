---
layout: "constellix"
page_title: "CONSTELLIX: constellix_aaaa_record_pool"
sidebar_current: "docs-constellix-resource-constellix_aaaa_record_pool"
description: |-
   Manages the pools of AAAA records.
---

# constellix_aaaa_record_pool
 Manages the pools of AAAA records.

## Example Usage ##

```hcl
resource "constellix_aaaa_record_pool" "firstrecord" {
  name                   = "firstrecord"
  num_return             = 1
  min_available_failover = 1
  values {
    value        = "0:0:0:0:0:0:0:12"
    weight       = 20
    policy       = "followsonar"
    disable_flag = false
  }
  failed_flag  = false
  disable_flag = false
  note         = "First record"
}


```

## Argument Reference ##
* `name` - (Required) Pool name should be unique.
* `num_return` - (Required) Minimum number of value object to return. Value must be in between `0` and `64`.
* `min_available_failover` - (Required) Minimum number of Available Failover. Value must be in between `0` and `64`.
* `failed_flag` - (Optional) Failed flag. Default is `false`.
* `disable_flag` - (Optional) Enable or disable pool values. Default is `false`.
* `values` - (Required) Object Number of IP/Hosts in a pool values cannot be less than the "Num Return" and "Min Available" values.
* `values.value` - (Required) IPv6 address.
* `values.weight` - (Required) Weight number to sort the priorty. Weight must be in between `1` and `1000000`.
* `values.disable_flag` - (Optional) Enable or disable pool values. Default is `false`.
* `values.checkid` - (Optional) Sonar check id is required when you want to apply the ITO feature on a pool.
* `values.policy` - (Required) "followsonar" for Follow sonar. "alwaysoff" for Always off. "alwayson" for Always on. "offonfailure" for Off on Failure.
* `note` - (Optional) Description.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the aaaa record pool resource.

## Importing ##

An existing Pool can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_aaaa_record_pool.example <pool-id>
```

Where pool-id is the Id of record calculated via Constellix API.