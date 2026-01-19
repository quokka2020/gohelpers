package mqtthelper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// When the cleanSession flag is set to true, the client explicitly requests a non-persistent session.
var clean_session = util.GetEnvBool("MQTT_CLEAN", true)

func mqtt_id() string {
	filename := filepath.Base(os.Args[0])
	return fmt.Sprintf("%s-%s", filename, util.GetEnv("HOSTNAME", "localhost"))
}

type Mqtt_Helper struct {
	client            MQTT.Client
	Prefix            string // also the name
	numberMapping     map[string]func(string, float64)
	numberMappingFull map[string]func(string, float64)
	stringMapping     map[string]func(string, string)
	stringMappingFull map[string]func(string, string)
}

func CreateMqttHelper(prefix string) *Mqtt_Helper {
	helper := Mqtt_Helper{
		Prefix:            prefix,
		numberMapping:     make(map[string]func(string, float64)),
		numberMappingFull: make(map[string]func(string, float64)),
		stringMapping:     make(map[string]func(string, string)),
		stringMappingFull: make(map[string]func(string, string)),
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
		"0",
		0,
		true,
	)
	opts.SetOnConnectHandler(helper.onConnect)

	// opts.SetDefaultPublishHandler(data.msgReceived)

	helper.client = MQTT.NewClient(opts)
	go helper.client.Connect()

	return &helper
}

func (helper *Mqtt_Helper) GetClient() MQTT.Client {
	return helper.client
}

func (helper *Mqtt_Helper) AddNumberSubscription(subtopic string, function func(string, float64)) {
	helper.numberMapping[subtopic] = function
	helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.numberReceived)
}

func (helper *Mqtt_Helper) AddStringSubscription(subtopic string, function func(string, string)) {
	helper.stringMapping[subtopic] = function
	helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.stringReceived)
}

func (helper *Mqtt_Helper) AddNumberSubscriptionFull(topic string, function func(string, float64)) {
	helper.numberMappingFull[topic] = function
	helper.client.Subscribe(topic, byte(0), helper.numberReceivedFull)
}

func (helper *Mqtt_Helper) AddStringSubscriptionFull(topic string, function func(string, string)) {
	helper.stringMappingFull[topic] = function
	helper.client.Subscribe(topic, byte(0), helper.stringReceivedFull)
}

func (helper *Mqtt_Helper) Close() {
	helper.client.Disconnect(250)
}

func (helper *Mqtt_Helper) topic(subtopic string) string {
	return fmt.Sprintf("%s/%s", helper.Prefix, subtopic)
}

func (helper *Mqtt_Helper) subtopic(topic string) string {
	prefix_len := len(helper.Prefix)
	if prefix_len+1 > len(topic) {
		return topic
	}
	return topic[prefix_len+1:]
}

func (helper *Mqtt_Helper) onConnect(client MQTT.Client) {
	log.Printf("Connected to %s", broker)
	helper.PublishRetained("connected", "1")
	for subtopic := range helper.numberMapping {
		if token := helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.numberReceived); token.Wait() && token.Error() != nil {
			log.Printf("failed to subscribe to %s err: %v",subtopic, token.Error())
		}
	}
	for subtopic := range helper.stringMapping {
		if token := helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.stringReceived); token.Wait() && token.Error() != nil {
			log.Printf("failed to subscribe to %s err: %v",subtopic, token.Error())
		}
	}
	for topic := range helper.numberMappingFull {
		if token := helper.client.Subscribe(topic, byte(0), helper.numberReceivedFull); token.Wait() && token.Error() != nil {
			log.Printf("failed to subscribe to %s err: %v",topic, token.Error())
		}
	}
	for topic := range helper.stringMappingFull {
		if token := helper.client.Subscribe(topic, byte(0), helper.stringReceivedFull); token.Wait() && token.Error() != nil {
			log.Printf("failed to subscribe to %s err: %v",topic, token.Error())
		}
	}
}

