package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentController struct {
	Db *gorm.DB
	Validate *validator.Validate
}

func NewPaymentControllerImpl(Db *gorm.DB, validate *validator.Validate) *PaymentController {
	return &PaymentController{Db: Db, Validate: validate}
}

type ResponseData struct {
    Data struct {
		AccessToken string `json:"accessToken"`
		ExpiresAt int `json:"expiresAt"`
		ExpiresIn int `json:"expiresIn"`
		TokenType string `json:"tokenType"`
    } `json:"data"`
}

func (c PaymentController) GenAccessToken(ctx *gin.Context) {
	client := &http.Client{}

	u := uuid.New()

	jsonData := []byte(`{
		"applicationKey" : "",
		"applicationSecret" : ""
	}`);
	req, err := http.NewRequest("POST", "https://api-sandbox.partners.scb/partners/sandbox/v1/oauth/token", bytes.NewBuffer(jsonData))
	if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

	req.Header.Set("resourceOwnerId", "")
	req.Header.Set("requestUId", u.String())
	req.Header.Set("accept-language", "EN")
    req.Header.Set("Content-Type", "application/json")


	// Send the request
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return
    }

	//Extract data from the response (assume it's JSON)
	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	// generating the qr code for payment response
	jsonData = []byte(`{
		"qrType": "PP", 
		"ppType": "BILLERID", 
		"ppId": "937520472133613", 
		"amount": "5000", 
		"ref1": "DDWDWDDW", 
		"ref2": "SSQSQSS", 
		"ref3": "SCB" 
	}`);
	req, err = http.NewRequest("POST", "https://api-sandbox.partners.scb/partners/sandbox/v1/payment/qrcode/create", bytes.NewBuffer(jsonData))
	if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

	req.Header.Set("authorization", "Bearer " + responseData.Data.AccessToken)
	req.Header.Set("resourceOwnerId", "")
	req.Header.Set("requestUId", u.String())
	req.Header.Set("accept-language", "EN")
    req.Header.Set("Content-Type", "application/json")


	// Send the request
    resp, err = client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return
    }

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "generated access token successfully",
		"body": response,
	})
}