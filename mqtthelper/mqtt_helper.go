package mqtthelper

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/quokka2020/gohelpers/util"
	// "strconv"
)

var qos = 0
var broker = flag.String("broker", "tcp://192.168.10.4:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
var password = flag.String("mqttpassword", "", "The password (optional)")
var user = flag.String("mqttuser", "", "The User (optional)")
var cleansess = flag.Bool("clean", false, "Set Clean Session (default false)")

type MqttHelper struct {
	client MQTT.Client
	Prefix string // also the name
}

func CreateMqttHelper(prefix string) *MqttHelper {
	helper := MqttHelper{
		Prefix: prefix,
	}
	id := fmt.Sprintf("%s-%s", filepath.Base(os.Args[0]), util.GetEnv("HOSTNAME", "localhost"))
	if util.Verbose() {
		log.Printf("myname==id %s", id)
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID(id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetConnectRetryInterval(30 * time.Second)
	opts.SetConnectionLostHandler(func(c MQTT.Client, err error) {
		log.Printf("just lost mqtt connection err:%v", err)
	})
	opts.SetWill(
		helper.topic("connected"),
		"offline",
		0,
		true,
	)
	opts.SetOnConnectHandler(helper.onConnect)

	// opts.SetDefaultPublishHandler(data.msgReceived)

	helper.client = MQTT.NewClient(opts)
	if token := helper.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// announce to HomeAssistant

	// registerConfig(client)

	return &helper
}

func (helper *MqttHelper) topic(subtopic string) string {
	return fmt.Sprintf("%s/%s", helper.Prefix, subtopic)
}

func (helper *MqttHelper) onConnect(client MQTT.Client) {
	log.Printf("Connect to %s", *broker)
	token := client.Publish("envoy/connected", byte(qos), true, "online")
	token.WaitTimeout(10 * time.Second)
}

func (helper *MqttHelper) PublishRetained(subtopic, message string) {
	token := helper.client.Publish(helper.topic(subtopic), byte(qos), true, message)
	if token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishRetained failed err:%v", token.Error())
	}
}

func (helper *MqttHelper) Publish(subtopic string, value any) {
	var message string
	switch val := value.(type) {
	case string:
		message = val
	case int:
		message = fmt.Sprintf("%d", val)
	case uint:
		message = fmt.Sprintf("%d", val)
	case float32:
		message = fmt.Sprintf("%f", val)
	case float64:
		message = fmt.Sprintf("%f", val)
	}
	token := helper.client.Publish(helper.topic(subtopic), byte(qos), false, message)
	if token.WaitTimeout(1 * time.Second) {
		log.Printf("Publish failed err:%v", token.Error())
	}
}

func (helper *MqttHelper) PublishJson(subtopic string, payload any) {
	message, _ := json.Marshal(payload)
	token := helper.client.Publish(helper.topic(subtopic), byte(qos), false, string(message))
	if token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishJson failed err:%v", token.Error())
	}
}
