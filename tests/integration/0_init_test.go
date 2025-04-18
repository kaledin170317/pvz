package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"pvZ/internal/adapters/api/rest"
	"pvZ/internal/app"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testServer *httptest.Server
	testClient *http.Client
	databaseX  *sqlx.DB
	database   *sql.DB
)

func TestMain(m *testing.M) {
	fmt.Println("Старт тестов")
	dbx := sqlx.MustConnect("postgres", "postgres://postgres:password@localhost:55555/pvz?sslmode=disable")
	defer dbx.Close()

	database = dbx.DB
	databaseX = dbx

	secretKey := []byte("test")
	router := app.SetupRouter(dbx, secretKey)

	testServer = httptest.NewServer(router)
	defer testServer.Close()

	testClient = testServer.Client()

	// Запуск всех тестов
	code := m.Run()
	os.Exit(code)
}

func getToken(t *testing.T, client *http.Client, baseURL string, login rest.DummyLoginRequest) string {
	body, _ := json.Marshal(login)
	resp, err := client.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	tokenBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	return string(bytes.Trim(tokenBytes, `"`)) // JWT в виде строки
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
