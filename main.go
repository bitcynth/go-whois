package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
)

var regexExtractWhoisServer = regexp.MustCompile(`whois:\s+([a-z0-9\-\.]+)`)

func queryWhois(query string, addr string) string {
	conn, _ := net.Dial("tcp", addr)
	fmt.Fprintf(conn, query+"\r\n")
	data, _ := ioutil.ReadAll(conn)
	return string(data)
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	query := os.Args[1]
	data := queryWhois(query, "whois.iana.org:43")
	response := "---- whois.iana.org -----\n\n"
	response = response + data + "\n\n"
	for {
		referralSearch := regexExtractWhoisServer.FindStringSubmatch(data)
		if len(referralSearch) < 2 {
			break
		}
		data = queryWhois(query, referralSearch[1]+":43")
		response = response + "---- " + referralSearch[1] + " -----\n\n"
		response = response + data + "\n\n"
	}
	fmt.Printf("%s\n", response)
}
