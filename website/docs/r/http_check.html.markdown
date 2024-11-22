---
layout: "constellix"
page_title: "Constellix: constellix_http_check"
sidebar_current: "docs-constellix-resource-constellix_http_check"
description: |-
    Manages one or more HTTP check resource
---
# constellix_HTTP_Check_resource #
Manages one or more HTTP check resource.

# Example Usage #
```hcl
        
resource "constellix_http_check" "first" {
  name = "http check"
  host = "constellix.com"
  ip_version = "IPV4"
  port = 443
  protocol_type = "HTTPS"
  check_sites = [1,2]
  notification_groups = [874, 875]
}

```

## Argument Reference ##
* `name` - (Required) Name of the resource. Name should be unique.
* `host` - (Required) Host for the resource, for example "constellix.com". It can be set only once.
* `ip_version` - (Required) Specifies the version of IP. It can be set only once.
* `port` - (Required) Specifies the port number.
* `protocol_type` - (Required) Specifies upper layer protocol like HTTP, HTTPs, etc.
* `check_sites` - (Required) Site ids to check.
* `notification_groups` - (Optional) List of group IDs for the notification group of HTTP Check.
* `interval` - (Optional) Check Interval. Allowed values are `THIRTYSECONDS`, `ONEMINUTE`, `TWOMINUTES`, `THREEMINUTES`, `FOURMINUTES`, `FIVEMINUTES`, `TENMINUTES`, `THIRTYMINUTES`, `HALFDAY` and `DAY`.
* `interval_policy` - (Optional) Agent Interval Run Policy. It specifies whether you want to run checks from one location or all. Allowed values are `PARALLEL`, `ONCEPERSITE` and `ONCEPERREGION`.
* `verification_policy` - (Optional) Specifies how the check should be validated. Allowed values are `SIMPLE` and `MAJORITY`. This parameter will only work with the `interval_policy` set to `PARALLEL`.
* `fqdn` - (Optional) Fully qualified domain name of the URL should be checked.
* `path` - (Optional) In case of multi-page site, which path should be checked.
* `search_string` - (Optional) String to search in the first 2KB of resonse received.
* `expected_status_code` - (Optional) Expected HTTP status code for this check.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of HTTP check resource.

## Importing ##

An existing check can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_http_check.example <check-id>
```

Where check-id is the Id of check calculated via Constellix API.
