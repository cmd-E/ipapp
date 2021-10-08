package ip

import (
	"fmt"
	"strconv"
	"strings"
)

var numOfDots = 0 // compensate indexes

// IP Represents IPv4 address
type IP struct {
	Part1   int
	Part2   int
	Part3   int
	Part4   int
	MaskNum int
}

// IsValid checks if adresses is ipv4
func IsValid(ip string) bool {
	splittedIP := strings.Split(ip, ".")
	if len(splittedIP) != 4 {
		return false
	}
	for _, v := range splittedIP {
		ipPart, err := strconv.Atoi(v)
		if err != nil {
			return false
		}
		if ipPart < 0 || ipPart > 255 {
			return false
		}
	}
	return true
}

// ParseIPFromDecimalString parses ip from string ip in decimal format
func ParseIPFromDecimalString(ipStr string, maskNum int) IP {
	var ip IP
	splittedIP := strings.Split(ipStr, ".")
	ip.Part1, _ = strconv.Atoi(splittedIP[0])
	ip.Part2, _ = strconv.Atoi(splittedIP[1])
	ip.Part3, _ = strconv.Atoi(splittedIP[2])
	ip.Part4, _ = strconv.Atoi(splittedIP[3])
	ip.MaskNum = maskNum
	return ip
}

// ParseIPFromBinString parses ip from string ip in binary format
func ParseIPFromBinString(ipStr string) IP {
	var ip IP
	splittedIP := strings.Split(ipStr, ".")
	t1, _ := strconv.ParseInt(splittedIP[0], 2, 32)
	t2, _ := strconv.ParseInt(splittedIP[1], 2, 32)
	t3, _ := strconv.ParseInt(splittedIP[2], 2, 32)
	t4, _ := strconv.ParseInt(splittedIP[3], 2, 32)
	ip.Part1 = int(t1)
	ip.Part2 = int(t2)
	ip.Part3 = int(t3)
	ip.Part4 = int(t4)
	return ip
}

// ConvertMask converts mask from number to binary representation
func ConvertMask(maskNum int) IP {
	var mask IP
	var ipPart string
	for i := 1; i <= 32; i++ {
		if i <= maskNum {
			ipPart += "1"
		} else {
			ipPart += "0"
		}
		if i == 8 || i == 16 || i == 24 {
			ipPart += "."
			numOfDots++
		}
	}
	splittedIP := strings.Split(ipPart, ".")
	t1, _ := strconv.ParseInt(splittedIP[0], 2, 32)
	t2, _ := strconv.ParseInt(splittedIP[1], 2, 32)
	t3, _ := strconv.ParseInt(splittedIP[2], 2, 32)
	t4, _ := strconv.ParseInt(splittedIP[3], 2, 32)
	mask.Part1 = int(t1)
	mask.Part2 = int(t2)
	mask.Part3 = int(t3)
	mask.Part4 = int(t4)
	return mask
}

// GetNetworkPart returns network of ip address
func (ip IP) GetNetworkPart() IP {
	binIP := ip.GetIPInBin()
	binIP = binIP[:ip.MaskNum+numOfDots] // TODO: check for out of range
	for i := ip.MaskNum + 1; i <= 32; i++ {
		binIP += "0"
		if i == 8 || i == 16 || i == 24 {
			binIP += "."
		}
	}
	splittedIP := strings.Split(binIP, ".")
	var networkIP IP
	t1, _ := strconv.ParseInt(splittedIP[0], 2, 32)
	t2, _ := strconv.ParseInt(splittedIP[1], 2, 32)
	t3, _ := strconv.ParseInt(splittedIP[2], 2, 32)
	t4, _ := strconv.ParseInt(splittedIP[3], 2, 32)
	networkIP.Part1 = int(t1)
	networkIP.Part2 = int(t2)
	networkIP.Part3 = int(t3)
	networkIP.Part4 = int(t4)
	networkIP.MaskNum = ip.MaskNum
	return networkIP
}

// GetMinIPInNetwork returns minimum address in network
func (ip IP) GetMinIPInNetwork() IP {
	ipSplit := []rune(ip.GetIPInBin())
	for i := ip.MaskNum + numOfDots; i < 32+numOfDots; i++ {
		if ipSplit[i] == '.' {
			i--
			continue
		}
		ipSplit[i] = '0'
	}
	ipSplit[len(ipSplit)-1] = '1'
	return ParseIPFromBinString(string(ipSplit))
}

// GetMaxIPInNetwork returns maximum address in network
func (ip IP) GetMaxIPInNetwork() IP {
	ipSplit := []rune(ip.GetIPInBin())
	for i := ip.MaskNum + numOfDots; i < 32+numOfDots; i++ {
		if ipSplit[i] == '.' {
			i--
			continue
		}
		ipSplit[i] = '1'
	}
	ipSplit[len(ipSplit)-1] = '0'
	return ParseIPFromBinString(string(ipSplit))
}

// GetBroadcastAddress returns broadcast address of the network
func (ip IP) GetBroadcastAddress() IP {
	ipSplit := []rune(ip.GetIPInBin())
	for i := ip.MaskNum + numOfDots; i < 32+numOfDots; i++ {
		if ipSplit[i] == '.' {
			i--
			continue
		}
		ipSplit[i] = '1'
	}
	return ParseIPFromBinString(string(ipSplit))
}

// GetIPInBin returns ip in binary representation
func (ip IP) GetIPInBin() string {
	return fmt.Sprintf("%08b.%08b.%08b.%08b", ip.Part1, ip.Part2, ip.Part3, ip.Part4)
}

// GetIPInDec returns ip in decimal representation
func (ip IP) GetIPInDec() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip.Part1, ip.Part2, ip.Part3, ip.Part4)
}

// func (ip IP) AreInTheSameNetwork(ipToCompare IP) bool {
// 	return ip.GetNetworkPart() == ipToCompare.GetNetworkPart()
// }

// IINArray Array of IpsInNetwork structs
type IINArray []IpsInNetwork

func (arr IINArray) networkIPIsInTheList(networkIPToCheck IP) bool {
	for _, v := range arr {
		if v.Network == networkIPToCheck {
			return true
		}
	}
	return false
}

func (arr *IINArray) insertToExistingNetworkIP(networkIP, ipToInsert IP) {
	for i, v := range *arr {
		if v.Network == networkIP {
			v.IPs = append(v.IPs, ipToInsert)
			(*arr)[i].IPs = v.IPs
		}
	}
}

// IpsInNetwork stores ips which are in the same network
type IpsInNetwork struct {
	Network IP
	IPs     []IP
}

// SortIpsByNetworks sorts ips by it's network parts
func SortIpsByNetworks(ips []IP) IINArray {
	var networksIps IINArray
	for _, ip := range ips {
		if networksIps.networkIPIsInTheList(ip.GetNetworkPart()) {
			networksIps.insertToExistingNetworkIP(ip.GetNetworkPart(), ip)
		} else {
			networksIps = append(networksIps, IpsInNetwork{Network: ip.GetNetworkPart(), IPs: []IP{ip}})
		}
	}
	return networksIps
}
