package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu-community/sensu-plugin-sdk/templates"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

// VOEvent is the JSON type for a VictorOps message
type VOEvent struct {
	MessageType    string `json:"message_type"`
	StateMessage   string `json:"state_message,omitempty"`
	EntityID       string `json:"entity_id,omitempty"`
	HostName       string `json:"host_name,omitempty"`
	MonitoringTool string `json:"monitoring_tool,omitempty"`
	// This is feature parity with Ruby sensu-plugins-victorops
	// which includes the check and client in its payload.
	Check  *corev2.Check  `json:"check,omitempty"`
	Entity *corev2.Entity `json:"entity,omitempty"`
}

// VOResopnse is the JSON type for a VictorOps response message
type VOResponse struct {
	Result   string `json:"result"`
	EntityID string `json:"entity_id"`
}

// HandlerConfig is needed for Sensu Go Handlers
type HandlerConfig struct {
	sensu.PluginConfig
	RoutingKey       string
	APIURL           string
	MessageTemplate  string
	EntityIDTemplate string
}

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-victorops-handler",
			Short:    "The Sensu Go VictorOps handler for sending notifications",
			Keyspace: "sensu.io/plugins/victorops/config",
		},
	}
	// VictorOpsConfigOptions contains the Sensu Plugin Config Options
	VictorOpsConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      "routingkey",
			Env:       "SENSU_VICTOROPS_ROUTINGKEY",
			Argument:  "routingkey",
			Shorthand: "r",
			Default:   "",
			Secret:    true,
			Usage:     "The VictorOps Routing Key",
			Value:     &config.RoutingKey,
		},
		{
			Path:      "api-url",
			Env:       "SENSU_VICTOROPS_APIURL",
			Argument:  "api-url",
			Shorthand: "a",
			Default:   "https://alert.victorops.com/integrations/generic/20131114/alert",
			Usage:     "The URL for the VictorOps API",
			Value:     &config.APIURL,
		},
		{
			Path:      "message-template",
			Env:       "",
			Argument:  "message-template",
			Shorthand: "m",
			Default:   "{{.Entity.Name}:{{.Check.Name}}:{{.Check.Output}}",
			Usage:     "The template for the message sent to VictorOps",
			Value:     &config.MessageTemplate,
		},
		{
			Path:      "entity-id-template",
			Env:       "",
			Argument:  "entity-id-template",
			Shorthand: "e",
			Default:   "{{.Entity.Name}/{{.Check.Name}}",
			Usage:     "The template for the Entity ID sent to VictorOps",
			Value:     &config.EntityIDTemplate,
		},
	}
)

func main() {

	goHandler := sensu.NewGoHandler(&config.PluginConfig, VictorOpsConfigOptions, CheckArgs, SendVictorOps)
	goHandler.Execute()

}

// CheckArgs checks that necessary arguments are set
func CheckArgs(_ *corev2.Event) error {

	if len(config.RoutingKey) == 0 {
		return errors.New("missing VictorOps Routing Key")
	}
	if len(config.APIURL) == 0 {
		return errors.New("missing VictorOps API URL specification")
	}
	if !govalidator.IsURL(config.APIURL) {
		return errors.New("invlaid VictorOps API URL specification")
	}
	config.APIURL = strings.TrimSuffix(config.APIURL, "/")

	return nil
}

// SendVictorOps builds the event message and sends it to VO
func SendVictorOps(event *corev2.Event) error {

	var msgType string

	switch eventStatus := event.Check.Status; eventStatus {
	case 0:
		msgType = "RECOVERY"
	case 1:
		msgType = "WARNING"
	default:
		msgType = "CRITICAL"
	}

	msgEntityID, err := templates.EvalTemplate("EntityID", config.EntityIDTemplate, event)
	if err != nil {
		return fmt.Errorf("Failed to create Entity ID from template: %v", err)
	}
	msgStateMessage, err := templates.EvalTemplate("Message", config.MessageTemplate, event)
	if err != nil {
		return fmt.Errorf("Failed to create message from template: %v", err)
	}

	message := VOEvent{
		MessageType:    msgType,
		StateMessage:   msgStateMessage,
		EntityID:       msgEntityID,
		HostName:       event.Entity.Name,
		MonitoringTool: "sensu",
		Check:          event.Check,
		Entity:         event.Entity,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Failed to marshal VictorOps message: %s", err)
	}

	url := fmt.Sprintf("%s/%s", config.APIURL, config.RoutingKey)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("Post to %s failed: %s", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("POST to %s failed with %v", url, resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body from %s: %v", config.APIURL, err)
	}

	voResponse := VOResponse{}
	err = json.Unmarshal(body, &voResponse)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response from VictorOps: %v", err)
	}

	// FUTURE: send to AH
	// VictorOps only sends back a "success" or "failure" message along with the Entity ID
	fmt.Printf("Submission to VictorOps for Entity ID %s result: %s\n", voResponse.EntityID, voResponse.Result)

	// Should we actually return error if Result above is "failure"?
	return nil
}
