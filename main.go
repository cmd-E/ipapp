package main

import (
	"flag"
	"fmt"
	"ipapp/packages/ip"
	"ipapp/packages/logger"
	"ipapp/packages/printer"
	"ipapp/packages/utils"
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
	providedIps = utils.FilterOutFlags(providedIps)
	if len(providedIps) < 1 {
		fmt.Println("Ip addresses weren't provided")
		os.Exit(1)
	}
	invalidIps := ""
	for i := 0; i < len(providedIps); i++ {
		if !ip.IsValid(providedIps[i]) {
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
	var ips []ip.IP
	for _, ipDecStr := range providedIps {
		ips = append(ips, ip.ParseIPFromDecimalString(ipDecStr, maskNum))
	}
	const maxBits = 32
	bitsReservedForHost := maxBits - maskNum
	const numberOfReservedIps = 2
	numberOfHosts := math.Pow(2, float64(bitsReservedForHost)) - numberOfReservedIps
	numberOfSubnetworks := math.Pow(2, float64(maskNum))
	printer.PrintFormatted("Mask in binary representation", ip.ConvertMask(maskNum).GetIPInBin())
	printer.PrintFormatted("Number of hosts", numberOfHosts)
	printer.PrintFormatted("Number of subnetworks", numberOfSubnetworks)
	fmt.Println()
	for _, ip := range ips {
		fmt.Printf("Examining %s...\n", ip.GetIPInDec())
		examineNetwork(ip.GetNetworkPart())
		fmt.Println()
	}
	for _, v := range ip.SortIpsByNetworks(ips) {
		fmt.Printf("Ips in %s network:\n", v.Network.GetIPInDec())
		for _, ip := range v.IPs {
			fmt.Println(ip.GetIPInDec())
		}
	}
}

func examineNetwork(networkIP ip.IP) {
	printer.PrintFormatted("Network ip address DEC | BIN", fmt.Sprintf("%s | %s", networkIP.GetIPInDec(), networkIP.GetIPInBin()))
	printer.PrintFormatted("Minimum ip address DEC | BIN", fmt.Sprintf("%s | %s", networkIP.GetMinIPInNetwork().GetIPInDec(), networkIP.GetMinIPInNetwork().GetIPInBin()))
	printer.PrintFormatted("Maximum ip address DEC | BIN", fmt.Sprintf("%s | %s", networkIP.GetMaxIPInNetwork().GetIPInDec(), networkIP.GetMaxIPInNetwork().GetIPInBin()))
	printer.PrintFormatted("Broadcast address  DEC | BIN", fmt.Sprintf("%s | %s", networkIP.GetBroadcastAddress().GetIPInDec(), networkIP.GetBroadcastAddress().GetIPInBin()))
}
