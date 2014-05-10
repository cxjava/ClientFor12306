package main

import "testing"

func TestReadConfig(t *testing.T) {
	ReadConfig()
	if len(Conf.CDN) < 1 {
		t.Fatal("ReadConfig failed!")
	}
	if len(Conf.ProxyUrl) < 1 {
		t.Fatal("ReadConfig don't contain proxyUrl!")
	}
}
