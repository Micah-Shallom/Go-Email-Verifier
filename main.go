package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type domainResults struct {
	hasMX, hasSPF, hasDMARC        bool
	spfRecord, dmarcRecord, domain string
}

func main() {
	var results domainResults

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	for scanner.Scan() {
		results.domain = scanner.Text()
		checkDomain(&results)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(results *domainResults) {
	const dmarcSubdomain = "_dmarc."

	go func() {
		mxRecords, err := net.LookupMX(results.domain)
	
		if len(mxRecords) > 0 {
			results.hasMX = true
		}
	
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}()


	txtRecords, err := net.LookupTXT(results.domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			results.hasSPF = true
			results.spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT(dmarcSubdomain + results.domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			results.hasDMARC = true
			results.dmarcRecord = record
			break
		}
	}

	fmt.Printf("Domain: %v\n, hasMX: %v\n, hasSPF: %v\n, spfRecord: %v\n, hasDMARC: %v\n, dmarcRecord: %v\n", results.domain, results.hasMX, results.hasSPF, results.spfRecord, results.hasDMARC, results.dmarcRecord)
}
