package actions

import (
	"os"

	"github.com/YaleSpinup/proctor/libs/s3"
	"github.com/YaleSpinup/proctor/proctor"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	paramlogger "github.com/gobuffalo/mw-paramlogger"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

var (
	app *buffalo.App

	// ENV is used to help switch settings based on where the
	// application is being run. Default is "development".
	ENV = envy.Get("GO_ENV", "development")

	// S3 has the initialized client session
	S3 s3.Client

	// Version is the main version number
	Version = proctor.Version

	// VersionPrerelease is a prerelease marker
	VersionPrerelease = proctor.VersionPrerelease

	// BuildStamp is the timestamp the binary was built, it should be set at buildtime with ldflags
	BuildStamp = proctor.BuildStamp

	// GitHash is the git sha of the built binary, it should be set at buildtime with ldflags
	GitHash = proctor.GitHash
)

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

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// initialize S3 client session
		S3 = s3.NewSession(os.Getenv("S3_API_KEY"), os.Getenv("S3_API_SECRET"), os.Getenv("S3_REGION"), os.Getenv("S3_BUCKET"))

		userAPI := app.Group("/v1/proctor")
		userAPI.GET("/ping", PingPong)
		userAPI.GET("/version", VersionHandler)

		userAPI.GET("/risklevels", RiskLevelsGet)
		userAPI.GET("/{campaign}/questions", QuestionsGet)
		userAPI.POST("/{campaign}/responses", ResponsesPost)
	}

	return app
}
