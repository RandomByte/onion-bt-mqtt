#!/bin/bash
GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc -o out/onion-bt-mqtt onion-bt-mqtt.go
