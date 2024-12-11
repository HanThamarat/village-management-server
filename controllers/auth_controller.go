package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
	"tectnexify.github.com/e-payment/data/request"
	"tectnexify.github.com/e-payment/hooks"
	"tectnexify.github.com/e-payment/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type AuthController struct {
	Db *gorm.DB
	Validate *validator.Validate
}

func NewAuthControllerImpl(Db *gorm.DB, validate *validator.Validate) *AuthController {
	return &AuthController{Db: Db, Validate: validate}
}

func (c AuthController) Register(ctx *gin.Context) {
	var reqBody request.RegisterReq
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors);
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field());
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return;
	}

	var existingUser models.Users;
	result := c.Db.Where("username = ?", reqBody.Username).First(&existingUser);
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
        return;
	}

	password, err := hooks.EncryptPassword(reqBody.Password);

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
  		return;
	}

	newUser := models.Users{
		Name: reqBody.Name,
		Username: reqBody.Username,
		Password: password,
		Roles: reqBody.RoleReq,
	}

	response := c.Db.Create(&newUser)

	if response.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user register successfully"})
}

func (c AuthController) SignIn(ctx *gin.Context) {
	var reqBody request.SigninReq
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validateErr := err.(validator.ValidationErrors)
		errMessage := fmt.Sprintf("Validation failed for field: %s", validateErr[0].Field());
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : errMessage});
		return;
	}

	var existingUser models.Users
	result := c.Db.Where("username = ?", reqBody.Username).Preload("Role").First(&existingUser);
	if result.RowsAffected < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username not found"})
		return;
	}

	valid := hooks.ComparePassword(reqBody.Password, existingUser.Password)

	if valid == true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password invalid"})
		return;
	}

	newToken := jwt.New(jwt.SigningMethodHS256);

	claims := newToken.Claims.(jwt.MapClaims);
	claims["userId"] = existingUser.ID;
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix();

	token, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")));

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
        return;
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "signin successfully!",
		"body": existingUser,
		"authToken": token,
		});
}

func (c AuthController) CurrentUser(ctx *gin.Context) {
	userid := ctx.MustGet("userID");

	var user models.Users
	result := c.Db.Where("id = ?", userid).Preload("Role").First(&user);

	if result.RowsAffected < 1 {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return;
    }

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get current user successfully",
		"body": user,
	})
}
