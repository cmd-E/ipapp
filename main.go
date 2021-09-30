package main

import (
	"flag"
	"fmt"
	"ipapp/packages/ip"
	"ipapp/packages/logger"
	"ipapp/packages/printer"
	"math"
	"os"
)

var ip1Str string
var ip2Str string
var maskNum int

func init() {
	logger.InitLogger()
	flag.StringVar(&ip1Str, "ip1", "", "first ip address")
	flag.StringVar(&ip2Str, "ip2", "", "second ip address")
	flag.IntVar(&maskNum, "m", 24, "mask number (24 by default)")
}

func main() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}
	if !ip.IpIsValid(ip1Str) || !ip.IpIsValid(ip2Str) {
		fmt.Println("One or more ip addresses are invalid")
		flag.Usage()
		os.Exit(1)
	}
	if maskNum < 0 || maskNum > 32 {
		fmt.Println("Mask cannot be less than 0 or greater than 32")
		flag.Usage()
		os.Exit(1)
	}
	ip1 := ip.ParseIpFromDecimalString(ip1Str)
	ip1.MaskNum = maskNum
	ip2 := ip.ParseIpFromDecimalString(ip2Str)
	ip2.MaskNum = maskNum
	mask := ip.MakeMask(maskNum)
	const maxBits = 32
	bitsReservedForHost := maxBits - maskNum
	numberOfHosts := math.Pow(2, float64(bitsReservedForHost)) - 2
	numberOfSubnetworks := math.Pow(2, float64(maskNum))
	printer.PrintFormatted("Mask in binary representation", mask.GetIpInBin())
	printer.PrintFormatted("Number of hosts", numberOfHosts)
	printer.PrintFormatted("Number of subnetworks", numberOfSubnetworks)
	fmt.Println()
	fmt.Printf("Examining %s...\n", ip1.GetIpInDec())
	ExamineNetwork(ip1.GetNetworkPart())
	fmt.Println()
	fmt.Printf("Examining %s...\n", ip2.GetIpInDec())
	ExamineNetwork(ip2.GetNetworkPart())
}

func ExamineNetwork(networkIp ip.Ip) {
	printer.PrintFormatted("Network ip address DEC", networkIp.GetIpInDec())
	printer.PrintFormatted("Network ip address BIN", networkIp.GetIpInBin())
	fmt.Println()
	printer.PrintFormatted("Minimal ip address DEC", networkIp.GetMinIpInNetwork().GetIpInDec())
	printer.PrintFormatted("Minimal ip address BIN", networkIp.GetMinIpInNetwork().GetIpInBin())
	fmt.Println()
	printer.PrintFormatted("Maximum ip address DEC", networkIp.GetMaxIpInNetwork().GetIpInDec())
	printer.PrintFormatted("Maximum ip address BIN", networkIp.GetMaxIpInNetwork().GetIpInBin())
	fmt.Println()
	printer.PrintFormatted("Broadcast address DEC", networkIp.GetBroadcastAddress().GetIpInDec())
	printer.PrintFormatted("Broadcast address BIN", networkIp.GetBroadcastAddress().GetIpInBin())
}
