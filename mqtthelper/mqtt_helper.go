package mqtthelper

import (
	"encoding/json"
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
var broker = util.GetEnv("MQTT_BROKER", "tcp://192.168.10.4:1883")
var password = util.GetEnv("MQTT_PASSWD", "")
var user = util.GetEnv("MQTT_USER", "")
var id = util.GetEnv("MQTT_ID", mqtt_id())
var clean_session = util.GetEnvBool("MQTT_CLEAN", false)

func mqtt_id() string {
	filename := filepath.Base(os.Args[0])
	return fmt.Sprintf("%s-%s", filename, util.GetEnv("HOSTNAME", "localhost"))
}

type Mqtt_Helper struct {
	client MQTT.Client
	Prefix string // also the name
}

func CreateMqttHelper(prefix string) (*Mqtt_Helper, error) {
	helper := Mqtt_Helper{
		Prefix: prefix,
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(clean_session)
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
		return nil, token.Error()
	}

	return &helper, nil
}

func (helper *Mqtt_Helper) Close() {
	helper.client.Disconnect(250)
}

func (helper *Mqtt_Helper) topic(subtopic string) string {
	return fmt.Sprintf("%s/%s", helper.Prefix, subtopic)
}

func (helper *Mqtt_Helper) onConnect(client MQTT.Client) {
	log.Printf("Connect to %s", broker)
	helper.PublishRetained("connected", "online")
}

func (helper *Mqtt_Helper) PublishRetained(subtopic, message string) {
	token := helper.client.Publish(helper.topic(subtopic), byte(qos), true, message)
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishRetained failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) Publish(subtopic string, value any) {
	topic := helper.topic(subtopic)
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
	case bool:
		if val {
			message = "1"
		} else {
			message = "0"
		}
	default:
		log.Printf("Type not implemented for topic %s",topic)
		message = fmt.Sprintf("%v", value)
	}
	if util.Verbose() {
		log.Printf("mqtt publish token:%s message:%s", topic, message)
	}
	token := helper.client.Publish(topic, byte(qos), false, message)
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("Publish failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) PublishJson(subtopic string, payload any) {
	message, _ := json.Marshal(payload)
	token := helper.client.Publish(helper.topic(subtopic), byte(qos), false, string(message))
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishJson failed err:%v", token.Error())
	}
}
