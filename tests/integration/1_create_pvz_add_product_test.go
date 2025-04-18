package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"pvZ/init/migrations"
	"pvZ/internal/adapters/api/rest"
	"testing"
)

func TestCreatePVZAddProduct(t *testing.T) {

	if err := migrations.RollbackMigrations(database); err != nil {
		assert.NoError(t, err)
	}

	if err := migrations.RunMigrations(database); err != nil {
		assert.NoError(t, err)
	}

	server := testServer
	client := testClient

	moderatorToken := getToken(t, client, server.URL, rest.DummyLoginRequest{Role: "moderator"})
	employeeToken := getToken(t, client, server.URL, rest.DummyLoginRequest{Role: "employee"})

	// 1. Создание ПВЗ
	pvzReq := rest.CreatePVZRequest{
		ID:               "11111111-1111-1111-1111-111111111111",
		City:             "Москва",
		RegistrationDate: "2025-04-18T21:00:00Z",
	}
	resp := doRequest(t, client, "POST", server.URL+"/pvz", moderatorToken, pvzReq)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var pvzResp rest.PVZResponse
	err := json.NewDecoder(resp.Body).Decode(&pvzResp)
	assert.NoError(t, err)
	assert.Equal(t, pvzReq.City, pvzResp.City)
	assert.Equal(t, pvzReq.ID, pvzResp.ID)
	assert.Equal(t, pvzReq.RegistrationDate, pvzResp.RegistrationDate)
	resp.Body.Close()

	// 2. Создание приёмки
	recReq := rest.CreateReceptionRequest{PVZID: pvzResp.ID}
	resp = doRequest(t, client, "POST", server.URL+"/receptions", employeeToken, recReq)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var recResp rest.ReceptionResponse
	err = json.NewDecoder(resp.Body).Decode(&recResp)
	assert.Equal(t, recReq.PVZID, recResp.PVZID)
	assert.NoError(t, err)
	resp.Body.Close()

	// 3. Добавление 50 товаров
	for i := 0; i < 1; i++ {
		tmpType := fmt.Sprintf("электроника")

		if i%2 == 0 {
			tmpType = fmt.Sprintf("электроника")
		} else {
			tmpType = fmt.Sprintf("обувь")
		}
		prodReq := rest.AddProductRequest{
			Type:  tmpType,
			PVZID: pvzResp.ID,
		}

		resp = doRequest(t, client, "POST", server.URL+"/products", employeeToken, prodReq)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var prodResp rest.ProductResponse
		err = json.NewDecoder(resp.Body).Decode(&prodResp)
		assert.NoError(t, err)
		assert.Equal(t, tmpType, prodResp.Type)
		assert.Equal(t, recResp.ID, prodResp.ReceptionID)
		resp.Body.Close()
	}

	// 4. Закрытие приёмки
	resp = doRequest(t, client, "POST", fmt.Sprintf("%s/pvz/%s/close_last_reception", server.URL, pvzResp.ID), employeeToken, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var closedRec rest.ReceptionResponse
	err = json.NewDecoder(resp.Body).Decode(&closedRec)
	assert.NoError(t, err)
	assert.Equal(t, "close", closedRec.Status)
	resp.Body.Close()
}
