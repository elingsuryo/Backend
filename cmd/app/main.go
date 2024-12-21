package main

import (
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/config"
	_ "github.com/TrinityKnights/Backend/docs"
)

// @title Trinity Knights API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase(viper, log)
	redis := config.NewRedisClient(viper, log)
	jwt := config.NewJWT(viper)
	validate := config.NewValidator()
	xendit := config.NewXendit(viper)
	app, log := config.NewEcho()
	gomail := config.NewGomail(viper, log)
	err := config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Cache:    redis,
		App:      app,
		Log:      log,
		Validate: validate,
		JWT:      jwt,
		Viper:    viper,
		Xendit:   xendit,
		Gomail:   gomail,
	})
	if err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	port := viper.GetString("APP_PORT")
	go func() {
		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatal("shutting down the server")
		}
	}()

	config.GracefulShutdown(app, log, 10*time.Second)
}
