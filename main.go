package main

import (
	_ "net/http/pprof"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var version = "devel" // -ldflags="-X 'main.version=X.X.X'"

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6068", nil))
	// }()

	// logger
	log := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger().Hook(SeverityHook{})
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// application
	app := cli.NewApp()
	cli.VersionFlag = &cli.BoolFlag{Name: "version", Aliases: []string{"V"}}

	app.Name = "nel-reports-controller"
	app.Version = version
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "MindHunter86",
			Email: "mindhunter86@vkom.cc",
		},
	}
	app.Copyright = "(c) 2022-2023 mindhunter86"
	app.Usage = "Backend for processing browsers NEL reports"

	app.Flags = []cli.Flag{
		// common flags
		&cli.StringFlag{
			Name:    "log-level",
			Aliases: []string{"l"},
			Value:   "debug",
			Usage:   "levels: trace, debug, info, warn, err, panic, disabled",
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.BoolFlag{
			Name:    "quite",
			Aliases: []string{"q"},
			Usage:   "Flag is equivalent to --log-level=quite",
		},
		// http client settings
		&cli.BoolFlag{
			Name:  "http-client-insecure",
			Usage: "Flag for TLS certificate verification disabling",
		},
		&cli.DurationFlag{
			Name:  "http-client-timeout",
			Usage: "Internal HTTP client connection `TIMEOUT` (format: 1000ms, 1s)",
			Value: 3 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-tcp-timeout",
			Usage: "",
			Value: 1 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-tls-handshake-timeout",
			Usage: "",
			Value: 1 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-idle-timeout",
			Usage: "",
			Value: 300 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-keepalive-timeout",
			Usage: "",
			Value: 300 * time.Second,
		},
		&cli.IntFlag{
			Name:  "http-max-idle-conns",
			Usage: "",
			Value: 100,
		},
		&cli.BoolFlag{
			Name:  "http-debug",
			Usage: "",
		},
	}

	app.Action = func(c *cli.Context) (e error) {
		var lvl zerolog.Level
		if lvl, e = zerolog.ParseLevel(c.String("log-level")); e != nil {
			log.Fatal().Err(e)
		}

		zerolog.SetGlobalLevel(lvl)
		if c.Bool("quite") {
			zerolog.SetGlobalLevel(zerolog.Disabled)
		}

		log.Debug().Msg("ready...")
		log.Debug().Strs("args", os.Args).Msg("")

		// return application.NewApp(c, &log).Bootstrap()
		return nil
	}

	// TODO sort.Sort of Flags uses too much allocs; temporary disabled
	// sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if e := app.Run(os.Args); e != nil {
		log.Fatal().Err(e).Msg("")
	}
}

type SeverityHook struct{}

func (SeverityHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if level > zerolog.DebugLevel || version != "devel" {
		return
	}

	rfn := "unknown"
	pcs := make([]uintptr, 1)

	if runtime.Callers(4, pcs) != 0 {
		if fun := runtime.FuncForPC(pcs[0] - 1); fun != nil {
			rfn = fun.Name()
		}
	}

	fn := strings.Split(rfn, "/")
	e.Str("func", fn[len(fn)-1:][0])
}
