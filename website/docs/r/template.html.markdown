---
layout: "constellix"
page_title: "CONSTELLIX: constellix_template"
sidebar_current: "docs-constellix-resource-constellix_template"
description: |-
   Manages the DNS record templates.
---

# constellix_template
  Manages the DNS record templates.

## Example Usage ##

```hcl
resource "constellix_template" "firsttemplate" {
  name            = "sample"
  has_gtd_regions = "true"
  has_geoip       = "false"
}



```

## Argument Reference ##
* `name` - (Required) Template names. e.g "sampletemplate".
* `domain` - (Optional) Id of domain to be applied.
* `has_gtd_regions` - (Optional) Enable/Disable GTD Region of the domain. The Default value is false.
* `has_geoip` - (Optional) Enable/Disable GEO IP. The Default value is false.
* `version` - (Optional) System generated template history version.

## Attributes Reference
This resource exports the following attributes:
* `id` - The constellix calculated id of the template resource.

## Importing ##

An existing Template can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_template.example <template-id>
```

Where template-id is the Id of template calculated via Constellix API.