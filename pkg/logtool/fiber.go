package logtool

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Custom Fiber logger middleware for zap
func FiberZapLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := c.Context().Time()
		err := c.Next()
		stop := c.Context().Time()
		status := c.Response().StatusCode()
		addBody := false
		title := http.StatusText(status)

		var logFn func(string, ...interface{})
		switch {
		case status >= 500:
			logFn = getLogger().Errorw
			addBody = true
		case status >= 400:
			logFn = getLogger().Warnw
			addBody = true
		default:
			logFn = getLogger().Infow
		}

		fields := []interface{}{
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"caller", "logtool.FiberZapLogger",
		}
		if !jsonFormat {
			fields = append(fields,
				"latency", stop.Sub(start).String(),
				"ip", c.IP(),
				"user_agent", c.Get("User-Agent"),
			)
		}

		if err != nil {
			fields = append(fields, "error", err.Error())
		}
		if addBody {
			fields = append(fields, "body", string(c.Response().Body()))
		}

		logFn(title, fields...)
		return err
	}
}
