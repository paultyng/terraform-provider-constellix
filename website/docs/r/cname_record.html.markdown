---
layout: "constellix"
page_title: "CONSTELLIX: constellix_cname_record"
sidebar_current: "docs-constellix-resource-constellix_cname_record"
description: |-
  Manages one or more domain CNAME records.
---

# constellix_cname_record
Manages one or more domain CNAME records.

## Example Usage ##

```hcl
resource "constellix_cname_record" "firstrecord" {
  domain_id     = "${constellix_domain.first_domain.id}"
  source_type   = "domains"
  record_option = "failover"
  ttl           = 100
  name          = "cnamerecord"
  host          = "abcd.com."
  geo_location = {
    geo_ip_user_region = 1
    drop               = "false"
  }
  pools       = [123]
  contact_ids = [1234]
  type        = "CNAME"
  gtd_region  = 1
  note        = "First record"
  noanswer    = false
  record_failover_values {
    value        = "abc.com."
    sort_order   = 1
    disable_flag = "false"
  }
  record_failover_values {
    value        = "ab.com."
    sort_order   = 2
    disable_flag = "false"
  }
  record_failover_failover_type = 2
  record_failover_disable_flag  = "false"
}


```

## Argument Reference ##
* `domain_id` - (Required) Domain ID under which CNAME record should be created.
* `source_type` - (Required) Type of the CName record. The values which can be applied are `domains` or `templates`.
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `name` - (Optional) Name of record. Name should be unique.
* `host` - (Required for standard CNAME) Value/"alias to" of the CNAME record.
* `geo_location` - (Optional) Details of IP filter / Geo proximity to be applied. Default is null.
* `geo_location.geo_ip_user_region` - (Optional) For Geo proximity to be applied. geoipUserRegion should not be provided.
* `geo_location.drop` - (Optional) Drop flag. Default is `false`.
* `geo_location.geo_ip_proximity` - (Optional) A valid geoipProximity id.
* `geo_location.geo_ip_user_region` - (Optional) For Geo IP Filter to be applied geo_ip_proximity must not be provided. Before applying a specific IP Filter you must first create a record with the same name that has IP Filter setting of "World Default". geoipUserRegion should be `1` for "World Default". Otherwise, use a valid IP Filter id number.
* `geo_location.drop` - (Optional) drop flag. Default is `false`.
* `geo_location.geo_ip_failover` - (Optional) Flag to enable/disable Failover to nearest proximity when all the host fails. Works with the record type pools and Failover. It requires Geo Proximity to be enabled at the Domain level and applied to the record you are enabeling the geo_ip_filter option on. Default is "false" mark "true" to enable. 
* `geo_location.geo_ip_proximity` - (Optional) For Geo IP Filter, geoipProximity must not be provided. please create an A record with "World (Default)" IP Filter first before a more specific IP Filter is applied. The "World (Default)" record would only be used if no matching Filter or Proximity records are found.
* `record_option` - (Optional) Type of record. "roundRobin" for Standard record (Default). "failover" for Failover. "pools" for Pools. "roundRobinFailover" for Round Robin with Failover.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is `false` (Active).
* `note` - (Optional) Record note.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `type` - (Optional) Record type CNAME.
* `contact_ids` - (Optional) Applied contact list id. Only applicable to record with type roundRobin with failover and failover.
* `pools` - (Optional) Ids of CNamepool.
* `record_failover` - (Optional) To create a record failover object pass the following attributes.
* `record_failover_values` - (Required for failover) Set. 
* `record_failover_values.value` - (Required for failover) Host name.
* `record_failover_values.check_id` - (Optional) Sonar check id.
* `record_failover_values.sort_order` - (Required for failover) Integer value which decides in which order the recordfailover should be sorted.
* `record_failover_values.disable_flag` - (Required for failover) Enable or disable the recordFailover value object. Default is `false` (Active). At least one recordFailover value object should be false.
* `record_failover_failover_type` - (Required for failover) `1` for Normal (always lowest level). `2` for Off on any Failover event. `3` for One Way (move to higher level).
* `record_failover_disable_flag` - (Required for failover) enable or disable the recordFailover object. Default is `false` (Active). At least one recordFailover object should be false.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the CNAME resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_cname_record.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.
