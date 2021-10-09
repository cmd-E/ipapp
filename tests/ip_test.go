package tests

import (
	"ipapp/packages/ip"
	"testing"
)

func TestGetMinIPInNetwork(t *testing.T) {
	tests := []struct {
		networkIP  string
		mask       int
		expectedIP string
	}{
		{"179.178.254.117", 3, "160.0.0.1"},
		{"223.110.137.184", 5, "216.0.0.1"},
		{"64.230.226.51", 9, "64.128.0.1"},
		{"71.85.245.180", 21, "71.85.240.1"},
		{"192.45.92.62", 6, "192.0.0.1"},
		{"184.59.34.204", 2, "128.0.0.1"},
		{"68.159.130.172", 0, "0.0.0.1"},
		{"208.253.84.188", 3, "192.0.0.1"},
		{"199.130.100.242", 11, "199.128.0.1"},
		{"151.205.197.104", 32, "151.205.197.104"},
	}
	for i, test := range tests {
		testedIP := ip.ParseIPFromDecimalString(test.networkIP, test.mask)
		if res := testedIP.GetMinIPInNetwork(); res.GetIPInDec() != test.expectedIP {
			t.Errorf("Test: %d. Expected: %s, got: %v", i+1, test.expectedIP, res)
		}
	}
}

func TestGetMaxIPInNetwork(t *testing.T) {
	tests := []struct {
		networkIP  string
		mask       int
		expectedIP string
	}{
		{"179.178.254.117", 3, "191.255.255.254"},
		{"223.110.137.184", 5, "223.255.255.254"},
		{"64.230.226.51", 9, "64.255.255.254"},
		{"71.85.245.180", 21, "71.85.247.254"},
		{"192.45.92.62", 6, "195.255.255.254"},
		{"184.59.34.204", 2, "191.255.255.254"},
		{"68.159.130.172", 0, "255.255.255.254"},
		{"208.253.84.188", 3, "223.255.255.254"},
		{"199.130.100.242", 11, "199.159.255.254"},
		{"151.205.197.104", 32, "151.205.197.104"},
	}
	for i, test := range tests {
		testedIP := ip.ParseIPFromDecimalString(test.networkIP, test.mask)
		if res := testedIP.GetMaxIPInNetwork(); res.GetIPInDec() != test.expectedIP {
			t.Errorf("Test: %d. Expected: %s, got: %v", i+1, test.expectedIP, res)
		}
	}
}

func TestGetBroadcastAddress(t *testing.T) {
	tests := []struct {
		networkIP  string
		mask       int
		expectedIP string
	}{
		{"179.178.254.117", 3, "191.255.255.255"},
		{"223.110.137.184", 5, "223.255.255.255"},
		{"64.230.226.51", 9, "64.255.255.255"},
		{"71.85.245.180", 21, "71.85.247.255"},
		{"192.45.92.62", 6, "195.255.255.255"},
		{"184.59.34.204", 2, "191.255.255.255"},
		{"68.159.130.172", 0, "255.255.255.255"},
		{"208.253.84.188", 3, "223.255.255.255"},
		{"199.130.100.242", 11, "199.159.255.255"},
		{"151.205.197.104", 32, "151.205.197.104"},
	}
	for i, test := range tests {
		testedIP := ip.ParseIPFromDecimalString(test.networkIP, test.mask)
		if res := testedIP.GetBroadcastAddress(); res.GetIPInDec() != test.expectedIP {
			t.Errorf("Test: %d. Expected: %s, got: %v", i+1, test.expectedIP, res)
		}
	}
}
