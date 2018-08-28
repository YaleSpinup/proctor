package grifts

import (
	"github.com/YaleSpinup/proctor/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
