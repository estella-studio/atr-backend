package bootstrap

import (
	"fmt"
	"log"

	userhandler "github.com/estella-studio/leon-backend/internal/app/user/interface/rest"
	userrepository "github.com/estella-studio/leon-backend/internal/app/user/repository"
	userusecase "github.com/estella-studio/leon-backend/internal/app/user/usecase"
	"github.com/estella-studio/leon-backend/internal/infra/env"
	"github.com/estella-studio/leon-backend/internal/infra/jwt"
	"github.com/estella-studio/leon-backend/internal/infra/mysql"
	"github.com/estella-studio/leon-backend/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		return err
	}

	log.Println("loaded config")

	database, err := mysql.New(fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	))
	if err != nil {
		return err
	}

	log.Println("database connected")

	err = mysql.Migrate(database)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("database migration complete")
	}

	val := validator.New()

	app := fiber.New()

	jwt := jwt.NewJWT(config)

	middleware := middleware.NewMiddleware(*jwt)

	app.Use(
		cors.New(cors.Config{
			AllowHeaders: "*",
			AllowOrigins: "*",
			AllowMethods: "*",
		}),
	)

	v1 := app.Group("/api/v1")

	userRepository := userrepository.NewUserMySQL(database)
	userUseCase := userusecase.NewUserUseCase(userRepository, jwt)
	userhandler.NewUserHandler(v1, val, middleware, userUseCase)

	log.Printf("listening on port %d", config.AppPort)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
