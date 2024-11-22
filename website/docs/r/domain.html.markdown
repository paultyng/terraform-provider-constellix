---
layout: "constellix"
page_title: "Constellix: constellix_domain"
sidebar_current: "docs-constellix-resource-constellix_domain"
description: |-
    Manages one or more domains within the account.
---
# constellix_domain #
Manages one or more domains within the account.

# Example Usage #
```hcl
resource "constellix_domain" "domain1" {
  name = "domain1.com"
  soa = {
    primary_nameserver = "ns41.constellix.com."
    email              = "dns.constellix.com."
    ttl                = 1800
    refresh            = 48100
    retry              = 7200
    expire             = 1209
    negcache           = 8000
  }
}

```

## Argument Reference ##

* `name` - (Required) Name of the domain.
* `disabled` - (Optional) Indicates if the domain is disabled. The Default value is `false`.
* `has_gtd_regions` - (Optional) GTD Region status of the domain. The Default value is `false`.
* `has_geoip` - (Optional) GTD Region status of the domain. The Default value is `false`.
* `vanity_nameserver` - (Optional) vanity nameserver of domain.
* `nameserver_group` - (Optional) Shows the nameserver group of domain. The Default nameserverGroup is `1`.
* `note` - (Optional) Notes while creating the domain. The maximum length will be 1000 characters.
* `tags` - (Optional) Id of tags applied on domain. The default value is empty.
* `soa` - (Optional) Object.
* `soa.primary_nameserver` - (Optional) The Primary Nameserver is of SOA. 
* `soa.email` - (Optional) An Email Address specifies the mailbox of the person responsible for this zone. 
* `soa.ttl` - (Optional) The number of seconds that this SOA record will be cached in other resolving name servers. 
* `soa.refresh` - (Optional) The time interval (in seconds) before the zone should be refreshed. The recommended value – 86400 (24 Hours). 
* `soa.retry` - (Optional) The time interval (in seconds) before a failed refresh should be retried. Recommended value – 7200 (2 Hours). 
* `soa.expire` - (Optional) The time internal (in seconds) that specifies the upper limit on the time internally that can elapse before the zone is no longer authoritative. This is when the secondary name servers will expire if they are unable to refresh. Recommended value – up to 1209600
* `soa.negcache` - (Optional) The amount of time a record not found is cached. Recommended values can vary, between `180` and `172800` (3 min – 2 days). 

## Attribute Reference ##
This resource exports the following attributes:
* `id` - The constellix calculated id of the domain resource.
* `soa.serial` - The starting serial number for the version of the zone. If the SOA record is applied to a domain that is already created (and thus already has a starting serial number), the existing serial number will be incremented by one. e.g 2015010196

## Importing ##

An existing Record can be [imported][docs-import] into this resource using its Id, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import constellix_domain.example <record-id>
```

Where record-id is the Id of record calculated via Constellix API.