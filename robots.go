package buffalo

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/sirupsen/logrus"
)

func NewRobotsHandler(r *render.Engine) Handler {
	return func(c Context) error {
		contents, err := r.AssetsBox.MustString("robots.txt")
		if err != nil {
			logrus.Error(err)
			return c.Render(200, r.String("User-agent: *\nDisallow:"))
		}

		return c.Render(200, r.String(contents))
	}
}
