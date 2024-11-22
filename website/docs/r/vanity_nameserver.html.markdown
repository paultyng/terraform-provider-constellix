---
layout: "constellix"
page_title: "Constellix: constellix_vanity_nameserver"
sidebar_current: "docs-constellix-resource-constellix_vanity_nameserver"
description: |-
    Manages Vanity Nameservers for domain records.
---
# constellix_vanity_nameserver #
Manages Vanity Nameservers for domain records.

# Example Usage #
```hcl
        
resource "constellix_vanity_nameserver" "vanitynameserver1" {
  name                   = "vanitynameserverrecord"
  nameserver_group       = 1
  nameserver_list_string = "www.google.com,\nwww.facebook.com,\nwww.instegram.com"
  is_default             = false
  is_public              = false
  nameserver_group_name  = "NS user group 1"
}


```

## Argument Reference ##
* `name` - (Required) Vanity nameserver name should be unique.
* `nameserver_group` - (Required) Name server group id. Available nameserver groups: `1`.
* `nameserver_list_string` - (Required) Comma separated name servers list
* `is_default` - (Optional) Default flag. Default is false.
* `is_public` - (Optional) isPublic flag. Default is false
* `nameserver_group_name` - (Optional) Name server group name.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of vanitynameserver resource.

## Importing ##

An existing Vanity Name Server can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_vanity_nameserver.example <nameserver-id>
```

Where nameserver-id is the Id of nameserver calculated via Constellix API.