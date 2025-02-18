package main

import (
	"log"

	"github.com/Tech-With-Tim/cdn/api"
	"github.com/Tech-With-Tim/cdn/server"
	"github.com/go-chi/chi/v5"

	// _ "net/http/pprof"  // only in use when profiling
	"os"

	"github.com/Tech-With-Tim/cdn/utils"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

func main() {
	//Export Env Variables If exist
	//err := utils.ExportVariables()
	config, err := utils.LoadConfig("./", "app")
	if err != nil {
		log.Fatalln(err.Error())
	}
	//Register Commands
	commands(config)

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadconfig(testconf bool) (conf utils.Config, err error) {
	if testconf {
		conf, err = utils.LoadConfig("./", "test")
	} else {
		conf, err = utils.LoadConfig("./", "app")
	}
	return
}

func commands(config utils.Config) {
	app.Commands = []*cli.Command{
		{
			Name:  "migrate_up",
			Usage: "Migrate DB to latest version",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "test",
					Aliases: []string{"t"},
					Usage:   "loads test.env instead of app.env",
				},
			},
			Action: func(c *cli.Context) error {
				conf, err := loadconfig(c.Bool("test"))
				if err != nil {
					return err
				}
				err = utils.MigrateUp(conf, "./models/migrations/")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "dropdb",
			Usage: "Drop the DB",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "test",
					Aliases: []string{"t"},
					Usage:   "loads test.env instead of app.env",
				},
			},
			Action: func(c *cli.Context) error {
				conf, err := loadconfig(c.Bool("test"))
				if err != nil {
					return err
				}
				err = utils.MigrateDown(conf, "./models/migrations/")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "migrate_steps",
			Usage: "Migrate with Steps",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "steps",
					Usage: "Number of steps of migrations to run",
				},
				&cli.BoolFlag{
					Name:    "test",
					Aliases: []string{"t"},
					Usage:   "loads test.env instead of app.env",
				},
			},
			Action: func(c *cli.Context) error {
				conf, err := loadconfig(c.Bool("test"))
				if err != nil {
					return err
				}
				err = utils.MigrateSteps(c.Int("steps"), conf, "./models/migrations/")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "runserver",
			Usage: "Run Api Server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "host",
					Usage:   "Host on which server has to be run",
					Value:   "localhost",
					Aliases: []string{"H"},
				},
				&cli.IntFlag{
					Name:    "port",
					Usage:   "Port on which server has to be run",
					Value:   5000,
					Aliases: []string{"P"},
				},
			},
			Action: func(c *cli.Context) error {
				s := server.NewServer(config)
				//Create Routers Here
				CdnRouter := chi.NewRouter()

				//Add Routes to Routers Here
				api.MainRouter(CdnRouter, s.Store, config)
				//Mount Routers here
				s.Router.Mount("/", CdnRouter)
				// r.Mount("/debug/", middleware.Profiler()) // Only in use when profiling
				//Store Router in Struct
				err := s.RunServer(c.String("host"), c.Int("port"))
				return err
			},
		},
	}
}
