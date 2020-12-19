package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedrobertao/backend-hackathon-klever-2020/database"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	router = gin.Default()
	router.GET("/user/:address", GetUserByAddress)
	router.GET("/user", GetUser)
	router.POST("/user/phone", PhoneVerify)
	router.PUT("/user", CreateUser)
	router.POST("/sms/transaction", SmsTransaction)

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(".env not found, using os variables")
	}
	database.Connect()

	code := m.Run()
	os.Exit(code)
}

func TestGetUserByAddress(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/user/bertao", nil)
	router.ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response interface{}
	err = json.Unmarshal(body, &response)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(resp.StatusCode, http.StatusOK)
}

func TestGetUser(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/user?search=roney", nil)
	router.ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response interface{}
	err = json.Unmarshal(body, &response)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(resp.StatusCode, http.StatusOK)
}

func TestPhoneVerify(t *testing.T) {
	t.Skip()
	w := httptest.NewRecorder()
	data := fmt.Sprintf(`{"username": "bertao"}`)
	req, err := http.NewRequest("POST", "/user/phone", strings.NewReader(data))
	router.ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response interface{}
	err = json.Unmarshal(body, &response)
	fmt.Println(response)
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(resp.StatusCode, http.StatusOK)
}

func TestCreateUser(t *testing.T) {
	t.Skip()
	w := httptest.NewRecorder()
	data := fmt.Sprintf(`
		{["addresses": "TYTRqwRGkLGWiUnb4TNW3Pyk5aLbtq438h]", 
		"mainAddress": "TYTRqwRGkLGWiUnb4TNW3Pyk5aLbtq438h",
		"username": "fulano",
		"email": "fulano@klever.io",
		"phone": "5585999997766",}
		`)
	req, err := http.NewRequest("PUT", "/user", strings.NewReader(data))
	router.ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response interface{}
	err = json.Unmarshal(body, &response)
	fmt.Println(response)
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(resp.StatusCode, http.StatusOK)
}

func TestSmsTransaction(t *testing.T) {
	t.Skip()
	w := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"from": "bertao", 
		"to": "roney", 
		"coin": "BTC", 
		"amount": "1"
		}`)
	req, err := http.NewRequest("POST", "/user/phone", strings.NewReader(data))
	router.ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response interface{}
	err = json.Unmarshal(body, &response)
	fmt.Println(response)
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(resp.StatusCode, http.StatusOK)
}
