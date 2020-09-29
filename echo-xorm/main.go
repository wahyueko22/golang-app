package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"

	"github.com/corvinusz/echo-xorm/app"
	"github.com/corvinusz/echo-xorm/app/ctx"
)

var (
	configFlag = flag.String("config",
		"/usr/local/etc/echo-xorm-config.toml",
		"-config=\"path-to-your-config-file\" ")
)

func main() {
	// parse flags
	flag.Parse()

	var (
		err error
		a   *app.Application
	)

	flags := &ctx.Flags{
		CfgFileName: *configFlag,
	}

	// create application
	a, err = app.New(flags)
	if err != nil {
		log.Fatal("error ", os.Args[0]+" initialization error: "+err.Error())
		os.Exit(1)
	}

	// log initialization
	a.Ctx.Logger.Info("appcontrol", "application initialized successfully")
	a.Ctx.Logger.Info("appcontrol", "CONFIG: "+fmt.Sprintf("%+v", a.Ctx.Config))
	a.Ctx.Logger.Info("appcontrol", "FLAGS: "+fmt.Sprintf("%+v", flags))
	a.Ctx.Logger.Info("appcontrol", "JWT Signing Key: "+
		base64.StdEncoding.EncodeToString(a.Ctx.JWTSignKey))

	go func() {
		// here we go
		a.Run()
	}()

	// signal control
	sigstop := make(chan os.Signal, 1)
	signal.Notify(sigstop, syscall.SIGTERM, os.Interrupt)

	sig := <-sigstop

	if a.Ctx.Logger != nil {
		a.Ctx.Logger.Info("appcontrol", os.Args[0]+" caught signal "+sig.String())
	}

	// shutdown server on signal
	err = a.Shutdown()
	if err != nil {
		if a.Ctx.Logger != nil {
			a.Ctx.Logger.Error("error stopping server", err)
		}
		os.Exit(1)
	}

}
