// Package ip implements validated IPv4 subnet arithmetic without I/O.
package ip

import (
	"encoding/binary"
	"fmt"
	"net/netip"
	"strings"
)

// IP holds an IPv4 address and its prefix length. Its zero value is 0.0.0.0/0.
// Private fields prevent octets and prefixes from becoming invalid after parsing.
type IP struct {
	address uint32
	prefix  uint8
}

// IsValid accepts strict, dotted-decimal IPv4, without signs or leading zeroes.
func IsValid(text string) bool {
	_, err := ParseIPFromDecimalString(text, 0)
	return err == nil
}

// ParseIPFromDecimalString validates both the address and the prefix length.
func ParseIPFromDecimalString(text string, prefix int) (IP, error) {
	if err := validatePrefix(prefix); err != nil {
		return IP{}, err
	}
	address, err := netip.ParseAddr(text)
	if err != nil || !address.Is4() {
		return IP{}, fmt.Errorf("invalid IPv4 address %q", text)
	}
	octets := address.As4()
	return IP{address: binary.BigEndian.Uint32(octets[:]), prefix: uint8(prefix)}, nil
}

// ParseIPFromBinString accepts four eight-bit octets. The returned prefix is /0,
// as in the original binary parser; subnet methods preserve their input prefix.
func ParseIPFromBinString(text string) (IP, error) {
	octets := strings.Split(text, ".")
	if len(octets) != 4 {
		return IP{}, fmt.Errorf("invalid binary IPv4 address %q", text)
	}
	var value uint32
	for _, octet := range octets {
		if len(octet) != 8 {
			return IP{}, fmt.Errorf("invalid binary IPv4 octet %q", octet)
		}
		for _, bit := range octet {
			if bit != '0' && bit != '1' {
				return IP{}, fmt.Errorf("invalid binary IPv4 octet %q", octet)
			}
			value = value<<1 | uint32(bit-'0')
		}
	}
	return IP{address: value}, nil
}

func validatePrefix(prefix int) error {
	if prefix < 0 || prefix > 32 {
		return fmt.Errorf("prefix length must be between 0 and 32, got %d", prefix)
	}
	return nil
}

func maskBits(prefix uint8) uint32 { return ^uint32(0) << (32 - uint(prefix)) }

// ConvertMask returns the dotted IPv4 mask for a validated prefix length.
func ConvertMask(prefix int) (IP, error) {
	if err := validatePrefix(prefix); err != nil {
		return IP{}, err
	}
	return IP{address: maskBits(uint8(prefix)), prefix: uint8(prefix)}, nil
}

// PrefixLength returns the address's prefix length.
func (address IP) PrefixLength() int { return int(address.prefix) }

func (address IP) GetNetworkPart() IP {
	return IP{address: address.address & maskBits(address.prefix), prefix: address.prefix}
}

// GetMinIPInNetwork returns the first conventional host address, or the first
// endpoint for a /31 point-to-point link or the sole /32 address.
func (address IP) GetMinIPInNetwork() IP {
	first := address.GetNetworkPart()
	if address.prefix < 31 {
		first.address++
	}
	return first
}

func (address IP) GetMaxIPInNetwork() IP {
	last := address.GetBroadcastAddress()
	if address.prefix < 31 {
		last.address--
	}
	return last
}

// GetBroadcastAddress returns the inclusive upper subnet boundary. /31 and /32
// have no directed broadcast: their upper boundary is itself an endpoint.
func (address IP) GetBroadcastAddress() IP {
	return IP{address: address.address | ^maskBits(address.prefix), prefix: address.prefix}
}

// HostCount follows conventional IPv4 subnet rules and RFC 3021 for /31 links.
// It describes subnet arithmetic, not address allocation or reachability policy.
func (address IP) HostCount() uint64 {
	count := uint64(1) << (32 - address.prefix)
	if address.prefix < 31 {
		count -= 2
	}
	return count
}

func (address IP) GetIPInBin() string {
	return fmt.Sprintf("%08b.%08b.%08b.%08b", byte(address.address>>24),
		byte(address.address>>16), byte(address.address>>8), byte(address.address))
}

func (address IP) GetIPInDec() string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(address.address>>24),
		byte(address.address>>16), byte(address.address>>8), byte(address.address))
}

type IINArray []IpsInNetwork

type IpsInNetwork struct {
	Network IP
	IPs     []IP
}

// SortIpsByNetworks groups addresses in expected linear time. Despite the legacy
// name, groups retain first-seen order, and addresses retain their original order.
func SortIpsByNetworks(addresses []IP) IINArray {
	var groups IINArray
	positions := make(map[IP]int)
	for _, address := range addresses {
		network := address.GetNetworkPart()
		if position, exists := positions[network]; exists {
			groups[position].IPs = append(groups[position].IPs, address)
		} else {
			positions[network] = len(groups)
			groups = append(groups, IpsInNetwork{Network: network, IPs: []IP{address}})
		}
	}
	return groups
}
