package app

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	gLog *zerolog.Logger
	gCtx *cli.Context
)

func init() {

}

type App struct {
	fb  *fiber.App
	nel *NelControllers
}

func NewApp() *App {
	return &App{}
}

func (m *App) Bootstrap() error {
	m.fb = fiber.New(fiber.Config{
		EnableTrustedProxyCheck: len(gCtx.String("http-trusted-proxies")) > 0,
		TrustedProxies:          strings.Split(gCtx.String("http-trusted-proxies"), ","),
		ProxyHeader:             fiber.HeaderXForwardedFor,

		AppName:      gCtx.App.Name,
		ServerHeader: gCtx.App.Name,

		StrictRouting:             true,
		DisableDefaultContentType: true,
		DisableDefaultDate:        true,

		Prefork:      gCtx.Bool("http-prefork"),
		IdleTimeout:  300 * time.Second,
		ReadTimeout:  1000 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,

		RequestMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodOptions,
		},
	})

	m.fiberConfigure()

	log.Fatal(m.fb.Listen("127.0.0.1:8088"))

	return nil
}

func (m *App) fiberConfigure() {

	// CORS serving
	if gCtx.Bool("http-cors") {
		m.fb.Use(cors.New(cors.Config{
			AllowHeaders: strings.Join([]string{
				fiber.HeaderContentType,
			}, ","),
			AllowOrigins: "*",
			AllowMethods: strings.Join([]string{
				fiber.MethodPost,
			}, ","),
		}))
	}

	// Routes
	m.fb.Get("/", m.nel.root)
	m.fb.Post("/nel-report", m.nel.nel_report)

	//
}

type NelControllers struct{}

func (m *NelControllers) root(c *fiber.Ctx) error {
	return c.SendString("Not found :(")
}

func (m *NelControllers) nel_report(c *fiber.Ctx) error {
	return c.SendString("1234")
}
