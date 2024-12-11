package controllers

import (
	_ "bytes"
	_ "encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/google/uuid"
	"gorm.io/gorm"
	"tectnexify.github.com/e-payment/data/request"
	"tectnexify.github.com/e-payment/models"
)

type HouseController struct {
	Db *gorm.DB
	Validate *validator.Validate
}

func NewHouseControllerImpl(Db *gorm.DB, validate *validator.Validate) *HouseController {
	return &HouseController{
		Db: Db,
        Validate: validate,
	}
}

func (c HouseController) CreateHouse(ctx *gin.Context) {
	var reqBody request.HouseReq
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return;
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors);
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
		ctx.JSON(http.StatusBadRequest, gin.H{"err": errorMessage});
		return;
	}

	userId, exists := ctx.Get("userID");
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized to create house"});
        return;
	}

	userIdInt, ok := userId.(uint);
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "invalid type params"});
		return;
	}

	var existingHouse models.VallageHouses
	result := c.Db.Where("house_name = ?", reqBody.HouseName).First(&existingHouse);
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "this house name already exist" });
		return;
	}

	newHouse := models.VallageHouses{
		HouseName: &reqBody.HouseName,
		HousePrice: &reqBody.HousePrice,
		WaterBill: reqBody.WaterBill,
		ElectricityBill: reqBody.ElectricityBill,
		OtherBill: reqBody.OtherBill,
		Active: &reqBody.Active,
		VallageOwnerShipID: reqBody.VallageOwnerShipID,
        CreateBy: &userIdInt,
        UpdateBy: &userIdInt,
	}

	userRes := c.Db.Create(&newHouse);
	if userRes.Error!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create house"});
        return;
    }

	lastInsertId := newHouse.ID;
	
	var HoouseResponse models.VallageHouses
	response := c.Db.Where("id = ?", lastInsertId).
	Preload("VallageOwnerShip").
	First(&HoouseResponse);

	if response.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "house not found"});
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "house created successfully",
		"body": HoouseResponse,
	});
}

func (c *HouseController) Gethouses(ctx *gin.Context) {
	villageId := ctx.Param("villageId");
	var houses []models.VallageHouses
	c.Db.Where("vallage_owner_ship_id = ?", villageId).
	Order("id DESC").
	Find(&houses)

	if len(houses) == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "house not found"})
        return
    }

	ctx.JSON(http.StatusOK, gin.H{
		"message": "getting all houses successfully",
		"body": houses,
	})
}

func (c *HouseController) GetHouse(ctx *gin.Context) {
	id := ctx.Param("id");
    var house models.VallageHouses
    c.Db.Where("id =?", id).
    Preload("VallageOwnerShip").
    First(&house)

    if house.ID == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "house not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "getting house successfully",
        "body": house,
    })
}

func (c *HouseController) UpdateHouse(ctx *gin.Context) {
	id := ctx.Param("id");
	var reqBody request.HouseReq
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return;
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors);
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
		ctx.JSON(http.StatusBadRequest, gin.H{"err": errorMessage});
		return;
	}

	userId, exists := ctx.Get("userID");
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized to create house"});
        return;
	}

	userIdInt, ok := userId.(uint);
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "invalid type params"});
		return;
	}

	var existingHouse models.VallageHouses
	result := c.Db.Where("id = ?", id).First(&existingHouse);
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "this house not found" });
		return;
	}

	existingHouse.HouseName = &reqBody.HouseName
	existingHouse.HousePrice = &reqBody.HousePrice
	existingHouse.WaterBill = reqBody.WaterBill
	existingHouse.ElectricityBill = reqBody.ElectricityBill
	existingHouse.OtherBill = reqBody.OtherBill
	existingHouse.Active = &reqBody.Active
	existingHouse.UpdateBy = &userIdInt

    userRes := c.Db.Save(&existingHouse);

	if userRes.Error!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create house"});
        return;
    }
	
	var HoouseResponse models.VallageHouses
	response := c.Db.Where("id = ?", id).
	Preload("VallageOwnerShip").
	First(&HoouseResponse);

	if response.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "house not found"});
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "house created successfully",
		"body": HoouseResponse,
	});
}

func (c *HouseController) DeleteHouse(ctx *gin.Context) {
	id := ctx.Param("id");
    var house models.VallageHouses
    result := c.Db.Where("id =?", id).Delete(&house)

    if result.Error != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "house not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "house deleted successfully"})
}