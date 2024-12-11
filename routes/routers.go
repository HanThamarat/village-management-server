package routes

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"tectnexify.github.com/e-payment/controllers"
	"tectnexify.github.com/e-payment/middlewares"
)

func Routers(
	authController *controllers.AuthController, 
	systemController *controllers.SystemController, 
	paymentController *controllers.PaymentController, 
	villageController *controllers.VillageController, 
	roleController *controllers.RoleController,
	houseController *controllers.HouseController,
	opnController *controllers.OPNController) *gin.Engine {

	service := gin.Default()

	// can configure cors middlewares is here
	service.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	secret := os.Getenv("JWT_SECRET")

	// doing routers prefix path here
	service.GET("/", indexPath)
	auth := service.Group("/auth")
	system := service.Group("/system")
	payment := service.Group("/payment")
	village := service.Group("/legal")
	opnPayment := service.Group("/payment/opn");

	system.Use(middlewares.JWTMiddleware(secret))
	auth.POST("/signup", authController.Register)
	auth.POST("/signin", authController.SignIn)

	// here u can create a new route service
	auth.Use(middlewares.JWTMiddleware(secret))
	auth.GET("/users", authController.CurrentUser)
	system.POST("/bank", systemController.CreateBank)
	system.GET("/bank", systemController.GetBanks)
	system.GET("/bank/:id", systemController.GetBank)
	system.PUT("/bank/:id", systemController.UpdatedBank)
	system.DELETE("/bank/:id", systemController.DeleteBank)
	payment.POST("/accesstoken", paymentController.GenAccessToken)
	village.Use(middlewares.JWTMiddleware(secret))
	village.POST("/village", villageController.CreateVillage)
	village.GET("/village", villageController.GetVillages)
	village.GET("/village/:id", villageController.GetVIllage)
	village.PUT("/village/:id", villageController.UpdateVillage)
	village.DELETE("/village/:id", villageController.DeleteVillage)
	village.POST("/house", houseController.CreateHouse)
	village.GET("/village_house/:villageId", houseController.Gethouses)
	village.GET("/house/:id", houseController.GetHouse)
	village.PUT("/house/:id", houseController.UpdateHouse)
	village.DELETE("/house/:id", houseController.DeleteHouse)
	service.Use(middlewares.JWTMiddleware(secret))
	service.POST("/role", roleController.CreateRole)
	service.GET("/role", roleController.GetRoles)
	service.GET("/role/:id", roleController.GetRole)
	service.PUT("/role/:id", roleController.UpdateRole)
	service.DELETE("/role/:id", roleController.DeleteRole)
	opnPayment.Use(middlewares.JWTMiddleware(secret))
	opnPayment.POST("/generateToken", opnController.CreateTokenOPN);

	return service
}

func indexPath(ctx *gin.Context) {
	ctx.String(200, "Hello, World!")
}
