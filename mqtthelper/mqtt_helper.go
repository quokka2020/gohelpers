package mqtthelper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
var verbose = util.GetEnvBool("MQTT_VERBOSE", false)

// When the cleanSession flag is set to true, the client explicitly requests a non-persistent session.
var clean_session = util.GetEnvBool("MQTT_CLEAN", true)

func mqtt_id() string {
	filename := filepath.Base(os.Args[0])
	return fmt.Sprintf("%s-%s", filename, util.GetEnv("HOSTNAME", "localhost"))
}

type Mqtt_Helper struct {
	client               MQTT.Client
	Prefix               string // also the name
	onConnectHandlers    []func(helper *Mqtt_Helper)
	onDisconnectHandlers []func()
	numberMapping        map[string]func(string, float64)
	numberMappingFull    map[string]func(string, float64)
	stringMapping        map[string]func(string, string)
	stringMappingFull    map[string]func(string, string)
}

func CreateMqttHelper(prefix string) *Mqtt_Helper {
	helper := Mqtt_Helper{
		Prefix:               prefix,
		onConnectHandlers:    make([]func(helper *Mqtt_Helper), 0),
		onDisconnectHandlers: make([]func(), 0),
		numberMapping:        make(map[string]func(string, float64)),
		numberMappingFull:    make(map[string]func(string, float64)),
		stringMapping:        make(map[string]func(string, string)),
		stringMappingFull:    make(map[string]func(string, string)),
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
	opts.SetConnectionLostHandler(helper.onDisconnect)
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
	token := helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.numberReceived)
	if verbose && !token.WaitTimeout(time.Second) {
		log.Printf("MQTT_HELPER AddNumberSubscription subtopic:%s err:%v", subtopic, token.Error())
	}
}

func (helper *Mqtt_Helper) AddStringSubscription(subtopic string, function func(string, string)) {
	helper.stringMapping[subtopic] = function
	token := helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.stringReceived)
	if verbose && !token.WaitTimeout(time.Second) {
		log.Printf("MQTT_HELPER AddStringSubscription subtopic:%s err:%v", subtopic, token.Error())
	}
}

func (helper *Mqtt_Helper) AddNumberSubscriptionFull(topic string, function func(string, float64)) {
	helper.numberMappingFull[topic] = function
	token := helper.client.Subscribe(topic, byte(0), helper.numberReceivedFull)
	if verbose && !token.WaitTimeout(time.Second) {
		log.Printf("MQTT_HELPER AddNumberSubscriptionFull subtopic:%s err:%v", topic, token.Error())
	}
}

func (helper *Mqtt_Helper) AddStringSubscriptionFull(topic string, function func(string, string)) {
	helper.stringMappingFull[topic] = function
	token := helper.client.Subscribe(topic, byte(0), helper.stringReceivedFull)
	if verbose && !token.WaitTimeout(time.Second) {
		log.Printf("MQTT_HELPER AddStringSubscriptionFull subtopic:%s err:%v", topic, token.Error())
	}
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
	log.Printf("MQTT_HELPER Connected to %s", broker)
	helper.PublishRetained("connected", "1")
	for subtopic := range helper.numberMapping {
		if token := helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.numberReceived); token.Wait() && token.Error() != nil {
			log.Printf("MQTT_HELPER failed to subscribe to %s err: %v", subtopic, token.Error())
		}
	}
	for subtopic := range helper.stringMapping {
		if token := helper.client.Subscribe(helper.topic(subtopic), byte(0), helper.stringReceived); token.Wait() && token.Error() != nil {
			log.Printf("MQTT_HELPER failed to subscribe to %s err: %v", subtopic, token.Error())
		}
	}
	for topic := range helper.numberMappingFull {
		if token := helper.client.Subscribe(topic, byte(0), helper.numberReceivedFull); token.Wait() && token.Error() != nil {
			log.Printf("MQTT_HELPER failed to subscribe to %s err: %v", topic, token.Error())
		}
	}
	for topic := range helper.stringMappingFull {
		if token := helper.client.Subscribe(topic, byte(0), helper.stringReceivedFull); token.Wait() && token.Error() != nil {
			log.Printf("MQTT_HELPER failed to subscribe to %s err: %v", topic, token.Error())
		}
	}
	for _, onConnectHandler := range helper.onConnectHandlers {
		onConnectHandler(helper)
	}
}

func (helper *Mqtt_Helper) onDisconnect(c MQTT.Client, err error) {
	log.Printf("just lost mqtt connection err:%v", err)
	for _, onDisconnectHandler := range helper.onDisconnectHandlers {
		onDisconnectHandler()
	}
}

func (helper *Mqtt_Helper) RegisterOnConnectHandler(onConnectHandler func(*Mqtt_Helper)) {
	helper.onConnectHandlers = append(helper.onConnectHandlers, onConnectHandler)
}

