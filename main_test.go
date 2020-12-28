package main

import "testing"

func TestMatchServiceName(t *testing.T) {
	if !regexpAutoCreatedIngressService.MatchString("ingress-3beae2ef678a239ca440081fa6d5663e") {
		t.Fatal("bad")
	}
}
