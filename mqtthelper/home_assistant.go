package mqtthelper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type ha_device struct {
	Model       string `json:"model"`
	Name        string `json:"name"`
	SwVersion   string `json:"sw_version"`
	Identifiers string `json:"identifiers"`
}

type ha_config struct {
	// Availability []ha_availability `json:"availability"`
	Availability_Topic string    `json:"availability_topic"`
	Device             ha_device `json:"device"`
	Name               string    `json:"name"`
	StateTopic         string    `json:"state_topic"`
	UnitOfMeasurement  string    `json:"unit_of_measurement"`
	ValueTemplate      string    `json:"value_template,omitempty"`
	UniqueId           string    `json:"unique_id"`
	StateClass         string    `json:"state_class,omitempty"`
	DeviceClass        string    `json:"device_class"`
	Icon               string    `json:"icon,omitempty"`
}

func (helper *MqttHelper) base_ha_config() ha_config {
	return ha_config{
		Availability_Topic: helper.topic("connected"),
		Device: ha_device{
			Model:       helper.Prefix,
			Name:        helper.Prefix,
			SwVersion:   "UNKNOWN",
			Identifiers: helper.Prefix,
		},
	}
}

func (helper *MqttHelper) HARegisterIncreasing(subtopic string, name string, unit string, device_class string, icon string) {
	state_topic := helper.topic(subtopic)
	topic := fmt.Sprintf("homeassistant/sensor/%s/config", state_topic)
	payload := helper.base_ha_config()
	payload.Name = name
	payload.StateTopic = state_topic
	payload.UnitOfMeasurement = unit
	payload.UniqueId = name
	payload.StateClass = "total_increasing"
	payload.DeviceClass = device_class
	payload.Icon = icon

	msg, _ := json.Marshal(payload)
	// log.Printf("Should publish %s %s", topic, msg)
	token := helper.client.Publish(topic, byte(qos), true, string(msg))
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishRetained failed err:%v", token.Error())
	}
}

func (helper *MqttHelper) HARegister(subtopic string, name string, unit string, device_class string, icon string) {
	state_topic := helper.topic(subtopic)
	topic := fmt.Sprintf("homeassistant/sensor/%s/config", state_topic)
	payload := helper.base_ha_config()
	payload.Name = name
	payload.StateTopic = state_topic
	payload.UnitOfMeasurement = unit
	payload.UniqueId = name
	payload.DeviceClass = device_class
	payload.Icon = icon

	msg, _ := json.Marshal(payload)
	// log.Printf("Should publish %s %s", topic, msg)
	token := helper.client.Publish(topic, byte(qos), true, string(msg))
	if !token.WaitTimeout(1 * time.Second) {
		log.Printf("PublishRetained failed err:%v", token.Error())
	}
}
