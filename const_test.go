package main

import "testing"

func TestParseStationNames(t *testing.T) {
	parseStationNames()
	if len(StationMap) < 1 {
		t.Fatal("parseStationNames failed!")
	}
	if StationMap["北京北"] != "VAP" {
		t.Fatal("parseStationNames don't contain 北京北!")
	}
}
