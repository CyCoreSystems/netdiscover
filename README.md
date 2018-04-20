# Netdiscover
[![](https://godoc.org/github.com/CyCoreSystems/netdiscover/discover?status.svg)](http://godoc.org/github.com/CyCoreSystems/netdiscover/discover)

Netdiscover is a Golang-based tool and library by which network information can
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

Currently, this tool supports three cloud providers:

  * `aws`: Amazon Web Services
  * `azure`: Microsoft Azure Cloud
  * `gcp`: Google Cloud Platform

I am happy to accept pull requests to implement more providers.

## Supported Fields

Currently, this tool supports four network data fields:

  * `hostname`: the public hostname of the node
  * `privatev4`: the private (internal) IPv4 address of the node
  * `publicv4`: the public (external) IPv4 address of the node
  * `publicv6`: the public (external) IPv6 address of the node


## Examples

Retrieve the public version 4 IP address of the node instance on GCP:

```
  netdiscover -provider gcp -field publicv4
```

Retrieve all network information on Amazon Wed Services platform:

```
  netdiscover -provider aws
```

Retrieve the version 6 IP address of a baremetal machine:

```
  netdiscover -field publicv6
```

