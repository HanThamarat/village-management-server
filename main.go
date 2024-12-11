package main

import (
	_"fmt"
	_"log"
	"net/http"
	"time"
	_"os"
	_"github.com/joho/godotenv"

	"github.com/go-playground/validator/v10"
	_"github.com/omise/omise-go"
	_"github.com/omise/omise-go/operations"
	"tectnexify.github.com/e-payment/config"
	"tectnexify.github.com/e-payment/controllers"
	"tectnexify.github.com/e-payment/routes"
)

func main() {
	// err := godotenv.Load()
	// if err!= nil {
    //     log.Fatal("Error loading.env file")
    // }

	// var(
	// 	OmisePublicKey = os.Getenv("OmisePublicKey")
	// 	OmiseSecretKey = os.Getenv("OmiseSecretKey")
	// )

	// client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	// if e != nil {
	// 	log.Fatal(e);
	// }

	// token := "tokn_test_6207qyx6y0jzabqd40m";

	// charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
	// 	Amount: 100000,
	// 	Currency: "thb",
	// 	Card: token,
	// }
	// if e := client.Do(charge, createCharge); e != nil {
	// 	log.Fatal(e)
	// }

	// fmt.Printf("charge: %s amount: %s %d\n", charge.ID, charge.Currency, charge.Amount)

	db := config.DatabaseConnection();
	client := config.InitOmise();

	validate := validator.New();

	//send prams db and validate lib to controller 
	authController := controllers.NewAuthControllerImpl(db, validate);
	SystemController := controllers.NewSystemController(db, validate);
	paymentController := controllers.NewPaymentControllerImpl(db, validate);
	villageController := controllers.NewVillageControllerImpl(db, validate);
	roleController := controllers.NewRoleControllerImpl(db, validate);
	houseController := controllers.NewHouseControllerImpl(db, validate);
	opnController := controllers.NewOPNControllerImpl(db, validate, client);

	routes := routes.Routers(
		authController, 
		SystemController, 
		paymentController, 
		villageController, 
		roleController, 
		houseController,
		opnController);

	server := &http.Server{
		Addr:           ":8080",
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe();
	if err != nil {
		panic(err)
	}
}