package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"ipapp/packages/ip"
	"os"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

// run owns its flags and output streams so invocations can be tested independently.
func run(args []string, stdout, stderr io.Writer) int {
	options := flag.NewFlagSet("ipapp", flag.ContinueOnError)
	prefix := options.Int("m", 24, "IPv4 prefix length (0-32; default 24)")
	// Flag parse errors may contain user-controlled control characters. Render the
	// returned error with quoting instead of echoing the parser's raw diagnostics.
	options.SetOutput(io.Discard)
	if err := options.Parse(args); err != nil {
		options.SetOutput(stderr)
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(stderr, "Usage: ipapp [-m prefix] IPv4 [IPv4 ...]")
			options.PrintDefaults()
			return 0
		}
		fmt.Fprintf(stderr, "Invalid arguments: %q\n", err.Error())
		return 2
	}
	mask, err := ip.ConvertMask(*prefix)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if options.NArg() == 0 {
		fmt.Fprintln(stderr, "IP addresses weren't provided")
		return 1
	}
	addresses := make([]ip.IP, 0, options.NArg())
	exitCode := 0
	for _, text := range options.Args() {
		address, err := ip.ParseIPFromDecimalString(text, *prefix)
		if err != nil {
			fmt.Fprintln(stderr, err)
			exitCode = 1
			continue
		}
		addresses = append(addresses, address)
	}
	if len(addresses) == 0 {
		fmt.Fprintln(stderr, "No IPs to examine")
		return 1
	}
	output := bufio.NewWriter(stdout)
	printFormatted(output, "Mask in binary representation", mask.GetIPInBin())
	printFormatted(output, "Number of hosts", addresses[0].HostCount())
	printFormatted(output, "Number of subnetworks", uint64(1)<<uint(*prefix))
	fmt.Fprintln(output)
	for _, address := range addresses {
		fmt.Fprintf(output, "Examining %s...\n", address.GetIPInDec())
		examineNetwork(output, address.GetNetworkPart())
		fmt.Fprintln(output)
	}
	for _, group := range ip.SortIpsByNetworks(addresses) {
		fmt.Fprintf(output, "Ips in %s network:\n", group.Network.GetIPInDec())
		for _, address := range group.IPs {
			fmt.Fprintln(output, address.GetIPInDec())
		}
	}
	if err := output.Flush(); err != nil {
		fmt.Fprintf(stderr, "Cannot write report: %v\n", err)
		return 1
	}
	return exitCode
}

func printFormatted(output io.Writer, definition string, value interface{}) {
	fmt.Fprintf(output, "%-29s %v\n", definition, value)
}

func examineNetwork(output io.Writer, network ip.IP) {
	format := func(address ip.IP) string {
		return fmt.Sprintf("%s | %s", address.GetIPInDec(), address.GetIPInBin())
	}
	printFormatted(output, "Network ip address DEC | BIN", format(network))
	printFormatted(output, "Minimum ip address DEC | BIN", format(network.GetMinIPInNetwork()))
	printFormatted(output, "Maximum ip address DEC | BIN", format(network.GetMaxIPInNetwork()))
	if network.PrefixLength() < 31 {
		printFormatted(output, "Broadcast address  DEC | BIN", format(network.GetBroadcastAddress()))
	} else {
		printFormatted(output, "Broadcast address  DEC | BIN", "N/A (no directed broadcast for /31 or /32)")
	}
}
