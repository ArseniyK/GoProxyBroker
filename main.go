package main

import (
	"ProxyBroker/types"
	"flag"
	"fmt"
	"os"
	"strings"
)

var levelMap = map[string]types.ProxyLevel{
	"transparent": types.TRANSPARENT,
	"anonymous":   types.ANONYMOUS,
	"high":        types.HIGH,
}

func main() {
	// Subcommands

	findCmd := flag.NewFlagSet("find", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	// Common flags (applies to both subcommands)
	var countriesString string
	var levels []types.ProxyLevel
	addCommonFlags := func(f *flag.FlagSet) {
		f.StringVar(&countriesString, "countries", "", "Comma-separated country codes (e.g., US,GB,DE)")
		f.Func("lvl", "Comma-separated proxy levels (transparent, anonymous, high)", func(s string) error {
			values := strings.Split(s, ",")
			for _, v := range values {
				v = strings.TrimSpace(v)
				level, exists := levelMap[v]
				if !exists {
					return fmt.Errorf("invalid level: %s (allowed: transparent, anonymous, high)", v)
				}
				levels = append(levels, level)
			}
			return nil
		})
	}

	host := serveCmd.String("host", "127.0.0.1", "The host of the server")
	port := serveCmd.Int("port", 8080, "The port on which the server listens")
	addCommonFlags(serveCmd)

	limit := findCmd.Int("limit", 0, "The maximum number of proxies")
	check := findCmd.Bool("check", true, "Check found proxies")
	addCommonFlags(findCmd)

	// Ensure a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Expected 'find' or 'serve' subcommands")
		printMainUsage()
		os.Exit(1)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "find":
		findCmd.Parse(os.Args[2:])
		countries := parseCountries(countriesString)
		executeFind(*limit, *check, countries, levels)
	case "serve":
		serveCmd.Parse(os.Args[2:])
		countries := parseCountries(countriesString)
		executeServe(*host, *port, countries, levels)
	default:
		fmt.Println("Expected 'find' or 'serve' subcommands")
		printMainUsage()
		os.Exit(1)
	}

}

// Display usage for the main program
func printMainUsage() {
	fmt.Println("Usage: ProxyBroker <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  find     Finds proxies with specified options")
	fmt.Println("  serve    Serves proxies over a network or host")
	fmt.Println("\nUse '<command> -h' for more details about a command.\n")
}

// Validates and parses common flags (countries and levels)
func parseCountries(countriesString string) []string {
	var countries []string
	if countriesString != "" {
		countries = strings.Split(countriesString, ",")
	} else {
		countries = make([]string, 0)
	}
	return countries
}

// Execute the 'find' subcommand
func executeFind(limit int, check bool, countries []string, levels []types.ProxyLevel) {
	b := Broker{}
	b.Init(countries, levels)
	b.Find(limit, check)
}

// Execute the 'serve' subcommand
func executeServe(host string, port int, countries []string, levels []types.ProxyLevel) {
	b := Broker{}
	b.Init(countries, levels)
	b.Serve(host, port)
}
