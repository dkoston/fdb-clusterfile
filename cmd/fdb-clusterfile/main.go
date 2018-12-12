package main

import (
	"io/ioutil"
	"github.com/jessevdk/go-flags"
	"log"
	"net"
	"os"
	"strings"
)

type CommandLineOptions struct {
	FDBAddr     string `long:"fdb_addr" description:"FDB connection address (fdb.cluster contents)" env:"FDB_ADDR" default:"fdb:fdb@localhost:4500"`
	ClusterFile string `long:"cluster_file" description:"full path to fdb.cluster" env:"FDB_CLUSTER_FILE" default:"fdb.cluster"`
}

const name = "fdb-clusterfile"
const version = "0.0.1"

func main() {
	var opts CommandLineOptions
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("[%s] Unable to parse command line options: %v", name, err)
	}

	contents, err := WriteClusterFile(opts.FDBAddr, opts.ClusterFile)
	if err != nil {
		log.Fatalf("[%s] Failed to write cluster file: %v", name, err)
		os.Exit(1)
	}

	log.Printf("[%s] Wrote cluster file (%s) with contents (%s)", name, opts.ClusterFile, contents)
}

func isIP(ip string) bool {
	ip = strings.Trim(ip, " ")
	r := net.ParseIP(ip)

	if r != nil {
		return true
	}
	return false
}

// Checks for hostnames in the fdbAddr and converts them to IPs
// as FDB doesn't allow hostnames in fdb.cluster
func TranslateFDBAddr(fdbAddr string) string {
	// fdb:fdb@127.0.0.1:4500:tcp,hostname:port:tls

	parts := strings.Split(fdbAddr, ",")
	if len(parts) > 1 {
		for i := 1; i < len(parts); i++ {
			parts[i] = TranslateHostToIP(parts[i])
		}
	}
	firstHostParts := strings.Split(parts[0], "@")
	firstHostParts[1] = TranslateHostToIP(firstHostParts[1])
	parts[0] = strings.Join(firstHostParts, "@")
	return strings.Join(parts, ",")
}

func TranslateHostArrayToIPs(addrArray []string) []string {
	for i := 0; i < len(addrArray); i++ {
		addrArray[i] = TranslateHostToIP(addrArray[i])
	}
	return addrArray
}

func TranslateHostToIP(addr string) string {
	parts := strings.Split(addr, ":")
	hostnameOrIP := parts[0]

	if hostnameOrIP == "localhost" {
		parts[0] = "127.0.0.1"
		hostnameOrIP = "127.0.0.1"
	}

	if !isIP(hostnameOrIP) {
		addr, err := net.LookupIP(hostnameOrIP)
		if err != nil {
			log.Fatalf("Invalid Hostname: %v", err)
		}
		parts[0] = addr[0].String()
	}
	return strings.Join(parts, ":")
}

func WriteClusterFile(fdbAddr string, clusterFilePath string) (string, error) {
	translatedAddr := TranslateFDBAddr(fdbAddr)
	err := ioutil.WriteFile(clusterFilePath, []byte(translatedAddr), 0644)
	return translatedAddr, err
}
