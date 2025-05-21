package bootstrap

import (
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	datahandler "github.com/estella-studio/leon-backend/internal/app/data/interface/rest"
	datarepository "github.com/estella-studio/leon-backend/internal/app/data/repository"
	datausecase "github.com/estella-studio/leon-backend/internal/app/data/usecase"
	pinghandler "github.com/estella-studio/leon-backend/internal/app/ping/interface/rest"
	userhandler "github.com/estella-studio/leon-backend/internal/app/user/interface/rest"
	userrepository "github.com/estella-studio/leon-backend/internal/app/user/repository"
	userusecase "github.com/estella-studio/leon-backend/internal/app/user/usecase"
	"github.com/estella-studio/leon-backend/internal/infra/env"
	"github.com/estella-studio/leon-backend/internal/infra/jwt"
	"github.com/estella-studio/leon-backend/internal/infra/mailer"
	"github.com/estella-studio/leon-backend/internal/infra/mysql"
	"github.com/estella-studio/leon-backend/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Start() (*fiber.App, uint, error) {
	config, err := env.New()
	if err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}

	log.Println("database connected")

	err = mysql.Migrate(database)
	if err != nil {
		return nil, 0, err
	} else {
		log.Println("database migration complete")
	}

	val := validator.New()

	app := fiber.New(
		fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		},
	)

	jwt := jwt.NewJWT(config)

	mailer := mailer.NewMailer(config)

	middleware := middleware.NewMiddleware(*jwt)

	app.Use(
		cache.New(),
		idempotency.New(),
		cors.New(
			cors.Config{
				AllowHeaders: "*",
				AllowOrigins: "*",
				AllowMethods: "*",
			}),
		limiter.New(
			limiter.Config{
				Max:               config.LimiterMax,
				Expiration:        time.Duration(config.LimiterExpirationMinute) * 60,
				LimiterMiddleware: limiter.SlidingWindow{},
			}),
		logger.New(
			logger.Config{
				Format: "${ip} - - [${time}] ${method} ${url} ${protocol} ${status} ${bytesSent} ${referer} ${ua}\n",
			},
		),
	)

	v1 := app.Group("/api/v1")

	pinghandler.NewPingHandler(v1)
	userRepository := userrepository.NewUserMySQL(database)
	userUseCase := userusecase.NewUserUseCase(userRepository, jwt)
	userhandler.NewUserHandler(v1, val, middleware, userUseCase, config, mailer)
	dataRepository := datarepository.NewDataMySQL(database)
	dataUseCase := datausecase.NewDataUseCase(dataRepository, jwt)
	datahandler.NewDataHandler(v1, val, middleware, dataUseCase)

	log.Printf("listening on port %d", config.AppPort)

	return app, config.AppPort, nil
}
