package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type device struct {
	Name string
}

var deviceRe = regexp.MustCompile("(?im)^[^0-9a-f]*((?:[0-9a-f]{2}:){5}[0-9a-f]{2})\\s*([^\\s].*)?$")
var endsWithSlashRe = regexp.MustCompile("\\/$")

func main() {
	mqttBroker := flag.String("mqtt-broker", "", "MQTT broker")
	mqttTopicPrefix := flag.String("mqtt-topic-prefix", "", "MQTT topic prefix")
	flag.Parse()

	if *mqttBroker == "" || *mqttTopicPrefix == "" {
		fmt.Println("Missing or empty mqtt-broker or mqtt-topic-prefix parameters")
		return
	}

	if !endsWithSlashRe.MatchString(*mqttTopicPrefix) {
		fmt.Println("mqtt-topic-prefix parameter value must end with a forward slash")
		return
	}

	setupBt()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer signal.Stop(sig)

	fmt.Println("Broker:", *mqttBroker)
	fmt.Println("Topic Prefix:", *mqttTopicPrefix)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(*mqttBroker)

	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case <-time.After(5000 * time.Millisecond):
			loop(mqttClient, *mqttTopicPrefix)
		case s := <-sig:
			fmt.Println("Got signal:", s)
			fmt.Println("Quitting...")
			return
		}
	}
}

func loop(mqttClient MQTT.Client, mqttTopicPrefix string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
		}
	}()
	result := scan()
	parsed := parse(result)

	for mac, device := range parsed {
		fmt.Printf("Detected device %s with MAC address %s\n", device.Name, mac)
		publish(mqttClient, mqttTopicPrefix, mac, device)
		blinkLed()
	}
}

func publish(mqttClient MQTT.Client, mqttTopicPrefix string, mac string, device device) {
	token := mqttClient.Publish(mqttTopicPrefix+mac, 2, false, device.Name)
	token.Wait()
}

func scan() string {
	// Create an *exec.Cmd
	cmd := exec.Command("hcitool", "scan", "--flush")

	// Stdout buffer
	cmdOutput := &bytes.Buffer{}
	// Attach buffer to command
	cmd.Stdout = cmdOutput

	err := cmd.Run() // will wait for command to return
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return cmdOutput.String()
}

func parse(rawScanResult string) map[string]device {
	matches := deviceRe.FindAllStringSubmatch(rawScanResult, -1)
	devices := make(map[string]device)
	for _, match := range matches {
		name := match[2]
		if name == "" {
			name = match[1]
		}
		devices[match[1]] = device{Name: name}
	}

	return devices
}
func setupBt() {
	cmd := exec.Command("hciconfig", "hci0", "up")

	err := cmd.Run() // will wait for command to return
	if err != nil {
		fmt.Println(err)
	}
}

func blinkLed() {
	cmdBlue := exec.Command("expled", "0x0000ff")

	err := cmdBlue.Run() // will wait for command to return
	if err != nil {
		fmt.Println(err)
	}

	cmdOff := exec.Command("expled", "0x000000")
	err = cmdOff.Run() // will wait for command to return
	if err != nil {
		fmt.Println(err)
	}
}
