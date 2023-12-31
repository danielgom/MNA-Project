// Package core in charge of initialising core configuration, DBs, repositories, services and handler.
package core

import (
	"errors"
	"fmt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	"MNA-project/pkg/config"
	"MNA-project/pkg/config/logger"
	"MNA-project/pkg/context"

	petdb "MNA-project/pkg/internal/pet/db"
	petservice "MNA-project/pkg/internal/pet/service"

	userdb "MNA-project/pkg/internal/user/db"
	userservice "MNA-project/pkg/internal/user/service"

	vacdb "MNA-project/pkg/internal/vaccine/db"
	vaccineservice "MNA-project/pkg/internal/vaccine/service"

	dewormingdb "MNA-project/pkg/internal/deworming/db"
	dewormingservice "MNA-project/pkg/internal/deworming/service"

	vetvisitsdb "MNA-project/pkg/internal/vet_visits/db"
	vetvisitsservice "MNA-project/pkg/internal/vet_visits/service"

	surgerydb "MNA-project/pkg/internal/surgeries/db"
	surgeryservice "MNA-project/pkg/internal/surgeries/service"

	dehandler "MNA-project/pkg/internal/deworming/routes"
	pethandler "MNA-project/pkg/internal/pet/routes"
	surhandler "MNA-project/pkg/internal/surgeries/routes"
	userhandler "MNA-project/pkg/internal/user/routes"
	vachandler "MNA-project/pkg/internal/vaccine/routes"
	vethandler "MNA-project/pkg/internal/vet_visits/routes"
)

var (
	errInvalidToken    = errors.New("invalid token")
	errParse           = errors.New("unable to parse token")
	errMalformedClaims = fmt.Errorf("malformed claims")
	errGetClaim        = fmt.Errorf("unable to get claim")
	errInvalidSigning  = fmt.Errorf("unexpected jwt signing method")
)

// Router initialises api and returns router to serve.
func Router() *echo.Echo {
	router := initialiseAPI()
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	router.Use(doJWTFilter())
	router.Use(doLoggerFilter())

	return router
}

func doJWTFilter() echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		switch c.Request().URL.RequestURI() {
		case "/v1/users/login", "/v1/users/signup":
			return true
		default:
			return false
		}
	}

	return echojwt.WithConfig(echojwt.Config{
		Skipper:                skipper,
		ContinueOnIgnoredError: false,
		ContextKey:             "sub",
		SigningKey:             []byte(config.LoadConfig().JWT.Key),
		ParseTokenFunc:         getParseTokenFunc(),
		NewClaimsFunc:          nil,
	})
}

func doLoggerFilter() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogURI:     true,
		LogURIPath: true,
		LogStatus:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request", zap.String("URI", v.URI),
				zap.Any("latency", v.Latency),
				zap.String("method", v.Method),
				zap.Int("status", v.Status))
			return nil
		},
	})
}

func getParseTokenFunc() func(c echo.Context, auth string) (interface{}, error) {
	c := config.LoadConfig()
	signingKey := []byte(c.JWT.Key)

	return func(c echo.Context, auth string) (any, error) {
		keyFunc := func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != "HS512" {
				return nil, errInvalidSigning
			}
			return signingKey, nil
		}

		token, err := jwt.Parse(auth, keyFunc)
		if err != nil {
			return nil, errParse
		}

		if !token.Valid {
			return nil, errInvalidToken
		}

		// check claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			err = errMalformedClaims
		}

		// func to get a claim by name
		getClaim := func(claim string) (string, error) {
			var str string
			if val, ok := claims[claim]; ok {
				if str, ok = val.(string); !ok {
					return str, errMalformedClaims
				}
			}
			if str == "" {
				err = errMalformedClaims
			}
			return str, err
		}

		userID, err := getClaim("sub")
		if err != nil {
			return nil, errGetClaim
		}

		c.Set("user", userID)

		return token, nil
	}
}

func initialiseAPI() *echo.Echo {
	c := config.LoadConfig()
	DBc := config.InitDatabase(c)
	router := echo.New()

	v := config.GetValidator()
	err := config.AddValidators(v.Validator)
	if err != nil {
		log.Fatalln("could not add validators")
	}

	router.Validator = v

	userRepository := userdb.NewUserRepository(DBc)
	petRepository := petdb.NewPetRepository(DBc)
	vaccineRepository := vacdb.NewVaccineRepository(DBc)
	dewormingRepository := dewormingdb.NewDewormingRepository(DBc)
	vetVisitRepository := vetvisitsdb.NewVetVisitRepository(DBc)
	surgeryRepository := surgerydb.NewSurgeryRepository(DBc)

	userService := userservice.NewUserService(userRepository)
	petService := petservice.NewPetService(petRepository, userService)
	vacService := vaccineservice.NewVaccineService(vaccineRepository, petService)
	dewormingService := dewormingservice.NewDewormingService(dewormingRepository, petService)
	vetVisitService := vetvisitsservice.NewVetVisitService(vetVisitRepository, petService)
	surService := surgeryservice.NewSurgeryService(surgeryRepository, petService)

	petHandler := pethandler.NewHandler(petService)
	userHandler := userhandler.NewHandler(userService)
	vacHandler := vachandler.NewHandler(vacService)
	deHandler := dehandler.NewHandler(dewormingService)

	vetVisitHandler := vethandler.NewHandler(vetVisitService)
	surgeryHandler := surhandler.NewHandler(surService)

	userHandler.Register(router, context.Handler)
	petHandler.Register(router, context.Handler)
	vacHandler.Register(router, context.Handler)
	deHandler.Register(router, context.Handler)

	vetVisitHandler.Register(router, context.Handler)
	surgeryHandler.Register(router, context.Handler)

	return router
}
