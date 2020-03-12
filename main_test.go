package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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
	config.APIURL = "http://sensu.example.com:3000"
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
			w.WriteHeader(http.StatusOK)
		}))

		_, err := url.ParseRequestURI(test.URL)
		require.NoError(t, err)
		config.APIURL = test.URL
		config.RoutingKey = "123"
		assert.NoError(SendVictorOps(event))
	}
}

func Testmain(t *testing.T) {
	assert := assert.New(t)
	file, _ := ioutil.TempFile(os.TempDir(), "sensu-victorops-handler")
	defer func() {
		_ = os.Remove(file.Name())
	}()

	event := corev2.FixtureEvent("entity1", "check1")
	event.Metrics = corev2.FixtureMetrics()
	eventJSON, _ := json.Marshal(event)
	_, err := file.WriteString(string(eventJSON))
	require.NoError(t, err)
	require.NoError(t, file.Sync())
	_, err = file.Seek(0, 0)
	require.NoError(t, err)
	os.Stdin = file
	requestReceived := false

	var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"ok": true}`))
		require.NoError(t, err)
	}))

	_, err = url.ParseRequestURI(test.URL)
	require.NoError(t, err)
	oldArgs := os.Args
	os.Args = []string{"sensu-victorops-handler", "--api-url", test.URL, "--routingkey", "123"}
	defer func() { os.Args = oldArgs }()

	main()
	assert.True(requestReceived)
}
