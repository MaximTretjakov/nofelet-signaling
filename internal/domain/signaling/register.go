package signaling

import (
	"nofelet/internal/dependency"
	sc "nofelet/internal/domain/signaling/controller"
)

func Register(deps *dependency.Container) {
	c := sc.NewController(deps.Logger, deps.Cfg)

	r := deps.Signaling.Routes
	r.GET("/connect/:uuid", c.GetConnection)
	r.GET("/turn-credentials/generate", c.GetCoTURNCredentials)
}
