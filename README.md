# Constellix DNS Provider

## Disclaimer
This is a fork of Constellix/terraform-provider-constellix with a small change to env variables and go releaser settings to enable builds on m1 macs.

### Environemnt variables
```
CONSTELLIX_API_KEY
CONSTELLIX_SECRET_KEY
```

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) Latest Version

- [Go](https://golang.org/doc/install) go1.17

## Building The Provider ##
Clone this repository to: `$GOPATH/src/github.com/Constellix/terraform-provider-constellix`.

```sh
$ mkdir -p $GOPATH/src/github.com/Constellix; cd $GOPATH/src/github.com/Constellix
$ git clone https://github.com/Constellix/terraform-provider-constellix.git
```

Enter the provider directory and run make build to build the provider binary.

```sh
$ cd $GOPATH/src/github.com/Constellix/terraform-provider-constellix
$ make build

```

Using The Provider
------------------
If you are building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory, run `terraform init` to initialize it.

ex.
```hcl
#configure provider with your Constellix  credentials.
provider "constellix" {
  # constellix Api key
  apikey = "apikey"
  # cosntellix secret key
  secretkey = "secretkey"
  insecure = true
  proxy_url = "https://proxy_server:proxy_port"
}

resource "constellix_domain" "domain1" {
  name = "domain1.com"
  soa = {
    primary_nameserver = "ns41.constellix.com."
    ttl                = 1800
    refresh            = 48100
    retry              = 7200
    expire             = 1209
    negcache           = 8000
  }
}


```
Note : If you are facing the issue of `409 conflict error` try running your Terraform configuration with parallelism set to one as mentioned below.

```
terraform plan -parallelism=1
terraform apply -parallelism=1
```  

Developing The Provider
-----------------------
If you want to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine. You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider with sanity checks present in scripts directory and put the provider binary in `$GOPATH/bin` directory.

# Importing A records

`terraform import "module.dev-hz1[0].module.something.constellix_a_record.something" domains:1035333:32661234`
Where domains is static, 1035333 is the domain id and 32661234 is the record id.
