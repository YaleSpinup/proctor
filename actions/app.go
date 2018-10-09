package actions

import (
	"os"

	"github.com/YaleSpinup/proctor/libs/s3"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

// S3 has the initialized client session
var S3 s3.Client

var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_proctor_session",
		})

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// initialize S3 client session
		S3 = s3.NewSession(os.Getenv("S3_API_KEY"), os.Getenv("S3_API_SECRET"), os.Getenv("S3_REGION"), os.Getenv("S3_BUCKET"))

		userAPI := app.Group("/v1/proctor")
		userAPI.GET("/ping", PingPong)
		userAPI.GET("/risklevels", RiskLevelsGet)
		userAPI.GET("/{campaign}/questions", QuestionsGet)
		userAPI.POST("/{campaign}/responses", ResponsesPost)
	}

	return app
}
