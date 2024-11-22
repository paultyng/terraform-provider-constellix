---
layout: "constellix"
page_title: "Constellix: constellix_http_redirection_record"
sidebar_current: "docs-constellix-resource-constellix_http_redirection_record"
description: |-
     Manages HTTP redirection for a specific domain.
---
# constellix_http_redirection_record #
 Manages HTTP redirection for a specific domain.

# Example Usage #
```hcl
        
resource "constellix_http_redirection_record" "http1" {
  domain_id        = "${constellix_domain.domain1.id}"
  source_type      = "domains"
  name             = "redirectionrecord"
  ttl              = 1800
  redirect_type_id = 1
  url              = "https://www.google.com"
  noanswer         = false
  note             = ""
  gtd_region       = 1
  type             = "HTTPRedirection"
  hardlink_flag    = false
  description      = false
  title            = ""
}


```

## Argument Reference ##
* `ttl` - (Required) TTL must be in between `0` and `2147483647`.
* `url` - (Required) URL link to redirect.
* `redirect_type_id` - (Required) `1` for Standard - 302, `2` for Hidden Frame Masked and `3` for Standard - 301. 
* `source_type` - (Required) `domains` for Domain records and `template` for Template records.
* `name` - (Optional) Name of record. Name should be unique.
* `noanswer` - (Optional) Shows if record is enabled or disabled. Default is `false` (Active).
* `note` - (Optional) Record note.
* `hardlink_flag` - (Optional) Hardlink flag. Default is `false`
* `description` - (Optional) Description.
* `title` - (Optional) Title of redirection.
* `gtd_region` - (Optional) Shows id of GTD region in which record is to be created. note: "gtdRegion" from 2 to 6 will be applied only when GTD region is enabled on domain. 
  * `1` for World (Default). 
  * `2` for Europe. 
  * `3` for US East. 
  * `4` for US West. 
  * `5` for Asia Pacific. 
  * `6` for Oceania.
* `type` - (Optional) Record type `HTTPRedirection`.

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of httpredirection resource.

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_http_redirection.example <source>:<parent-id>:<record-id>
```

Where source can be either domains or templates; parent-id is domain-id or template-id based on the source provided and record-id is the Id of record calculated via Constellix API.