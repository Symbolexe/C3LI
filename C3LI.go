package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Certificate struct {
	IssuerCaID              int       `json:"issuer_ca_id"`
	IssuerName              string    `json:"issuer_name"`
	CommonName              string    `json:"common_name"`
	NameValue               string    `json:"name_value"`
	ID                      int       `json:"id"`
	EntryTimestamp          CustomTime `json:"entry_timestamp"`
	NotBefore               CustomTime `json:"not_before"`
	NotAfter                CustomTime `json:"not_after"`
	SerialNumber            string    `json:"serial_number"`
	Source                  string    `json:"source"`
	AllDomains              []string  `json:"all_domains"`
	ValidationMethods       []string  `json:"validation_methods"`
	ValidationTimestamp     CustomTime `json:"validation_timestamp"`
	RevocationStatus        string    `json:"revocation_status"`
	RevocationReason        string    `json:"revocation_reason"`
	RevocationTimestamp     CustomTime `json:"revocation_timestamp"`
	SHA256Fingerprint       string    `json:"sha256_fingerprint"`
	SHA1Fingerprint         string    `json:"sha1_fingerprint"`
	MD5Fingerprint          string    `json:"md5_fingerprint"`
	SubjectAlternativeNames []string  `json:"subject_alternative_names"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse("2006-01-02T15:04:05.999", s)
	return
}

func fetchCertificates(domain string) ([]Certificate, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from crt.sh: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from crt.sh: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var certs []Certificate
	err = json.Unmarshal(body, &certs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return certs, nil
}

func printCertificates(w io.Writer, certs []Certificate, verbose bool) {
	for _, cert := range certs {
		fmt.Fprintf(w, "Common Name: %s\n", cert.CommonName)
		fmt.Fprintf(w, "Issuer Name: %s\n", cert.IssuerName)
		fmt.Fprintf(w, "Serial Number: %s\n", cert.SerialNumber)
		fmt.Fprintf(w, "Not Before: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(w, "Not After: %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))
		if verbose {
			fmt.Fprintf(w, "SHA256 Fingerprint: %s\n", cert.SHA256Fingerprint)
			fmt.Fprintf(w, "SHA1 Fingerprint: %s\n", cert.SHA1Fingerprint)
			fmt.Fprintf(w, "MD5 Fingerprint: %s\n", cert.MD5Fingerprint)
			fmt.Fprintf(w, "Subject Alternative Names: %v\n", cert.SubjectAlternativeNames)
		}
		fmt.Fprintln(w)
	}
}

func main() {
	domain := flag.String("url", "", "Target domain to search certificates for")
	verbose := flag.Bool("v", false, "Enable verbose output")
	sortBy := flag.String("sort", "", "Sort results by (issuer, expiration)")
	outputFile := flag.String("output", "", "File to write output (default is stdout)")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Usage: ./C3LI --url target.com [--output filename] [--v] [--sort issuer|expiration]")
		os.Exit(1)
	}

	fmt.Printf("Searching certificates for domain: %s\n", *domain)

	certs, err := fetchCertificates(*domain)
	if err != nil {
		fmt.Printf("Error fetching certificates: %s\n", err)
		os.Exit(1)
	}

	if *sortBy != "" {
		switch strings.ToLower(*sortBy) {
		case "issuer":
			sort.Slice(certs, func(i, j int) bool {
				return certs[i].IssuerName < certs[j].IssuerName
			})
		case "expiration":
			sort.Slice(certs, func(i, j int) bool {
				return certs[i].NotAfter.Time.Before(certs[j].NotAfter.Time)
			})
		default:
			fmt.Println("Invalid sort option. Available options: issuer, expiration")
			os.Exit(1)
		}
	}

	var out io.Writer = os.Stdout
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		out = file
	}

	printCertificates(out, certs, *verbose)

	fmt.Println("Results have been saved.")
}
