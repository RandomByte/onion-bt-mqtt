# onion-bt-mqtt
Scan for Bluetooth devices and log them via MQTT on an **[Onion Omega](https://onion.io/)**

Based off [github.com/RandomByte/onion-bt-wardriving](https://github.com/RandomByte/onion-bt-wardriving)

## Dependencies
* **hcitool**  
    On an Onion Omega this can be installed by executing the following commands:
	```sh
	opkg update
	opkg install bluez-libs bluez-utils
	reboot
	```
	*Note: You might need to add the official OpenWRT repositories in `/etc/opkg/distfeeds.conf` by uncommented the respective lines*

## Startup script
1. Copy `onion-bt-mqtt.init_file` to `/etc/init.d/onion-bt-mqtt`
1. Change `PROJECT_DIR`, `MQTT_BROKER`, `MQTT_TOPIC_PREFIX` and other settings if required
1. Do `chmod +x /etc/init.d/onion-bt-mqtt` to make it executable
1. Do `/etc/init.d/onion-bt-mqtt enable` to enable autostart on boot
1. Do `/etc/init.d/onion-bt-mqtt start` to start it for the first time
1. Do `tail -f /var/log/onion-bt-mqtt.log` and check for any errors. If you see nothing, wait a moment as nodejs takes some time to launch. When everything's good, the last line should read "Running...".

## MQTT Message Examples

| Topic        | Payload
| ------------- |-------------|
| `Home/Bluetooth/CA:D8:17:E3:C1:1B` | `Bernds Speakers` |
| `Home/Bluetooth/E2:27:DE:F2:C7:B1` | `Somsong TV` |
