package controllers

import (
	"fmt"
	"net/http"
	_"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"tectnexify.github.com/e-payment/data/request"
	"tectnexify.github.com/e-payment/models"
)

type VillageController struct {
	Db       *gorm.DB
	Validate *validator.Validate
}

func NewVillageControllerImpl(Db *gorm.DB, validate *validator.Validate) *VillageController {
	return &VillageController{Db: Db, Validate: validate}
}

func (c *VillageController) CreateVillage(ctx *gin.Context) {
	var reqBody request.VillageOwerReq
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

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "no userID found"})
		return
	}

	strUserID, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Valind UserID type"})
		return
	}

	var existingUser models.Users
	if err := c.Db.First(&existingUser, userID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Legal entity does not exist"})
		return
	}

	var existingVillageOwer models.VallageOwnerShips
	result := c.Db.Where("village_name = ?", reqBody.VillageName).First(&existingVillageOwer)
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Village with the same name and legal entity already exists"})
		return
	}

	NewVillage := models.VallageOwnerShips{
		VillageName:   &reqBody.VillageName,
		LegalEntityID: strUserID,
		CreateByID:    strUserID,
		UpdateByID:    strUserID,
	}

	response := c.Db.Create(&NewVillage)
	if response.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create village", "error": response.Error.Error()})
		return
	}

	lastInsertId := NewVillage.ID;

	var vill models.VallageOwnerShips
	responses := c.Db.Where("id = ?", lastInsertId).
		Preload("LegalEntity").
		Preload("Create").
		Preload("Update").
		First(&vill)

	if responses.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "village not found"})
		return
	}
	

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create village successfully",
		"body": vill,
	})
}

func (c *VillageController) GetVillages(ctx *gin.Context) {

	value, exists := ctx.Get("userID")
	if !exists {
		fmt.Print("not has userID")
	}
	fmt.Println(value)
	var villages []models.VallageOwnerShips
	c.Db.Where("legal_entity_id = ?", value).
		Order("id DESC").
		Preload("LegalEntity").
		Preload("Create").
		Preload("Update").
		Find(&villages)

	if len(villages) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "villages not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"villages": villages})
}

func (c *VillageController) GetVIllage(ctx *gin.Context) {
	id := ctx.Param("id")
	var village models.VallageOwnerShips
	result := c.Db.Where("id =?", id).
		Preload("LegalEntity").
		Preload("Create").
		Preload("Update").
		First(&village)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "village not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"village": village})
}

func (c *VillageController) UpdateVillage(ctx *gin.Context) {
	id := ctx.Param("id")
	var reqBody request.VillageOwerReq
	if err := ctx.Bind(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	var village models.VallageOwnerShips
	result := c.Db.Where("id =?", id).First(&village)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "village not found"})
		return
	}

	village.VillageName = &reqBody.VillageName
	// village.UpdateByID = reqBody.LegalEntity

	if response := c.Db.Save(&village); response.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update village"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Updated village successfully"})
}

func (c *VillageController) DeleteVillage(ctx *gin.Context) {
	id := ctx.Param("id")
	var village models.VallageOwnerShips
	reCkeck := c.Db.Where("id =?", id).First(&village)
	result := c.Db.Where("id =?", id).Delete(&village)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "village not found"})
		return
	}

	if reCkeck.Error != nil || reCkeck.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "village not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "village delated successfully"})
}
