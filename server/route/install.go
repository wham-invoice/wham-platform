package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/util"
)

// Installer configures a gin.RouterGroup.
type Installer interface {
	Install(*gin.RouterGroup)
}

// Group is an Installer that's convenient to specify.
type Group struct {
	Path       string
	Prereqs    []gin.HandlerFunc
	Installers []Installer
}

// Install is part of the Installer interface.
func (g Group) Install(rg *gin.RouterGroup) {
	gr := rg.Group(g.Path)
	gr.Use(g.Prereqs...)
	for _, inst := range g.Installers {
		inst.Install(gr)
	}
}

// Endpoint is an Installer that's convenient to specify.
type Endpoint struct {
	Method  string
	Path    string
	Prereqs []gin.HandlerFunc
	Do      func(*gin.Context) (interface{}, error)
}

// Install is part of the Installer interface.
func (ep Endpoint) Install(rg *gin.RouterGroup) {
	handlers := append(ep.Prereqs, func(c *gin.Context) {
		resp, err := ep.Do(c)
		if err != nil {
			util.Logger.Errorf("request failed: %s", errors.ErrorStack(err))
			Abort(c, err)
			return
		}
		if resp == nil {
			c.AbortWithStatus(http.StatusNoContent)
			return
		} else {
			c.JSON(http.StatusOK, resp)
		}
	})
	rg.Handle(ep.Method, ep.Path, handlers...)

	// If we've already got OPTIONS handled, don't worry; assume it was another
	// Endpoint with a different method that set up the handler.
	defer func() { recover() }()
	rg.Handle("OPTIONS", ep.Path, func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
}

// Prereqs is useful for setting Prereqs without explicit types.
func Prereqs(p ...gin.HandlerFunc) []gin.HandlerFunc {
	return p
}

// Installers is userful for setting Installers without explicit types.
func Installers(i ...Installer) []Installer {
	return i
}
