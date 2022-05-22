package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/saying-yan/embedded_system_course_project_backend/internal"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type Options struct {
	flags.Options
	EnableProfile bool   `short:"p" long:"enable-pprof"  description:"enable pprof"`
	ConfigFile    string `short:"c" long:"config" description:"server config file" default:"./config/server.yaml"`
}

func main() {
	// parse command line args
	var options Options
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			fmt.Printf("parse flag error: %s\n", err.Error())
			os.Exit(1)
		default:
			fmt.Printf("parse flag error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if options.EnableProfile {
		go func() {
			fmt.Println(http.ListenAndServe(":8088", nil))
		}()
	}

	s, err := internal.NewServer(options.ConfigFile)
	if err != nil {
		fmt.Printf("server start error: %s\n", err.Error())
		os.Exit(1)
	}
	s.Serve()
}
