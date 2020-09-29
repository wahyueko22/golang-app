package app

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/app/server"
	"github.com/corvinusz/echo-xorm/app/server/users"
	"github.com/corvinusz/echo-xorm/pkg/errors"
	"github.com/corvinusz/echo-xorm/pkg/logger"

	"github.com/BurntSushi/toml"
	"github.com/go-xorm/xorm"
)

// Application define a mode of running app
type Application struct {
	Ctx *ctx.Context
	srv *server.Server
}

// New constructor
func New(flags *ctx.Flags) (*Application, error) {
	a := new(Application)
	a.Ctx = new(ctx.Context)
	// read config file
	err := a.initConfigFromFile(flags.CfgFileName)
	if err != nil {
		return nil, err
	}

	// init Logger
	a.initLogger()

	// set JWTSignKey
	a.Ctx.Logger.Info("appcontrol", "started generation JWT signing key")
	err = a.setJWTSigningKey()
	if err != nil {
		return nil, err
	}
	a.Ctx.Logger.Info("appcontrol", "JWT signing key generated successfully")

	// connect to Db
	a.Ctx.Logger.Info("appcontrol", "started connection to database")
	err = a.initOrm()
	if err != nil {
		return nil, err
	}
	a.Ctx.Logger.Info("appcontrol", "connected to database successfully")

	a.Ctx.Config.Version = "0.1.0-dev"

	return a, nil
}

// Run starts application
func (a *Application) Run() {
	a.srv = server.New(a.Ctx)
	a.srv.Run()
}

// Shutdown gracefully stops server
func (a Application) Shutdown() error {
	// stop server
	a.Ctx.Logger.Info("appcontrol", "stopping server")
	if a.srv != nil {
		err := a.srv.Shutdown()
		if err != nil {
			return err
		}
	}

	// close database connection
	if a.Ctx.Orm != nil {
		a.Ctx.Logger.Info("appcontrol", "closing db connection")
		err := a.Ctx.Orm.Close()
		if err != nil {
			return err
		}
	}

	// close logger
	if a.Ctx.Logger != nil {
		a.Ctx.Logger.Info("appcontrol", "logger stopped")
		a.Ctx.Logger.Info("appcontrol", "quitting")
		a.Ctx.Logger.Close()
	}
	return nil
}

//-----------------------------------------------------------------------------

// readConfig reads configuration file into application Config structure and inits in-memory token storage
func (a *Application) initConfigFromFile(cfgFileName string) error {
	// read config
	tomlData, err := ioutil.ReadFile(cfgFileName) //nolint
	if err != nil {
		return errors.New("Configuration file read error: " + cfgFileName + "\nError:" + err.Error())
	}
	_, err = toml.Decode(string(tomlData), &a.Ctx.Config)
	if err != nil {
		return errors.New("Configuration file decoding error: " + cfgFileName + "\nError:" + err.Error())
	}
	// init Logging data
	if a.Ctx.Config.Logging.ID == "" {
		a.Ctx.Config.Logging.ID = strconv.Itoa(os.Getpid())
	}
	if a.Ctx.Config.Logging.LogTag == "" {
		a.Ctx.Config.Logging.LogTag = os.Args[0]
	}
	return nil
}

// setupLogger sets apllication Logger up according to configuration settings
func (a *Application) initLogger() {
	if a.Ctx.Config.Logging.LogMode == "nil" || a.Ctx.Config.Logging.LogMode == "null" {
		a.Ctx.Logger = logger.NewNilLogger()
		return
	}
	a.Ctx.Logger = logger.NewStdLogger(a.Ctx.Config.Logging.ID, a.Ctx.Config.Logging.LogTag)
}

// setJWTSigningKey sets key for JWT.
func (a *Application) setJWTSigningKey() error {
	// generate random bytes
	seed := make([]byte, 10)
	_, err := cryptorand.Read(seed) // get 10 crypto random bytes
	if err != nil {
		return errors.New("crypto random byte generation error: " + err.Error())
	}
	hasher := sha256.New()
	_, err = hasher.Write(seed)
	if err != nil {
		return err
	}
	a.Ctx.JWTSignKey = hasher.Sum(nil)
	return nil
}

// init database
func (a *Application) initOrm() error {
	var err error
	// open database
	a.Ctx.Orm, err = xorm.NewEngine(a.Ctx.Config.Database.Db, a.Ctx.Config.Database.Dsn)
	if err != nil {
		return err
	}
	// turn on logs
	ormLogger := logger.NewOrmLogger(a.Ctx.Logger)
	a.Ctx.Orm.SetLogger(ormLogger)
	a.Ctx.Orm.ShowSQL(true)
	// migrate
	err = a.migrateDb()
	if err != nil {
		return err
	}
	// init data
	err = a.initDbData()
	return err
}

// migrate database
func (a *Application) migrateDb() error {
	// migrate tables
	return a.Ctx.Orm.Sync(&users.User{})
}

// initDbData installs hardcoded data from config
func (a *Application) initDbData() error {
	user := &users.User{Email: "admin", DisplayName: "admin", Password: "admin"} // aaaa, backdoor
	err := user.Save(a.Ctx.Orm)
	if err == nil {
		return nil
	}
	status, _ := errors.Decompose(err)
	if status == http.StatusConflict {
		return nil
	}
	err = errors.NewWithPrefix(err, "database error")
	a.Ctx.Logger.Error("application init error", err.Error())
	return err
}
