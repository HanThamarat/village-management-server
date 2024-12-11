package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"tectnexify.github.com/e-payment/data/request"
	"tectnexify.github.com/e-payment/models"
)

type SystemController struct {
	Db *gorm.DB
	Validate *validator.Validate
}

func NewSystemController(Db *gorm.DB, validate *validator.Validate) *SystemController {
	return &SystemController{Db: Db, Validate: validate}
}

func (c SystemController) CreateBank(ctx *gin.Context) {
	var reqBody request.BankReq
	if err := ctx.Bind(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors);
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return;
	}

	var existingBank models.BankCredentials
	result := c.Db.Where("api_key = ? or api_secret = ?", reqBody.ApiKey, reqBody.ApiSecret).First(&existingBank);
	if result.RowsAffected > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "API Key or API Secret already registered"})
        return;
    }

	newBank := models.BankCredentials{
		BankName: &reqBody.BankName,
		AppName: &reqBody.AppName,
		API_KEY: &reqBody.ApiKey,
		API_SECRET: &reqBody.ApiSecret,
		Biller_ID: &reqBody.BillerID,
		Merchant_ID: &reqBody.MerchantID,
		Terminal_ID: &reqBody.TerminalID,
	}

	// create a new data to database
	response := c.Db.Create(&newBank);

	if response.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"messsage": "Failed to create bank"})
        return
	}

	ctx.JSON(http.StatusOK, gin.H{"messsage": "created bank successfully"})
}

func (c SystemController) GetBanks(ctx *gin.Context) {
	var banks []models.BankCredentials
    c.Db.Order("id DESC").Find(&banks)

	if len(banks) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "bank not found"})
		return
	}

    ctx.JSON(http.StatusOK, gin.H{"banks": banks})
}

func (c SystemController) GetBank(ctx *gin.Context) {
	id := ctx.Param("id")
    var bank models.BankCredentials
    result := c.Db.Where("id =?", id).First(&bank)

    if result.Error!= nil {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "bank not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"bank": bank})
}

func (c SystemController) UpdatedBank(ctx *gin.Context) {
	id := ctx.Param("id")
    var reqBody request.BankReq
    if err := ctx.Bind(&reqBody); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return;
    }

    if err := c.Validate.Struct(reqBody); err!= nil {
        validationErrors := err.(validator.ValidationErrors);
        errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
        ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
        return;
    }

    var bank models.BankCredentials
    result := c.Db.Where("id =?", id).First(&bank)
	if result.RowsAffected == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "bank not found"})
        return
    }

	bank.BankName = &reqBody.BankName
	bank.AppName = &reqBody.AppName
	bank.API_KEY = &reqBody.ApiKey
	bank.API_SECRET = &reqBody.ApiSecret
	bank.Biller_ID = &reqBody.BillerID
	bank.Merchant_ID = &reqBody.MerchantID
	bank.Terminal_ID = &reqBody.TerminalID
	
	response := c.Db.Save(&bank)

	if response.Error != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"messsage": "Failed to update bank"})
        return
    }

	ctx.JSON(http.StatusOK, gin.H{"messsage": "updated bank successfully"})
}

func (c SystemController) DeleteBank(ctx *gin.Context) {
	id := ctx.Param("id")
    var bank models.BankCredentials
	reCheck := c.Db.Where("id =?", id).First(&bank)
    result := c.Db.Where("id =?", id).Delete(&bank)

    if result.Error != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "bank not found"})
        return
    }

	if reCheck.Error != nil || reCheck.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "bank not found"})
        return
	}

    ctx.JSON(http.StatusOK, gin.H{"message": "bank deleted successfully"})
}