func (helper *Mqtt_Helper) PublishFullRetained(topic, message string) {
	token := helper.client.Publish(topic, byte(qos), true, message)
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishRetained failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) PublishRetained(subtopic, message string) {
	helper.PublishFullRetained(helper.topic(subtopic), message)
}

func (helper *Mqtt_Helper) PublishFullTopic(topic string, value any) {
	var message string
	switch val := value.(type) {
	case string:
		message = val
	case int:
		message = fmt.Sprintf("%d", val)
	case int8:
		message = fmt.Sprintf("%d", val)
	case int16:
		message = fmt.Sprintf("%d", val)
	case int32:
		message = fmt.Sprintf("%d", val)
	case int64:
		message = fmt.Sprintf("%d", val)
	case uint:
		message = fmt.Sprintf("%d", val)
	case uint8:
		message = fmt.Sprintf("%d", val)
	case uint16:
		message = fmt.Sprintf("%d", val)
	case uint32:
		message = fmt.Sprintf("%d", val)
	case uint64:
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
		log.Printf("Type not implemented for topic %s", topic)
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

func (helper *Mqtt_Helper) Publish(subtopic string, value any) {
	helper.PublishFullTopic(helper.topic(subtopic), value)
}

func (helper *Mqtt_Helper) PublishFullJson(topic string, payload any) {
	message, _ := json.Marshal(payload)
	token := helper.client.Publish(topic, byte(qos), false, string(message))
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishJson failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) PublishJson(subtopic string, payload any) {
	helper.PublishFullJson(helper.topic(subtopic), payload)
}

func (helper *Mqtt_Helper) numberReceived(client MQTT.Client, msg MQTT.Message) {
	// if *VERBOSE {
	// 	log.Printf("MQTT Number received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	// }

	i, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		log.Printf("Got a non-number from %s with payload [%s] %v", msg.Topic(), string(msg.Payload()), err)
		return
	}
	for subtopic, function := range helper.numberMapping {
		if match(msg.Topic(), helper.topic(subtopic)) {
			function(helper.subtopic(msg.Topic()), i)
			return
		}
	}
	log.Printf("Got an unmapped number from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func (helper *Mqtt_Helper) stringReceived(client MQTT.Client, msg MQTT.Message) {
	// if *VERBOSE {
	// 	log.Printf("MQTT String received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	// }
	for subtopic, function := range helper.stringMapping {
		if match(msg.Topic(), helper.topic(subtopic)) {
			function(helper.subtopic(msg.Topic()), string(msg.Payload()))
			return
		}
	}
	log.Printf("Got an unmapped string from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func (helper *Mqtt_Helper) numberReceivedFull(client MQTT.Client, msg MQTT.Message) {
	// if *VERBOSE {
	// 	log.Printf("MQTT Number received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	// }

	i, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		log.Printf("Got a non-number from %s with payload [%s] %v", msg.Topic(), string(msg.Payload()), err)
		return
	}
	for topic, function := range helper.numberMappingFull {
		if match(msg.Topic(), topic) {
			function(msg.Topic(), i)
			return
		}
	}
	log.Printf("Got an unmapped number from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func (helper *Mqtt_Helper) stringReceivedFull(client MQTT.Client, msg MQTT.Message) {
	// if *VERBOSE {
	// 	log.Printf("MQTT String received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	// }
	for topic, function := range helper.stringMappingFull {
		if match(msg.Topic(), topic) {
			function(msg.Topic(), string(msg.Payload()))
			return
		}
	}
	log.Printf("Got an unmapped string from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func match(topic, full_possible_match string) bool {
	if topic == full_possible_match {
		return true
	}
	topic_p := strings.Split(topic, "/")
	match_p := strings.Split(full_possible_match, "/")
	for nr, p := range match_p {
		if p == "+" {
			continue
		}
		if p == "#" {
			return true
		}
		if len(topic_p) < nr || topic_p[nr] != match_p[nr] {
			return false
		}
	}
	return true
}
