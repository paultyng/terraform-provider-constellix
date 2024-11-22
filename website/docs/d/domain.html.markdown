---
layout: "constellix"
page_title: "Constellix: constellix_domain"
sidebar_current: "docs-constellix-data-source-constellix_domain"
description: |-
  Data source for Constellix domain
---

# constellix_domain #
Data source for Constellix domain

## Example Usage ##

```hcl
data "constellix_domain" "domain1" {
  name = "datasourcedomain.com"
}
 
```
## Argument Reference ##

* `name` - (Required) Name of the domain.

## Attribute Reference

* `disabled` - (Optional) Indicates if the domain is disabled. The Default value is false.
* `has_gtd_regions` - (Optional) GTD Region status of the domain. The Default value is false.
* `has_geoip` - (Optional) GTD Region status of the domain. The Default value is false.
* `vanity_nameserver` - (Optional) vanity nameserver of domain.
* `nameserver_group` - (Optional) Shows the nameserver group of domain. The Default nameserverGroup is 1.
* `note` - (Optional) Notes while creating the domain. The maximum length will be 1000 characters.
* `tags` - (Optional) Id of tags applied on domain. The default value is empty.
* `soa` - (Optional) Object
* `soa.primary_nameserver` - (Optional) The Default value of SOA Primary Nameserver is "ns0.constellix.com.". However, it is possible to create a custom SOA record with differing values if required.
* `soa.email` - (Optional) An Email Address specifies the mailbox of the person responsible for this zone. The default value is "dns.constellix.com."
* `soa.ttl` - (Optional) The number of seconds that this SOA record will be cached in other resolving name servers. The Default value is "86400".TTL must be in between 0 and 2147483647
* `soa.refresh` - (Optional) The time interval (in seconds) before the zone should be refreshed. The recommended value – 86400 (24 Hours). The default value is 43200 (12 hours)
* `soa.serial` - (Optional) The starting serial number for the version of the zone. If the SOA record is applied to a domain that is already created (and thus already has a starting serial number), the existing serial number will be incremented by one. e.g 2015010196
* `soa.retry` - (Optional) The time interval (in seconds) before a failed refresh should be retried. Recommended value – 7200 (2 Hours). The default value is 1 hour
* `soa.expire` - (Optional) The time internal (in seconds) that specifies the upper limit on the time internally that can elapse before the zone is no longer authoritative. This is when the secondary name servers will expire if they are unable to refresh. Recommended value – up to 1209600
* `soa.negcache` - (Optional) The amount of time a record not found is cached. Recommended values can vary, between 180 and 172800 (3 min – 2 days). The default value is 180