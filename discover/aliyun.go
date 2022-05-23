package discover

import (
	"net"
)

const (
	aliPrivateIPv4URL = "http://100.100.100.200/latest/meta-data/private-ipv4"
	aliPublicIPv4URL  = "http://100.100.100.200/latest/meta-data/eipv4"
	aliHostnameURL    = "http://100.100.100.200/latest/meta-data/hostname"
)

// NewAliYunDiscoverer returns a new Alibaba Cloud Services network discoverer
func NewAliYunDiscoverer() Discoverer {
	return NewDiscoverer(
		PrivateIPv4DiscovererOption(aliPrivateIPv4),
		PublicIPv4DiscovererOption(aliPublicIPv4),
		PublicHostnameDiscovererOption(aliHostname),
	)

}

func aliPrivateIPv4() (net.IP, error) {
	return StandardIPFromHTTP(aliPrivateIPv4URL, nil)
}

func aliPublicIPv4() (net.IP, error) {
	return StandardIPFromHTTP(aliPublicIPv4URL, nil)
}

func aliHostname() (string, error) {
	return StandardHostnameFromHTTP(aliHostnameURL, nil)
}
