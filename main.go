package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	geo "github.com/kellydunn/golang-geo"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	topic := flag.String("topic", "nguyenvanduocit/feeds/tracker.gps", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://io.adafruit.com:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "nguyenvanduocit", "The User (optional)")
	flag.Parse()
	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID(*password).SetUsername(*user).SetPassword(*password)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	center := geo.NewPoint(10.773418, 106.674822)

	for {
		center = center.PointAtDistanceAndBearing(0.05, float64(rand.Int63n(360)))

		token := c.Publish(*topic, 0, false, `{"time":"`+time.Now().Format(time.RFC3339)+`", "tracked_satellites":7, "total_satellites":15, "is_fixed":"3D fix", "latitude":`+strconv.FormatFloat(center.Lat(), 'f', 6, 64)+`, "longitude": `+strconv.FormatFloat(center.Lng(), 'f', 6, 64)+`, "altitude":0.000000, "battery_voltage":3288, "battery_percent": 0}`)
		token.Wait()
		time.Sleep(2 * time.Second)
	}
}
