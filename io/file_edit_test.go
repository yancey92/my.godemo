package myio

import "testing"

func TestInjectionIngress(t *testing.T) {
	err := InjectionIngress("dashboard.lenovo.com", "10.121.111.134")
	if err != nil {
		t.Fatal(err)
	}
}
