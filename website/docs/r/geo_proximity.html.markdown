---
layout: "constellix"
page_title: "CONSTELLIX: constellix_geo_proximity"
sidebar_current: "docs-constellix-resource-constellix_geo_proximity"
description: |-
  Manage Geoproximity for A, AAAA, CNAME or ANAME records. 
---

# constellix_geo_proximity
 Manage Geoproximity for A, AAAA, CNAME or ANAME records. 

## Example Usage ##

```hcl
resource "constellix_geo_proximity" "firstgeoproximity" {
  name      = "practice"
  latitude  = "22.7"
  longitude = "56.8333"
  region    = "05"
  city      = "273890"
  country   = "OM"
}


```

## Argument Reference ##
* `name` - (Required) Geo Proximity name should be unique.
* `country` - (Optional) Country code. Default is `null`.
* `region` - (Optional) Region or state or province code. Default is `null`.
* `latitude` - (Optional) Latitude value.
* `longitude` - (Optional) Longitude value.
* `city` - (Optional) City code. Default is `null`.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the Geo Proximity resource.

## Importing ##

An existing Geo Proximity can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_geo_proximity.example <geoproxy-id>
```

Where geoproxy-id is the Id of geoproxy calculated via Constellix API.