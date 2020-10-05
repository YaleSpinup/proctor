module github.com/YaleSpinup/proctor

go 1.12

require (
	github.com/aws/aws-sdk-go v1.23.2
	github.com/gobuffalo/buffalo v0.15.5
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/mw-paramlogger v0.0.0-20190224201358-0d45762ab655
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/gobuffalo/suite/v3 v3.0.0
	github.com/gobuffalo/x v0.1.0
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/rs/cors v1.7.0
)

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