func (helper *Mqtt_Helper) UnregisterOnConnectHandler(toRemove func(*Mqtt_Helper)) {
	helper.onConnectHandlers = slices.DeleteFunc(helper.onConnectHandlers, func(cur func(*Mqtt_Helper)) bool {
		cur_val := reflect.ValueOf(cur)
		toRemove_val := reflect.ValueOf(toRemove)
		return cur_val == toRemove_val
	})
}

func (helper *Mqtt_Helper) RegisterOnDisconnectHandler(onDisconnectHandler func()) {
	helper.onDisconnectHandlers = append(helper.onDisconnectHandlers, onDisconnectHandler)
}

func (helper *Mqtt_Helper) UnregisterOnDisconnectHandler(toRemove func()) {
	helper.onDisconnectHandlers = slices.DeleteFunc(helper.onDisconnectHandlers, func(cur func()) bool {
		cur_val := reflect.ValueOf(cur)
		toRemove_val := reflect.ValueOf(toRemove)
		return cur_val == toRemove_val
	})
}

func (helper *Mqtt_Helper) PublishFullRetained(topic, message string) {
	token := helper.client.Publish(topic, byte(qos), true, message)
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("MQTT_HELPER PublishRetained failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) PublishRetained(subtopic, message string) {
	helper.PublishFullRetained(helper.topic(subtopic), message)
}

func ValueToMessage(value any) []byte {
	var message string
	if value == nil {
		return nil
	}
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
		log.Printf("MQTT_HELPER Type not implemented for topic %v", value)
		message = fmt.Sprintf("%v", value)
	}
	return []byte(message)
}

func (helper *Mqtt_Helper) PublishFull(topic string, value any) {
	message := ValueToMessage(value)
	if util.Verbose() {
		log.Printf("MQTT_HELPER publish token:%s message:%s", topic, string(message))
	}
	token := helper.client.Publish(topic, byte(qos), false, message)
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("MQTT_HELPER Publish failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) Publish(subtopic string, value any) {
	helper.PublishFull(helper.topic(subtopic), value)
}

func (helper *Mqtt_Helper) PublishFullJson(topic string, payload any) {
	message, _ := json.Marshal(payload)
	token := helper.client.Publish(topic, byte(qos), false, string(message))
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("MQTT_HELPER PublishJson failed err:%v", token.Error())
	}
}

func (helper *Mqtt_Helper) PublishJson(subtopic string, payload any) {
	helper.PublishFullJson(helper.topic(subtopic), payload)
}

func (helper *Mqtt_Helper) numberReceived(client MQTT.Client, msg MQTT.Message) {
	if verbose {
		log.Printf("MQTT_HELPER Number received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	}

	i, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		log.Printf("MQTT_HELPER Got a non-number from %s with payload [%s] %v", msg.Topic(), string(msg.Payload()), err)
		return
	}
	for subtopic, function := range helper.numberMapping {
		if match(msg.Topic(), helper.topic(subtopic)) {
			function(helper.subtopic(msg.Topic()), i)
			return
		}
	}
	log.Printf("MQTT_HELPER Got an unmapped number from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func (helper *Mqtt_Helper) stringReceived(client MQTT.Client, msg MQTT.Message) {
	if verbose {
		log.Printf("MQTT_HELPER String received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	}
	for subtopic, function := range helper.stringMapping {
		if match(msg.Topic(), helper.topic(subtopic)) {
			function(helper.subtopic(msg.Topic()), string(msg.Payload()))
			return
		}
	}
	log.Printf("MQTT_HELPER Got an unmapped string from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func (helper *Mqtt_Helper) numberReceivedFull(client MQTT.Client, msg MQTT.Message) {
	if verbose {
		log.Printf("MQTT_HELPER Number received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	}

	i, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		log.Printf("MQTT_HELPER Got a non-number from %s with payload [%s] %v", msg.Topic(), string(msg.Payload()), err)
		return
	}
	for topic, function := range helper.numberMappingFull {
		if match(msg.Topic(), topic) {
			function(msg.Topic(), i)
			return
		}
	}
	log.Printf("MQTT_HELPER Got an unmapped number from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
}

func (helper *Mqtt_Helper) stringReceivedFull(client MQTT.Client, msg MQTT.Message) {
	if verbose {
		log.Printf("MQTT_HELPER String received %s with payload:[%s]", msg.Topic(), string(msg.Payload()))
	}
	for topic, function := range helper.stringMappingFull {
		if match(msg.Topic(), topic) {
			function(msg.Topic(), string(msg.Payload()))
			return
		}
	}
	log.Printf("MQTT_HELPER Got an unmapped string from %s with payload [%s]", msg.Topic(), string(msg.Payload()))
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
