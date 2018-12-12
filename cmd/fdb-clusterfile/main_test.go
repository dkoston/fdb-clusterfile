package main

import (
	"github.com/go-test/deep"
	"testing"
)

// helper
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Test_isIP(t *testing.T) {
	type IPTest struct {
		HostOrIP string
		IP       bool
	}

	var testCases = []IPTest{
		{"www.google.com", false},
		{"foundationdb", false},
		{"127.0.0.1", true},
		{"localhost", false},
		{"10.253.1.24", true},
	}

	for i := 0; i < len(testCases); i++ {
		ip := isIP(testCases[i].HostOrIP)

		if ip != testCases[i].IP {
			t.Errorf("Expected %s to be IP: %t, Got: %t", testCases[i].HostOrIP, testCases[i].IP, ip)
		}
	}
}

func Test_TranslateHostArrayToIPs(t *testing.T) {
	type HostArrayTest struct {
		HostArray []string
		Expected  [][]string
	}

	one := make([][]string, 1)
	one = append(one, []string{"127.0.0.1:4500:tcp"})

	two := make([][]string, 1)
	two = append(one, []string{"127.0.0.1:9999"})

	three := make([][]string, 1)
	three = append(three, []string{"104.27.130.63:1234", "127.0.0.1:9999"})
	three = append(three, []string{"104.27.131.63:1234", "127.0.0.1:9999"})

	var testCases = []HostArrayTest{
		{[]string{"127.0.0.1:4500:tcp"}, one},
		{[]string{"127.0.0.1:9999"}, two},
		{[]string{"davekoston.com:1234", "127.0.0.1:9999"}, three},
	}

	for i := 0; i < len(testCases); i++ {
		translated := TranslateHostArrayToIPs(testCases[i].HostArray)

		found := false

		for j := 0; j < len(testCases[i].Expected); j++ {
			if diff := deep.Equal(translated, testCases[i].Expected[j]); diff == nil {
				found = true
			}
		}

		if !found {
			t.Errorf("Expected %s to be translated to : %s, Got: %s", testCases[i].HostArray, testCases[i].Expected, translated)
		}
	}
}

func Test_TranslateHostToIP(t *testing.T) {
	type HostTest struct {
		Host     string
		Expected []string
	}

	var testCases = []HostTest{
		{"localhost:1234", []string{"127.0.0.1:1234"}},
		{"127.0.0.1:4500:tcp", []string{"127.0.0.1:4500:tcp"}},
		{"127.0.0.1:9999", []string{"127.0.0.1:9999"}},
		{"davekoston.com:1234", []string{"104.27.130.63:1234", "104.27.131.63:1234"}},
		{"davekoston.com:4321:tcp", []string{"104.27.130.63:4321:tcp", "104.27.131.63:4321:tcp"}},
	}

	for i := 0; i < len(testCases); i++ {
		translated := TranslateHostToIP(testCases[i].Host)

		if !stringInSlice(translated, testCases[i].Expected) {
			t.Errorf("Expected %s to be translated to : %s, Got: %s", testCases[i].Host, testCases[i].Expected, translated)
		}
	}
}

func Test_TranslateFDBAddr(t *testing.T) {
	type testCase struct {
		FDBAddr  string
		Expected string
	}

	testCases := []testCase{
		{
			FDBAddr:  "fdb:fdb@localhost:4500:tcp",
			Expected: "fdb:fdb@127.0.0.1:4500:tcp",
		},
		{
			FDBAddr:  "fdb:fdb@localhost:4500:tcp,localhost:4501:tls",
			Expected: "fdb:fdb@127.0.0.1:4500:tcp,127.0.0.1:4501:tls",
		},
		{
			FDBAddr:  "fdb:fdb@localhost:4500:tcp,localhost:4501:tls,localhost:4502:tls",
			Expected: "fdb:fdb@127.0.0.1:4500:tcp,127.0.0.1:4501:tls,127.0.0.1:4502:tls",
		},
	}

	for i := 0; i < len(testCases); i++ {
		result := TranslateFDBAddr(testCases[i].FDBAddr)
		if result != testCases[i].Expected {
			t.Fatalf("expected (%s) got (%s)", testCases[i].Expected, result)
		}
	}
}
