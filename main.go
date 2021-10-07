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

var maskNum int

func init() {
	logger.InitLogger()
	flag.IntVar(&maskNum, "m", 24, "mask number (24 by default)")
}

func main() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}
	providedIps := os.Args[1:]
	if len(providedIps) < 1 {
		fmt.Println("Ip addresses weren't provided")
		flag.Usage()
		os.Exit(1)
	}
	invalidIps := ""
	for i := 0; i < len(providedIps); i++ {
		if !ip.IpIsValid(providedIps[i]) {
			invalidIps += fmt.Sprintf("%s\n", providedIps[i])
			providedIps = append(providedIps[:i], providedIps[i+1:]...) // TODO: check for out of range
			i--
		}
	}
	if len(invalidIps) != 0 {
		fmt.Printf("Some of the provided IPs are invalid:\n%s", invalidIps)
	}
	if len(providedIps) == 0 {
		fmt.Println("No IPs to examine")
		os.Exit(0)
	}
	if maskNum < 0 || maskNum > 32 {
		fmt.Println("Mask cannot be less than 0 or greater than 32")
		os.Exit(1)
	}
	var ips []ip.Ip
	for _, ipDecStr := range providedIps {
		ips = append(ips, ip.ParseIpFromDecimalString(ipDecStr, maskNum))
	}
	const maxBits = 32
	bitsReservedForHost := maxBits - maskNum
	const numberOfReservedIps = 2
	numberOfHosts := math.Pow(2, float64(bitsReservedForHost)) - numberOfReservedIps
	numberOfSubnetworks := math.Pow(2, float64(maskNum))
	printer.PrintFormatted("Mask in binary representation", ip.MakeMask(maskNum).GetIpInBin())
	printer.PrintFormatted("Number of hosts", numberOfHosts)
	printer.PrintFormatted("Number of subnetworks", numberOfSubnetworks)
	fmt.Println()
	for _, ip := range ips {
		fmt.Printf("Examining %s...\n", ip.GetIpInDec())
		ExamineNetwork(ip.GetNetworkPart())
		fmt.Println()
	}
	for _, v := range ip.SortIpsByNetworks(ips) {
		fmt.Printf("Ips in %s network:\n", v.Network.GetIpInDec())
		for _, ip := range v.IPs {
			fmt.Println(ip.GetIpInDec())
		}
	}
}

func ExamineNetwork(networkIp ip.Ip) {
	printer.PrintFormatted("Network ip address DEC | BIN", fmt.Sprintf("%s | %s", networkIp.GetIpInDec(), networkIp.GetIpInBin()))
	printer.PrintFormatted("Minimum ip address DEC | BIN", fmt.Sprintf("%s | %s", networkIp.GetMinIpInNetwork().GetIpInDec(), networkIp.GetMinIpInNetwork().GetIpInBin()))
	printer.PrintFormatted("Maximum ip address DEC | BIN", fmt.Sprintf("%s | %s", networkIp.GetMaxIpInNetwork().GetIpInDec(), networkIp.GetMaxIpInNetwork().GetIpInBin()))
	printer.PrintFormatted("Broadcast address DEC | BIN", fmt.Sprintf("%s | %s", networkIp.GetBroadcastAddress().GetIpInDec(), networkIp.GetBroadcastAddress().GetIpInBin()))
}
