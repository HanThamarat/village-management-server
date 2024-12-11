package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gorm.io/gorm"
	"tectnexify.github.com/e-payment/data/request"
	_"tectnexify.github.com/e-payment/models"
)

type OPNController struct {
	client *omise.Client
	Db *gorm.DB
	Validate *validator.Validate
}

func NewOPNControllerImpl(Db *gorm.DB, validate *validator.Validate, client *omise.Client) *OPNController {
	return &OPNController{
		Db: Db,
		Validate: validate,
		client: client,
	}
}

func (c OPNController) CreateTokenOPN(ctx *gin.Context) {
    var reqBody request.CreateTokenOPN;
	if err := ctx.Bind(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationError := err.(validator.ValidationErrors)
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationError[0].Field())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	result := &omise.Card{}

	err := c.client.Do(result, &operations.CreateToken{
		Name: reqBody.Name,
		Number: reqBody.Number,
		ExpirationMonth: time.Month(reqBody.ExpirationMonth),
		ExpirationYear: int(reqBody.ExpirationYear),
		City: reqBody.City,
		PostalCode: reqBody.PostalCode,
		SecurityCode: reqBody.SecurityCode,
	})

	if err != nil {
		log.Fatal(err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to create token",
        });
        return;
	}

	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": "Create token successfully",
		"body": result,
	});
};

