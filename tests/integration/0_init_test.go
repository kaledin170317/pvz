package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"pvZ/internal/adapters/api/rest"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "pvz")
	os.Setenv("APP_PORT", "1488")
	os.Setenv("JWT_SECRET", "test-secret")
}

func getToken(t *testing.T, client *http.Client, baseURL string, login rest.DummyLoginRequest) string {
	body, _ := json.Marshal(login)
	resp, err := client.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	tokenBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	return string(bytes.Trim(tokenBytes, `"`))
}

func doRequest(t *testing.T, client *http.Client, method, url, token string, body interface{}) *http.Response {
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		assert.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, &buf)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}
