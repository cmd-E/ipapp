package ip

import (
	"fmt"
	"strconv"
	"strings"
)

var numOfDots = 0 // compensate indexes

type Ip struct {
	Part1 int
	Part2 int
	Part3 int
	Part4 int
}

// Only ipv4 adresses are considered valid
func IpIsValid(ip string) bool {
	splittedIp := strings.Split(ip, ".")
	if len(splittedIp) != 4 {
		return false
	}
	for _, v := range splittedIp {
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

func ParseIpFromDecimalString(ipStr string) Ip {
	var ip Ip
	splittedIp := strings.Split(ipStr, ".")
	ip.Part1, _ = strconv.Atoi(splittedIp[0])
	ip.Part2, _ = strconv.Atoi(splittedIp[1])
	ip.Part3, _ = strconv.Atoi(splittedIp[2])
	ip.Part4, _ = strconv.Atoi(splittedIp[3])
	return ip
}

func ParseIpFromBinString(ipStr string) Ip {
	var ip Ip
	splittedIp := strings.Split(ipStr, ".")
	t1, _ := strconv.ParseInt(splittedIp[0], 2, 32)
	t2, _ := strconv.ParseInt(splittedIp[1], 2, 32)
	t3, _ := strconv.ParseInt(splittedIp[2], 2, 32)
	t4, _ := strconv.ParseInt(splittedIp[3], 2, 32)
	ip.Part1 = int(t1)
	ip.Part2 = int(t2)
	ip.Part3 = int(t3)
	ip.Part4 = int(t4)
	return ip
}

func MakeMask(maskNum int) Ip {
	var mask Ip
	var ipPart string
	for i := 1; i <= 32; i++ {
		if i <= maskNum {
			ipPart += "1"
		} else {
			ipPart += "0"
		}
		if i == 8 || i == 16 || i == 24 {
			ipPart += "."
			numOfDots += 1
		}
	}
	splittedIp := strings.Split(ipPart, ".")
	t1, _ := strconv.ParseInt(splittedIp[0], 2, 32)
	t2, _ := strconv.ParseInt(splittedIp[1], 2, 32)
	t3, _ := strconv.ParseInt(splittedIp[2], 2, 32)
	t4, _ := strconv.ParseInt(splittedIp[3], 2, 32)
	mask.Part1 = int(t1)
	mask.Part2 = int(t2)
	mask.Part3 = int(t3)
	mask.Part4 = int(t4)
	return mask
}

func (ip Ip) GetNetworkPart(maskNum int) Ip {
	binIp := ip.GetIpInBin()
	binIp = binIp[:maskNum+numOfDots] // TODO: check for out of range
	for i := maskNum + 1; i <= 32; i++ {
		binIp += "0"
		if i == 8 || i == 16 || i == 24 {
			binIp += "."
		}
	}
	splittedIp := strings.Split(binIp, ".")
	var networkIp Ip
	t1, _ := strconv.ParseInt(splittedIp[0], 2, 32)
	t2, _ := strconv.ParseInt(splittedIp[1], 2, 32)
	t3, _ := strconv.ParseInt(splittedIp[2], 2, 32)
	t4, _ := strconv.ParseInt(splittedIp[3], 2, 32)
	networkIp.Part1 = int(t1)
	networkIp.Part2 = int(t2)
	networkIp.Part3 = int(t3)
	networkIp.Part4 = int(t4)
	return networkIp
}

func (ip Ip) GetMinIpInNetwork(mask int) Ip {
	ipSplit := []rune(ip.GetIpInBin())
	for i := mask + numOfDots; i < 32+numOfDots; i++ {
		if ipSplit[i] == '.' {
			i -= 1
			continue
		}
		ipSplit[i] = '0'
	}
	ipSplit[len(ipSplit)-1] = '1'
	return ParseIpFromBinString(string(ipSplit))
}

func (ip Ip) GetMaxIpInNetwork(mask int) Ip {
	ipSplit := []rune(ip.GetIpInBin())
	for i := mask + numOfDots; i < 32+numOfDots; i++ {
		if ipSplit[i] == '.' {
			i -= 1
			continue
		}
		ipSplit[i] = '1'
	}
	ipSplit[len(ipSplit)-1] = '0'
	return ParseIpFromBinString(string(ipSplit))
}

func (ip Ip) GetBroadcastAddress(mask int) Ip {
	ipSplit := []rune(ip.GetIpInBin())
	for i := mask + numOfDots; i < 32+numOfDots; i++ {
		if ipSplit[i] == '.' {
			i -= 1
			continue
		}
		ipSplit[i] = '1'
	}
	return ParseIpFromBinString(string(ipSplit))
}

func (ip Ip) GetIpInBin() string {
	return fmt.Sprintf("%08b.%08b.%08b.%08b", ip.Part1, ip.Part2, ip.Part3, ip.Part4)
}

func (ip Ip) GetIpInDec() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip.Part1, ip.Part2, ip.Part3, ip.Part4)
}
