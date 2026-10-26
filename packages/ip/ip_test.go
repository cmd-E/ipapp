package ip_test

import (
	"encoding/binary"
	"fmt"
	"ipapp/packages/ip"
	"math/rand"
	"net"
	"reflect"
	"testing"
)

func mustParse(t *testing.T, text string, prefix int) ip.IP {
	t.Helper()
	address, err := ip.ParseIPFromDecimalString(text, prefix)
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func TestStrictIPv4Parsing(t *testing.T) {
	for _, text := range []string{"0.0.0.0", "255.255.255.255", "192.0.2.1", "127.0.0.1"} {
		t.Run(text, func(t *testing.T) {
			address := mustParse(t, text, 24)
			if !ip.IsValid(text) || address.GetIPInDec() != text || address.PrefixLength() != 24 {
				t.Fatal(address)
			}
		})
	}
	for _, text := range []string{"", "1", "1.2.3", "1.2.3.4.5", "1..2.3", "256.1.2.3",
		"-1.2.3.4", "+1.2.3.4", "01.2.3.4", "1.002.3.4", "1.2.3.4 ", " 1.2.3.4",
		"1.2.3.4/24", "1.2.3.4\n", "::1", "::ffff:192.0.2.1", "1.2.3.4%zone", "999999999999999999999.2.3.4"} {
		t.Run(fmt.Sprintf("invalid-%q", text), func(t *testing.T) {
			if _, err := ip.ParseIPFromDecimalString(text, 24); err == nil || ip.IsValid(text) {
				t.Fatalf("accepted %q", text)
			}
		})
	}
}

func TestPrefixValidationAndMaskRepresentation(t *testing.T) {
	maximumInt := int(^uint(0) >> 1)
	for _, prefix := range []int{-1, 33, maximumInt, -maximumInt - 1} {
		if _, err := ip.ParseIPFromDecimalString("192.0.2.1", prefix); err == nil {
			t.Fatalf("accepted prefix %d", prefix)
		}
		if _, err := ip.ConvertMask(prefix); err == nil {
			t.Fatalf("accepted mask %d", prefix)
		}
	}
	for _, item := range []struct {
		prefix int
		mask   string
	}{
		{0, "0.0.0.0"}, {1, "128.0.0.0"}, {8, "255.0.0.0"}, {16, "255.255.0.0"},
		{24, "255.255.255.0"}, {30, "255.255.255.252"}, {31, "255.255.255.254"}, {32, "255.255.255.255"},
	} {
		mask, err := ip.ConvertMask(item.prefix)
		if err != nil || mask.GetIPInDec() != item.mask || mask.PrefixLength() != item.prefix {
			t.Fatal(item, mask, err)
		}
	}
}

func TestBinaryParsingAndFormatting(t *testing.T) {
	for _, item := range []struct{ decimal, binary string }{
		{"0.0.0.0", "00000000.00000000.00000000.00000000"},
		{"255.255.255.255", "11111111.11111111.11111111.11111111"},
		{"192.0.2.129", "11000000.00000000.00000010.10000001"},
	} {
		address := mustParse(t, item.decimal, 24)
		parsed, err := ip.ParseIPFromBinString(item.binary)
		if err != nil || parsed.GetIPInDec() != item.decimal || address.GetIPInBin() != item.binary || parsed.PrefixLength() != 0 {
			t.Fatal(parsed, err)
		}
	}
	for _, text := range []string{"", "00000000", "0.0.0.0", "00000000.00000000.00000000.000000000",
		"00000000.00000000.00000000.00000002", "00000000.00000000.00000000.+0000000",
		"00000000.00000000.00000000.00000000.00000000", "00000000.00000000.00000000.0000000é"} {
		if _, err := ip.ParseIPFromBinString(text); err == nil {
			t.Fatalf("accepted %q", text)
		}
	}
}

func TestHostRangeGoldenCases(t *testing.T) {
	for _, item := range []struct {
		address                         string
		prefix                          int
		network, minimum, maximum, last string
		hosts                           uint64
	}{
		{"192.0.2.129", 0, "0.0.0.0", "0.0.0.1", "255.255.255.254", "255.255.255.255", 4294967294},
		{"192.0.2.129", 8, "192.0.0.0", "192.0.0.1", "192.255.255.254", "192.255.255.255", 16777214},
		{"192.0.2.129", 24, "192.0.2.0", "192.0.2.1", "192.0.2.254", "192.0.2.255", 254},
		{"192.0.2.129", 25, "192.0.2.128", "192.0.2.129", "192.0.2.254", "192.0.2.255", 126},
		{"192.0.2.129", 30, "192.0.2.128", "192.0.2.129", "192.0.2.130", "192.0.2.131", 2},
		{"192.0.2.129", 31, "192.0.2.128", "192.0.2.128", "192.0.2.129", "192.0.2.129", 2},
		{"192.0.2.129", 32, "192.0.2.129", "192.0.2.129", "192.0.2.129", "192.0.2.129", 1},
		{"255.255.255.255", 32, "255.255.255.255", "255.255.255.255", "255.255.255.255", "255.255.255.255", 1},
	} {
		t.Run(fmt.Sprintf("%s/%d", item.address, item.prefix), func(t *testing.T) {
			address := mustParse(t, item.address, item.prefix)
			derived := []ip.IP{address.GetNetworkPart(), address.GetMinIPInNetwork(), address.GetMaxIPInNetwork(), address.GetBroadcastAddress()}
			for i, want := range []string{item.network, item.minimum, item.maximum, item.last} {
				if derived[i].GetIPInDec() != want || derived[i].PrefixLength() != item.prefix {
					t.Fatalf("derived %d: %s/%d; want %s/%d", i, derived[i].GetIPInDec(), derived[i].PrefixLength(), want, item.prefix)
				}
				if derived[i].GetNetworkPart() != address.GetNetworkPart() {
					t.Fatal("derived prefix lost")
				}
			}
			if address.HostCount() != item.hosts {
				t.Fatal(address.HostCount())
			}
		})
	}
}

func TestEveryPrefixAgainstStandardLibrary(t *testing.T) {
	random := rand.New(rand.NewSource(3021))
	addresses := []net.IP{net.IPv4(0, 0, 0, 0).To4(), net.IPv4(255, 255, 255, 255).To4()}
	for i := 0; i < 128; i++ {
		bytes := make(net.IP, 4)
		binary.BigEndian.PutUint32(bytes, random.Uint32())
		addresses = append(addresses, bytes)
	}
	for prefix := 0; prefix <= 32; prefix++ {
		t.Run(fmt.Sprintf("/%d", prefix), func(t *testing.T) {
			for _, original := range addresses {
				address := mustParse(t, original.String(), prefix)
				network := net.IPNet{IP: original.Mask(net.CIDRMask(prefix, 32)), Mask: net.CIDRMask(prefix, 32)}
				if address.GetNetworkPart().GetIPInDec() != network.IP.String() {
					t.Fatal("network disagrees with net.IP.Mask")
				}
				first, last := net.ParseIP(address.GetMinIPInNetwork().GetIPInDec()), net.ParseIP(address.GetMaxIPInNetwork().GetIPInDec())
				if !network.Contains(first) || !network.Contains(last) || address.GetNetworkPart().GetNetworkPart() != address.GetNetworkPart() {
					t.Fatal("out-of-network or non-idempotent result")
				}
				if binary.BigEndian.Uint32(last.To4())-binary.BigEndian.Uint32(first.To4())+1 != uint32(address.HostCount()) {
					t.Fatal("range and host count disagree")
				}
			}
		})
	}
}

func TestGroupingOrderPrefixIdentityAndInputOwnership(t *testing.T) {
	input := []ip.IP{mustParse(t, "192.0.2.2", 24), mustParse(t, "10.0.0.1", 8), mustParse(t, "192.0.2.1", 24), mustParse(t, "192.0.2.2", 24), mustParse(t, "192.0.2.1", 32)}
	before := append([]ip.IP(nil), input...)
	groups := ip.SortIpsByNetworks(input)
	if len(groups) != 3 || groups[0].Network.GetIPInDec() != "192.0.2.0" || groups[1].Network.GetIPInDec() != "10.0.0.0" || groups[2].Network.PrefixLength() != 32 {
		t.Fatal(groups)
	}
	if !reflect.DeepEqual(groups[0].IPs, []ip.IP{input[0], input[2], input[3]}) {
		t.Fatal("order or duplicates lost")
	}
	groups[0].IPs[0] = ip.IP{}
	if !reflect.DeepEqual(input, before) {
		t.Fatal("group storage aliases input")
	}
	if len(ip.SortIpsByNetworks(nil)) != 0 {
		t.Fatal("empty grouping")
	}
	overlap := ip.SortIpsByNetworks([]ip.IP{mustParse(t, "0.0.0.0", 0), mustParse(t, "0.0.0.0", 8)})
	if len(overlap) != 2 {
		t.Fatal("distinct prefixes merged")
	}
}

func TestZeroValue(t *testing.T) {
	var address ip.IP
	if address.GetIPInDec() != "0.0.0.0" || address.GetNetworkPart() != address || address.PrefixLength() != 0 || address.HostCount() != 4294967294 {
		t.Fatal(address)
	}
}

func FuzzDecimalParsing(f *testing.F) {
	for _, text := range []string{"192.0.2.1", "255.255.255.255", "", "::1", "+1.2.3.4"} {
		f.Add(text, 24)
	}
	f.Add("192.0.2.1", -1)
	f.Add("192.0.2.1", 32)
	f.Fuzz(func(t *testing.T, text string, prefix int) {
		address, err := ip.ParseIPFromDecimalString(text, prefix)
		if err != nil {
			return
		}
		if prefix < 0 || prefix > 32 || address.PrefixLength() != prefix || !ip.IsValid(address.GetIPInDec()) {
			t.Fatal("invalid parsed state")
		}
		parsed, err := ip.ParseIPFromBinString(address.GetIPInBin())
		if err != nil || parsed.GetIPInDec() != address.GetIPInDec() {
			t.Fatal("binary round trip", err)
		}
		if address.GetMinIPInNetwork().GetNetworkPart() != address.GetNetworkPart() || address.GetMaxIPInNetwork().GetNetworkPart() != address.GetNetworkPart() {
			t.Fatal("derived address lost prefix")
		}
	})
}

func FuzzBinaryParsing(f *testing.F) {
	for _, text := range []string{"00000000.00000000.00000000.00000000", "11111111.11111111.11111111.11111111", "0.0.0.0", "", "garbage"} {
		f.Add(text)
	}
	f.Fuzz(func(t *testing.T, text string) {
		address, err := ip.ParseIPFromBinString(text)
		if err == nil && address.GetIPInBin() != text {
			t.Fatalf("accepted noncanonical binary %q", text)
		}
	})
}
