package main

import (
	"fmt"
	"ipapp/packages/ip"
	"ipapp/packages/logger"
	"math"
	"os"
	"strconv"
	"strings"
)

func init() {
	logger.InitLogger()
}

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("Not enough arguments provided")
		os.Exit(0)
	}
	ip1Str := args[0]
	ip2Str := args[1]
	maskStr := args[2]
	if !ipIsValid(ip1Str) || !ipIsValid(ip2Str) {
		fmt.Println("One or more ip addresses were invalid")
		os.Exit(2)
	}
	ip1 := ip.Ip{}
	ip1.MakeIpDecStruct(ip1Str)
	ip2 := ip.Ip{}
	ip2.MakeIpDecStruct(ip2Str)
	maskNum, err := strconv.Atoi(maskStr)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}
	mask := ip.Ip{}
	mask.MakeMask(maskNum)
	const maxBits = 32
	bitsReservedForHost := maxBits - maskNum
	numberOfHosts := math.Pow(2, float64(bitsReservedForHost)) - 2
	numberOfSubnetworks := math.Pow(2, float64(maskNum))
	logger.InfoLogger.Println("mask", mask)
	logger.InfoLogger.Println("numberOfHosts", numberOfHosts)
	logger.InfoLogger.Println("numberOfSubnetworks", numberOfSubnetworks)

	ExamineNetwork(ip1.GetNetworkPart(maskNum))
	ExamineNetwork(ip2.GetNetworkPart(maskNum))

}

func ExamineNetwork(networkIp ip.Ip) {
	minIpInNetwork := networkIp.GetMinIpInNetwork()
	maxIpInNetwork := networkIp.GetMaxIpInNetwork()
	broadcastAddress := networkIp.GetBroadcastAddress()
	logger.InfoLogger.Println("minIpInNetwork", minIpInNetwork)
	logger.InfoLogger.Println("maxIpInNetwork", maxIpInNetwork)
	logger.InfoLogger.Println("broadcastAddress", broadcastAddress)
	logger.InfoLogger.Println("minIpInNetworkBin", minIpInNetwork.GetIpInBin())
	logger.InfoLogger.Println("maxIpInNetworkBin", maxIpInNetwork.GetIpInBin())
	logger.InfoLogger.Println("broadcastAddressBin", broadcastAddress.GetIpInBin())
}

// Only ipv4 adresses are considered valid
func ipIsValid(ip string) bool {
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
