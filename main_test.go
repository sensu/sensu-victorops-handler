package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	assert.Error(CheckArgs(event))
	config.RoutingKey = "123"
	assert.Error(CheckArgs(event))
	config.APIURL = "InvalidURL"
	assert.Error(CheckArgs(event))
	config.APIURL = "https://alert.victorops.com/integrations/generic/20131114/alert"
	assert.NoError(CheckArgs(event))
}

func TestSendVictorOps(t *testing.T) {
	testcases := []struct {
		status  uint32
		msgtype string
	}{
		{0, "RECOVERY"},
		{1, "WARNING"},
		{2, "CRITICAL"},
		{127, "CRITICAL"},
	}

	for _, tc := range testcases {
		assert := assert.New(t)
		event := corev2.FixtureEvent("entity1", "check1")
		event.Check.Status = tc.status

		var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(err)
			msg := &VOEvent{}
			err = json.Unmarshal(body, msg)
			require.NoError(t, err)
			expectedEntityID := "entity1/check1"
			assert.Equal(expectedEntityID, msg.EntityID)
			expectedMessageType := tc.msgtype
			assert.Equal(expectedMessageType, msg.MessageType)
			resp := VOResponse{
				Result:   "success",
				EntityID: msg.EntityID,
			}
			respj, err := json.Marshal(resp)
			require.NoError(t, err)
			_, err = w.Write(respj)
			require.NoError(t, err)
		}))

		_, err := url.ParseRequestURI(test.URL)
		require.NoError(t, err)
		config.APIURL = test.URL
		config.RoutingKey = "123"
		config.MessageTemplate = "{{.Entity.Name}}:{{.Check.Name}}:{{.Check.Output}}"
		config.EntityIDTemplate = "{{.Entity.Name}}/{{.Check.Name}}"
		assert.NoError(SendVictorOps(event))
	}
}
