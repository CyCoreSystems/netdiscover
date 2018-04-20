package discover

// TODO: decide out how, when, and if to return the Floating IP (URLs noted below)

import (
	"net"
)

const (
	doHostnameURL                  = "http://169.254.169.254/metadata/v1/hostname"
	doPrivateIPv4URL               = "http://169.254.169.254/metadata/v1/interfaces/private/0/ipv4/address"
	doPublicIPv4URL                = "http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address"
	doFloatingPublicIPv4EnabledURL = "http://169.254.169.254/metadata/v1/floating_ip/ipv4/active"
	doFloatingPublicIPv4URL        = "http://169.254.169.254/metadata/v1/interfaces/public/0/anchor_ipv4/address"
	doPublicIPv6URL                = "http://169.254.169.254/metadata/v1/interfaces/public/0/ipv6/address"
)

// NewDigitalOceanDiscoverer returns a new Digital Ocean network discoverer
func NewDigitalOceanDiscoverer() Discoverer {
	return NewDiscoverer(
		PublicHostnameDiscovererOption(doHostname),
		PublicIPv4DiscovererOption(doPublicIPv4),
		PublicIPv6DiscovererOption(doPublicIPv6),
	)

}

func doHostname() (string, error) {
	return StandardHostnameFromHTTP(doHostnameURL, nil)
}

// FIXME:  this URL never seems to be populated
func doPrivateIPv4() (net.IP, error) {
	return StandardIPFromHTTP(doPrivateIPv4URL, nil)
}

func doPublicIPv4() (net.IP, error) {
	return StandardIPFromHTTP(doPublicIPv4URL, nil)
}

func doPublicIPv6() (net.IP, error) {
	return StandardIPFromHTTP(doPublicIPv6URL, nil)
}
