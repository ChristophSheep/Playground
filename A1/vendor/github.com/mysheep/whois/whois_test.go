package whois_test

import (
	"fmt"
	"testing"

	"github.com/mysheep/whois"
)

func TestWhoisWithIPV4(t *testing.T) {
	m, err := whois.Get("199.232.18.133")

	if err != nil {
		t.Errorf("Map error")
	}

	for k, v := range m {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	org := m["Organization"]

	if org != "Fastly (SKYCA-3)" {
		t.Errorf("'Fastly (SKYCA-3)' was expected, but was %s", org)
	}
}
