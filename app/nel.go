package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/nel-collector/pkg/collector"
	_ "github.com/google/nel-collector/pkg/core"
)

func init() {

}

type NelController struct {
	//
}

func (m *NelController) processReports(c *fiber.Ctx) error {

	contentType := c.Get("Content-Type", "text/html")
	if contentType != "application/reports+json" {
		code, msg := fiber.ErrBadRequest.Code,
			"Must use application/reports+json to upload reports"

		gLog.Warn().Str("sip", c.IP()).Int("status", code).Msg(msg)
		return fiber.NewError(code, msg)
	}

	//

	return nil
}

func (m *NelController) newReportBatch(c *fiber.Ctx, clock collector.Clock) (e error) {
	// var rrl *url.URL
	// if rrl, e = url.Parse(c.Request().URI().String()); e != nil {
	// 	return e
	// }

	// var rbatch collector.ReportBatch
	// rbatch.Time = clock.Now()
	// rbatch.CollectorURL = *rrl
	// rbatch.ClientIP = c.IP()
	// rbatch.ClientUserAgent = c.Get("User-Agent", "undefined")
	// rbatch.Header

	return
}
