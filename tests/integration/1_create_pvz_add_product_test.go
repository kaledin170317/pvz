package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"pvZ/init/migrations"
	"pvZ/internal/adapters/api/rest"
	"pvZ/internal/app"
	"pvZ/internal/config"
	"pvZ/internal/logger"
	"testing"
)

func TestCreatePVZAddProduct(t *testing.T) {
	cfg := config.Load()
	logger.Init()
	fmt.Println("DSN:", cfg.DB.DSN())

	dbx := sqlx.MustConnect("postgres", cfg.DB.DSN())
	defer dbx.Close()

	db := dbx.DB

	if err := migrations.RollbackMigrations(db); err != nil {
		log.Println(err)
	}

	if err := migrations.RunMigrations(db); err != nil {
		log.Println(err)
	}

	secretKey := []byte(cfg.App.JWTSecret)

	deps := app.SetupDependencies(dbx, secretKey)

	router := app.SetupRoutes(deps.UserUC, deps.PVZUC, deps.ReceptionUC, deps.ProductUC, secretKey)

	if err := migrations.RollbackMigrations(db); err != nil {
		assert.NoError(t, err)
	}

	if err := migrations.RunMigrations(db); err != nil {
		assert.NoError(t, err)
	}

	testServer := httptest.NewServer(router)
	defer testServer.Close()
	testClient := testServer.Client()

	moderatorToken := getToken(t, testClient, testServer.URL, rest.DummyLoginRequest{Role: "moderator"})
	employeeToken := getToken(t, testClient, testServer.URL, rest.DummyLoginRequest{Role: "employee"})

	// 1. Создание ПВЗ
	pvzReq := rest.CreatePVZRequest{
		ID:               "11111111-1111-1111-1111-111111111111",
		City:             "Москва",
		RegistrationDate: "2025-04-18T21:00:00Z",
	}
	resp := doRequest(t, testClient, "POST", testServer.URL+"/pvz", moderatorToken, pvzReq)
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
	resp = doRequest(t, testClient, "POST", testServer.URL+"/receptions", employeeToken, recReq)
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

		resp = doRequest(t, testClient, "POST", testServer.URL+"/products", employeeToken, prodReq)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var prodResp rest.ProductResponse
		err = json.NewDecoder(resp.Body).Decode(&prodResp)
		assert.NoError(t, err)
		assert.Equal(t, tmpType, prodResp.Type)
		assert.Equal(t, recResp.ID, prodResp.ReceptionID)
		resp.Body.Close()
	}

	// 4. Закрытие приёмки
	resp = doRequest(t, testClient, "POST", fmt.Sprintf("%s/pvz/%s/close_last_reception", testServer.URL, pvzResp.ID), employeeToken, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var closedRec rest.ReceptionResponse
	err = json.NewDecoder(resp.Body).Decode(&closedRec)
	assert.NoError(t, err)
	assert.Equal(t, "close", closedRec.Status)
	resp.Body.Close()
}
