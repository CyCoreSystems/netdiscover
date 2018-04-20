package discover

import (
	"net"
)

const (
	awsPrivateIPv4URL = "http://169.254.169.254/latest/meta-data/local-ipv4"
	awsPublicIPv4URL  = "http://169.254.169.254/latest/meta-data/public-ipv4"
	awsHostnameURL    = "http://169.254.169.254/latest/meta-data/public-hostname"
)

// NewAWSDiscoverer returns a new Amazon Web Services network discoverer
func NewAWSDiscoverer() Discoverer {
	return NewDiscoverer(
		PrivateIPv4DiscovererOption(awsPrivateIPv4),
		PublicIPv4DiscovererOption(awsPublicIPv4),
		PublicHostnameDiscovererOption(awsHostname),
	)

}

func awsPrivateIPv4() (net.IP, error) {
	return StandardIPFromHTTP(awsPrivateIPv4URL, nil)
}

func awsPublicIPv4() (net.IP, error) {
	return StandardIPFromHTTP(awsPublicIPv4URL, nil)
}

func awsHostname() (string, error) {
	return StandardHostnameFromHTTP(awsHostnameURL, nil)
}
