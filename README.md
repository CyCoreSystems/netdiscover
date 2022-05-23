# Netdiscover
[![Build Status](https://travis-ci.org/CyCoreSystems/netdiscover.png)](https://travis-ci.org/CyCoreSystems/netdiscover) [![](https://godoc.org/github.com/CyCoreSystems/netdiscover/discover?status.svg)](http://godoc.org/github.com/CyCoreSystems/netdiscover/discover)

Netdiscover is a CLI tool and Golang library by which network information may
be discovered on various cloud platforms and bare metal installations.  The
typical use case is, when running inside Kubernetes or a container, to discover
the public IP and/or hostname of the node on which the container is running.
This is commonly necessary to configure VoIP applications.

## CLI tool

In the root directory can be found a CLI tool which can be used to query network
information.  There are two options:

  * `-provider <name>`:  if a cloud provider is specified, that provider's
    metadata services will be used to determine the network information.
    Otherwise, a best-effort approach will be used.
  * `-field <name>`: if a field is specified, only that particular network
    detail will be returned.  Otherwise, a JSON object will be returned
    containing all network information which was discovered.

## Supported providers

Currently, this tool supports four cloud providers:

  * `aws`: Amazon Web Services
  * `azure`: Microsoft Azure Cloud
  * `do`: Digital Ocean
  * `gcp`: Google Cloud Platform
  * `aliyun`: Alibaba Cloud
  * ``: general discovery (baremetal)

Not all providers support all network details.  General discovery should be used
for baremetal or other unsupported environments.  It will use the public service
(jsonip.io) to determine the necessary network data.

I am happy to accept pull requests to implement more providers.

## Supported Fields

Currently, this tool supports four network data fields:

  * `hostname`: the public hostname of the node
  * `privatev4`: the private (internal) IPv4 address of the node
  * `publicv4`: the public (external) IPv4 address of the node
  * `publicv6`: the public (external) IPv6 address of the node

Note that for DigitalOcean\'s hostname feature to work as expected, you need to
change the Droplet name to the fully-qualified domain name of the host.  Doing
so will also cause DigitalOcean to register the reverse DNS lookup for the
Droplet\'s IP address, so this should generally be done, anyway.


## Examples

Retrieve the public version 4 IP address of the node instance on GCP:

```
  netdiscover -provider gcp -field publicv4
```

Retrieve all network information on Amazon Web Services platform:

```
  netdiscover -provider aws
```

Retrieve the version 6 IP address of a baremetal machine:

```
  netdiscover -field publicv6
```

