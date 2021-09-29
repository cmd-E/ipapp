package ip

import (
	"fmt"
	"strconv"
	"strings"
)

type Ip struct {
	Part1 int
	Part2 int
	Part3 int
	Part4 int
}

func (ip *Ip) MakeIpDecStruct(ipStr string) {
	splittedIp := strings.Split(ipStr, ".")
	ip.Part1, _ = strconv.Atoi(splittedIp[0])
	ip.Part2, _ = strconv.Atoi(splittedIp[1])
	ip.Part3, _ = strconv.Atoi(splittedIp[2])
	ip.Part4, _ = strconv.Atoi(splittedIp[3])
}

func (ip *Ip) MakeMask(maskNum int) {
	var ipPart string
	for i := 0; i <= 32; i++ {
		if i <= maskNum {
			ipPart += "1"
		} else {
			ipPart += "0"
		}
		if i == 8 || i == 16 || i == 24 || i == 32 {
			ipPart += "."
		}
	}
	ipPart = ipPart[:len(ipPart)-1]
	splittedIp := strings.Split(ipPart, ".")
	ip.Part1, _ = strconv.Atoi(splittedIp[0])
	ip.Part2, _ = strconv.Atoi(splittedIp[1])
	ip.Part3, _ = strconv.Atoi(splittedIp[2])
	ip.Part4, _ = strconv.Atoi(splittedIp[3])
}

func (ip Ip) GetNetworkPart(maskNum int) Ip {
	binIp := fmt.Sprintf("%08b.%08b.%08b.%08b", ip.Part1, ip.Part2, ip.Part3, ip.Part4)
	binIp = binIp[:maskNum+1]
	for i := maskNum + 1; i <= 32; i++ {
		binIp += "0"
		if i == 8 || i == 16 || i == 24 || i == 32 {
			binIp += "."
		}
	}
	binIp = binIp[:len(binIp)-1]
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

func (ip Ip) GetMinIpInNetwork() Ip {
	ip.Part4 = ip.Part4 + 1
	return ip
}

func (ip Ip) GetMaxIpInNetwork() Ip {
	ip.Part4 = 254
	return ip
}

func (ip Ip) GetBroadcastAddress() Ip {
	ip.Part4 = 255
	return ip
}

func (ip Ip) GetIpInBin() string {
	return fmt.Sprintf("%08b.%08b.%08b.%08b", ip.Part1, ip.Part2, ip.Part3, ip.Part4)
}
