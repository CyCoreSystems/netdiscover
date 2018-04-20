package discover

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Discoverer describes an interface by which network details may be discovered
type Discoverer interface {

	// Hostname returns the public hostname.
	Hostname() (string, error)

	// PrivateIPv4 returns the private IPv4 address.
	PrivateIPv4() (net.IP, error)

	// PublicIPv4 returns the public IPv4 address.
	PublicIPv4() (net.IP, error)

	// PublicIPv6 returns the public IPv6 address.
	PublicIPv6() (net.IP, error)
}

// GenericDiscoverer provides a flexible network detail discoverer, allowing
// one to supply URLs and metadata for specific providers.  Any value which is
// not set will be attempted by the default process.
type genericDiscoverer struct {
	opts *DiscovererOptions
}

// DiscovererOptions describes the options used to construct a generic discoverer.
type DiscovererOptions struct {
	getHostname    func() (string, error)
	getPrivateIPv4 func() (net.IP, error)
	getPublicIPv4  func() (net.IP, error)
	getPublicIPv6  func() (net.IP, error)
}

// DiscovererOption is a function which configures a GenericDiscoverer.
type DiscovererOption func(o *DiscovererOptions)

// PrivateIPv4DiscovererOption retuns a DiscovererOption which sets the Private
// IPv4 address discovery function.
func PrivateIPv4DiscovererOption(f func() (net.IP, error)) func(*DiscovererOptions) {
	return func(o *DiscovererOptions) {
		if f == nil {
			return
		}
		o.getPrivateIPv4 = f
	}
}

// PublicIPv4DiscovererOption returns a DiscovererOption which sets the Public
// IPv4 address discovery function.
func PublicIPv4DiscovererOption(f func() (net.IP, error)) func(*DiscovererOptions) {
	return func(o *DiscovererOptions) {
		if f == nil {
			return
		}
		o.getPublicIPv4 = f
	}
}

// PublicIPv6DiscovererOption returns a DiscovererOption which sets the Public
// IPv6 address discovery function.
func PublicIPv6DiscovererOption(f func() (net.IP, error)) func(*DiscovererOptions) {
	return func(o *DiscovererOptions) {
		if f == nil {
			return
		}
		o.getPublicIPv6 = f
	}
}

// PublicHostnameDiscovererOption returns a DiscovererOption which sets the Public
// Hostname discovery function.
func PublicHostnameDiscovererOption(f func() (string, error)) func(*DiscovererOptions) {
	return func(o *DiscovererOptions) {
		if f == nil {
			return
		}
		o.getHostname = f
	}
}

// NewDiscoverer creates a new network discoverer
func NewDiscoverer(opts ...DiscovererOption) Discoverer {
	o := &DiscovererOptions{
		getPrivateIPv4: defaultPrivateIPv4,
		getPublicIPv4:  defaultPublicIPv4,
		getPublicIPv6:  defaultPublicIPv6,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}

	if o.getHostname == nil {
		o.getHostname = func() (string, error) {
			return defaultHostname(o.getPublicIPv4)
		}
	}

	return &genericDiscoverer{opts: o}
}

func (d *genericDiscoverer) Hostname() (string, error) {
	if d == nil || d.opts == nil {
		return "", errors.New("discoverer not created")
	}
	if d.opts.getHostname == nil {
		return "", errors.New("no hostname discoverer")
	}
	return d.opts.getHostname()
}

func (d *genericDiscoverer) PrivateIPv4() (net.IP, error) {
	if d == nil || d.opts == nil {
		return nil, errors.New("discoverer not created")
	}
	if d.opts.getPrivateIPv4 == nil {
		return nil, errors.New("no private IPv4 discoverer")
	}
	return d.opts.getPrivateIPv4()
}

func (d *genericDiscoverer) PublicIPv4() (net.IP, error) {
	if d == nil || d.opts == nil {
		return nil, errors.New("discoverer not created")
	}
	if d.opts.getPublicIPv4 == nil {
		return nil, errors.New("no public IPv4 discoverer")
	}
	return d.opts.getPublicIPv4()
}

func (d *genericDiscoverer) PublicIPv6() (net.IP, error) {
	if d == nil || d.opts == nil {
		return nil, errors.New("discoverer not created")
	}
	if d.opts.getPublicIPv6 == nil {
		return nil, errors.New("no public IPv6 discoverer")
	}
	return d.opts.getPublicIPv6()
}

func defaultPublicIPv4() (net.IP, error) {
	resp, err := http.Get("http://ipv4.jsonip.io")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	type response struct {
		Address string `json:"address"`
	}
	data := new(response)

	err = dec.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	ip := net.ParseIP(data.Address)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse address: %s", data.Address)
	}

	return ip, nil
}

func defaultPublicIPv6() (net.IP, error) {
	resp, err := http.Get("http://ipv6.jsonip.io")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	type response struct {
		Address string `json:"address"`
	}
	data := new(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	err = dec.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	ip := net.ParseIP(data.Address)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse address: %s", data.Address)
	}

	return ip, nil
}

func defaultPrivateIPv4() (net.IP, error) {
	netifs, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get list of local interfaces: %v", err)
	}

	for _, i := range netifs {
		addresses, err := i.Addrs()
		if err != nil {
			continue // next interface
		}

		for _, addrIf := range addresses {
			a := addrIf.String()

			tmpIP, _, err := net.ParseCIDR(a)
			if err != nil || tmpIP == nil {
				continue // next address
			}

			if tmpIP.IsLoopback() {
				break // next interface
			}

			if !tmpIP.IsGlobalUnicast() {
				// multicast, link-local, etc
				continue // next address
			}

			if toIPv4 := tmpIP.To4(); toIPv4 == nil {
				// Not an IPv4 address
				continue // next address
			}

			return tmpIP, nil
		}
	}

	return nil, errors.New("valid address not found")
}

func defaultHostname(ipFunc func() (net.IP, error)) (string, error) {

	if ipFunc == nil {
		return "", errors.New("no public IP discovery function")
	}

	ip, err := ipFunc()
	if err != nil {
		return "", fmt.Errorf("failed to obtain public IP: %v", err)
	}

	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return "", fmt.Errorf("failed to reverse-lookup ip address: %v", err)
	}

	for _, name := range names {
		if len(name) < 6 {
			// implausibly short name
			continue
		}
		if !strings.Contains(name, ".") {
			// implausible TLD or local-only hostname
			continue
		}
		if strings.HasSuffix(name, ".local") {
			continue
		}

		return strings.TrimSuffix(name, "."), err
	}

	return "", fmt.Errorf("failed to discover valid public name")
}
