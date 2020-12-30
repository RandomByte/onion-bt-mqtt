package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	devices := parse("Scanning ...\n	12:34:56:78:90:42	Device 1\n	13:37:13:37:13:37	//Device $('2 ")

	if len(devices) != 2 {
		t.Errorf("Expected two devices, got %v", len(devices))
	}
	device1, ok := devices["12:34:56:78:90:42"]
	if ok != true || device1.Name != "Device 1" {
		t.Error("Wrong values for Device 1")
	}
	device2, ok := devices["13:37:13:37:13:37"]
	if ok != true || device2.Name != "//Device $('2 " {
		t.Error("Wrong values for Device 2")
	}

	devices = parse("    12:34:56:78:90:42      ")
	if len(devices) != 1 {
		t.Errorf("Expected one device, got %v", len(devices))
	}
	device3, ok := devices["12:34:56:78:90:42"]
	if ok != true || device3.Name != "12:34:56:78:90:42" { // name should fall back to MAC
		t.Error("Wrong values for Device 3")
	}
}
