package discover

import (
	"net"
)

const (
	azurePrivateIPv4URL = "http://169.254.169.254/metadata/instance/network/interface/0/ipv4/ipAddress/0/privateIpAddress?api-version=2017-08-01&format=text"
	azurePublicIPv4URL  = "http://169.254.169.254/metadata/instance/network/interface/0/ipv4/ipAddress/0/publicIpAddress?api-version=2017-08-01&format=text"
)

// NewAzureDiscoverer returns a new Google Cloud Platform network discoverer
func NewAzureDiscoverer() Discoverer {
	return NewDiscoverer(
		PrivateIPv4DiscovererOption(azurePrivateIPv4),
		PublicIPv4DiscovererOption(azurePublicIPv4),
	)
}

func azurePrivateIPv4() (net.IP, error) {
	return StandardIPFromHTTP(azurePrivateIPv4URL, map[string]string{"Metadata": "true"})
}

func azurePublicIPv4() (net.IP, error) {
	return StandardIPFromHTTP(azurePublicIPv4URL, map[string]string{"Metadata": "true"})
}
