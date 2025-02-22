package main

import (
	"github.com/oschwald/maxminddb-golang/v2"
	"io"
	"log"
	"net/http"
	"net/netip"
)

func Any(arr []bool) bool {
	for _, v := range arr {
		if v {
			return true
		}
	}
	return false
}

func getPublicIP() (string, error) {
	// Use an external service to fetch the public IP
	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert response to string
	return string(body), nil
}

func getGeoIP(IP string) string {
	db, err := maxminddb.Open("data/geolite2-country-ipv4.mmdb")
	if err != nil {
		println(err)
	}
	defer db.Close()

	addr := netip.MustParseAddr(IP)

	var record struct {
		Code string `maxminddb:"country_code"`
	}

	err = db.Lookup(addr).Decode(&record)
	if err != nil {
		log.Panic(err)
	}

	return record.Code
}
