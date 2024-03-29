package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "sensu-handler-discord-")
	defer func() {
		_ = os.Remove(file.Name())
	}()

	event := types.FixtureEvent("entity1", "check1")
	eventJSON, _ := json.Marshal(event)
	_, err := file.WriteString(string(eventJSON))
	require.NoError(t, err)
	require.NoError(t, file.Sync())
	_, err = file.Seek(0, 0)
	require.NoError(t, err)
	os.Stdin = file
	requestReceived := false

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"ok": true}`))
		require.NoError(t, err)
	}))

	oldArgs := os.Args
	os.Args = []string{"discord", "-w", apiStub.URL}
	defer func() { os.Args = oldArgs }()

	main()
	assert.True(t, requestReceived)
}
