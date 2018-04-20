package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CyCoreSystems/netdiscover/discover"
)

var (
	debug    bool
	provider string
	retField string
)

func init() {
	flag.BoolVar(&debug, "debug", false, `debug mode`)
	flag.StringVar(&provider, "provider", "", `provider type.  Options are: "aws", "azure", "gcp"`)
	flag.StringVar(&retField, "field", "", `return only a single field.  Options are: "hostname", "publicv4", publicv6", "privatev4"`)
}

// Response describes the response from an unlimited discovery
type Response struct {
	// Hostname is the public hostname of the node
	Hostname string `json:"hostname"`

	// PrivateIPv4 is the private (internal) IPv4 address of the node
	PrivateIPv4 string `json:"private_ipv4"`

	// PublicIPv4 is the public (external) IPv4 address of the node
	PublicIPv4 string `json:"public_ipv4"`

	// PublicIPv6 is the public (external) IPv6 address of the node
	PublicIPv6 string `json:"public_ipv6"`
}

func main() {
	if os.Getenv("CLOUD_PROVIDER") != "" {
		provider = os.Getenv("CLOUD_PROVIDER")
	}

	flag.Parse()

	var discoverer discover.Discoverer
	switch provider {
	case "aws":
		discoverer = discover.NewAWSDiscoverer()
	case "azure":
		discoverer = discover.NewAzureDiscoverer()
	case "gcp":
		discoverer = discover.NewGCPDiscoverer()
	default:
		discoverer = discover.NewDiscoverer()
	}

	if retField != "" {
		switch retField {
		case "hostname":
			h, err := discoverer.Hostname()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(h)
			os.Exit(0)
		case "privatev4":
			ip, err := discoverer.PrivateIPv4()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(ip.String())
			os.Exit(0)
		case "publicv4":
			ip, err := discoverer.PublicIPv4()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(ip.String())
			os.Exit(0)
		case "publicv6":
			ip, err := discoverer.PublicIPv6()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(ip.String())
			os.Exit(0)
		default:
			fmt.Println("valid fields are: hostname, privatev4, publicv4, publicv6")
			os.Exit(1)
		}
	}

	var err error
	ret := new(Response)

	ret.Hostname, err = discoverer.Hostname()
	if err != nil && debug {
		log.Println("failed to get hostname:", err)
	}

	privateIP, err := discoverer.PrivateIPv4()
	if err == nil {
		ret.PrivateIPv4 = privateIP.String()
	} else if debug {
		log.Println("failed to get private IPv4 address:", err)
	}

	publicIP, err := discoverer.PublicIPv4()
	if err == nil {
		ret.PublicIPv4 = publicIP.String()
	} else if debug {
		log.Println("failed to get public IPv4 address:", err)
	}

	publicIPv6, err := discoverer.PublicIPv6()
	if err == nil {
		ret.PublicIPv6 = publicIPv6.String()
	} else if debug {
		log.Println("failed to get public IPv6 address:", err)
	}

	enc := json.NewEncoder(os.Stdout)

	if err := enc.Encode(ret); err != nil {
		fmt.Printf("failed to encode response: %+v\n", ret)
		os.Exit(1)
	}
}
