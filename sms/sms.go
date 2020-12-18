package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pedrobertao/backend-hackathon-klever-2020/models"
)

var (
	accountSid        = ""
	authToken         = ""
	accountServiceSid = ""
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func Config() error {
	accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
	authToken = os.Getenv("TWILIO_AUTH_TOKEN")
	accountServiceSid = os.Getenv("TWILIO_SERVICE_SID")
	return nil
}

// VerifyCodeSMS ...
func VerifyCodeSMS(phone string, code string) error {
	msgData := url.Values{}
	msgData.Set("To", phone)
	msgData.Set("Code", code)

	urlStr := fmt.Sprintf("https://verify.twilio.com/v2/Services/%s/VerificationCheck", accountServiceSid)

	msgDataReader := strings.NewReader(msgData.Encode())
	req, _ := http.NewRequest("POST", urlStr, msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("Request failed with %+v", data)
}

// SendVerifySMS ...
func SendVerifySMS(phone string) error {
	msgData := url.Values{}
	msgData.Set("To", phone)
	msgData.Set("Channel", "sms")

	urlStr := fmt.Sprintf("https://verify.twilio.com/v2/Services/%s/Verifications", os.Getenv("TWILIO_SERVICE_SID"))

	msgDataReader := strings.NewReader(msgData.Encode())
	req, _ := http.NewRequest("POST", urlStr, msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("Request failed with %+v", data)
}

//SendSMS send sms
func SendSMS(sms models.SMS) error {
	// Set account keys & information
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", sms.To)
	msgData.Set("From", sms.From)
	msgData.Set("Body", sms.Body)
	msgDataReader := strings.NewReader(msgData.Encode())

	req, _ := http.NewRequest("POST", urlStr, msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("Request failed with %+v", data)

}
