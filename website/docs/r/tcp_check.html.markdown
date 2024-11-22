---
layout: "constellix"
page_title: "Constellix: constellix_tcp_check"
sidebar_current: "docs-constellix-resource-constellix_tcp_check"
description: |-
    Manages one or more TCP check resource
---
# constellix_tcp_check #
Manages one or more TCP check resource.

# Example Usage #
```hcl
        
resource "constellix_tcp_check" "first" {
  name = "tcp check"
  host = "constellix.com"
  ip_version = "IPV4"
  port = 443
  check_sites = [1,2]
  notification_groups = [874, 875]
}

```

## Argument Reference ##
* `name` - (Required) Name of the resource. Name should be unique.
* `host` - (Required) Host for the resource, for example "constellix.com". It can be set only once.
* `ip_version` - (Required) Specifies the version of IP. It can be set only once.
* `port` - (Required) Specifies the port number.
* `check_sites` - (Required) Site ids to check.
* `notification_groups` - (Optional) List of group IDs for the notification group of TCP Check.
* `interval` - (Optional) Check Interval. Allowed values are `THIRTYSECONDS`, `ONEMINUTE`, `TWOMINUTES`, `THREEMINUTES`, `FOURMINUTES`, `FIVEMINUTES`, `TENMINUTES`, `THIRTYMINUTES`, `HALFDAY` and `DAY`.
* `interval_policy` - (Optional) Agent Interval Run Policy. It specifies whether you want to run checks from one location or all. Allowed values are `PARALLEL`, `ONCEPERSITE` and `ONCEPERREGION`.
* `verification_policy` - (Optional) Specifies how the check should be validated. Allowed values are `SIMPLE` and `MAJORITY`. This parameter will only work with the `interval_policy` set to `PARALLEL`.
* `string_to_send` - (Optional) String to send along with the check. It can be any parameter to the endpoint.
* `string_to_receive` - (Optional) String which should be received as a result of TCP check.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of TCP check resource.

## Importing ##

An existing Check can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_tcp_check.example <check-id>
```

Where check-id is the Id of check calculated via Constellix API.