---
layout: "constellix"
page_title: "Constellix: constellix_tags"
sidebar_current: "docs-constellix-resource-constellix_tags"
description: |-
    Manages tags for Constellix DNS.
---
# constellix_tags #
Manages tags for Constellix DNS.

# Example Usage #
```hcl
        
resource "constellix_tags" "tags1" {
  name = "tagsdns"
}

```

## Argument Reference ##
* `name` - (Required) Name of record. Name should be unique.


## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of tag resource.
## Importing ##

An existing Tag can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_tags.example <tag-id>
```

Where tag-id is the Id of tag calculated via Constellix API.