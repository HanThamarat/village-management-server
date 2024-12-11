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

type RoleController struct {
	Db       *gorm.DB
	Validate *validator.Validate
}

func NewRoleControllerImpl(Db *gorm.DB, vaidate *validator.Validate) *RoleController {
	return &RoleController{
		Db: Db,
		Validate: vaidate,
	}
}

func (c *RoleController) CreateRole(ctx *gin.Context) {
	var reqBody request.RoleReq
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return;
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors);
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
		ctx.JSON(http.StatusBadRequest, gin.H{"err": errorMessage})
		return;
	}

	var existingRole models.Roles
	result := c.Db.Where("name_th = ? or name_en = ?", reqBody.Name_TH, reqBody.Name_EN).First(&existingRole);
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "this role already exists"})
		return;
	}

	newRole := models.Roles{
		Name_TH: &reqBody.Name_TH,
		Name_EN: &reqBody.Name_EN,
		Active: &reqBody.Active,
	}

	response := c.Db.Create(&newRole);

	if response.Error!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create role"})
        return;
    }

	ctx.JSON(http.StatusOK, gin.H{"message": "created role successfully"})
}

func (c *RoleController) GetRoles(ctx *gin.Context) {
	var roles []models.Roles
    c.Db.Order("id DESC").Find(&roles)

    if len(roles) == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "role not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (c *RoleController) GetRole(ctx *gin.Context) {
	id := ctx.Param("id");
	var role models.Roles
	result := c.Db.Where("id = ?", id).First(&role);

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "role not found"})
		return;
	}

	// value, exists := ctx.Get("userID")
	// if !exists {
	// 	fmt.Print("not has userID");
	// }

	// fmt.Print(value);

	ctx.JSON(http.StatusOK, gin.H{"role": role})
}

func (c *RoleController) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id");
    var reqBody request.RoleReq
    if err := ctx.BindJSON(&reqBody); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
        return;
    }

    if err := c.Validate.Struct(reqBody); err!= nil {
        validationErrors := err.(validator.ValidationErrors);
        errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
        ctx.JSON(http.StatusBadRequest, gin.H{"err": errorMessage})
        return;
    }

    var role models.Roles
    result := c.Db.Where("id =?", id).First(&role);
	if result.Error!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": "role not found"})
        return;
    }

    role.Name_TH = &reqBody.Name_TH
	role.Name_EN = &reqBody.Name_EN
	role.Active = &reqBody.Active

	response := c.Db.Save(&role);

	if response.Error!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role"})
        return;
    }

	ctx.JSON(http.StatusOK, gin.H{"message": "updated role successfully"})
}

func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id");
    var role models.Roles
    result := c.Db.Where("id =?", id).First(&role);

    if result.Error != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": "role not found"})
        return;
    }

    response := c.Db.Delete(&role)

    if response.Error!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete role"})
        return;
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "deleted role successfully"})
}