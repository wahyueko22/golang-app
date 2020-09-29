package bdd_test

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/resty.v1"
	"gopkg.in/testfixtures.v2"

	"github.com/corvinusz/echo-xorm/app"
	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/app/server/auth"
)

// TestData defines format of data to privide for POST/PUT tests
type TestData struct {
	Comment       string
	JsonIn        string
	JsonOut       string
	ID            int
	HttpCode      int
	HaveToCheckDb bool
}

func TestBddtests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bddtests Suite")
}

var suite *EchoTestSuite

var _ = BeforeSuite(func() {
	suite = new(EchoTestSuite)
	err := suite.setupSuite()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	suite.app.Shutdown() //nolint
})

const (
	cfgFileName    = "./test-config/echo-xorm-test-config.toml"
	fixturesFolder = "./fixtures"
)

// EchoTestSuite is testing context for app
type EchoTestSuite struct {
	app     *app.Application
	baseURL string
	rc      *resty.Client
}

// SetupTest called once before test
func (s *EchoTestSuite) setupSuite() error {
	err := s.setupServer()
	if err != nil {
		return err
	}
	s.baseURL = "http://localhost:" + s.app.Ctx.Config.Port
	// create and setup resty client
	s.rc = resty.DefaultClient
	s.rc.SetHeader("Content-Type", "application/json")
	s.rc.SetHostURL(s.baseURL)
	// get auth token
	return s.authorizeMe("admin", "admin")
}

//------------------------------------------------------------------------------
func (s *EchoTestSuite) setupServer() error {
	var err error
	// init test application
	s.app, err = app.New(&ctx.Flags{CfgFileName: cfgFileName})
	if err != nil {
		return err
	}
	// load fixtures
	err = s.setupFixtures()
	if err != nil {
		return err
	}
	// start test server with go routine
	go s.app.Run()
	// wait til server started then return
	return s.waitServerStart(3 * time.Second)
}

//------------------------------------------------------------------------------
func (s *EchoTestSuite) setupFixtures() error {
	fixtures, err := testfixtures.NewFolder(s.app.Ctx.Orm.DB().DB, &testfixtures.SQLite{}, fixturesFolder)
	if err != nil {
		return err
	}
	return fixtures.Load()
}

//------------------------------------------------------------------------------
func (s *EchoTestSuite) waitServerStart(timeout time.Duration) error {
	const sleepTime = 300 * time.Millisecond
	dialer := &net.Dialer{
		DualStack: false,
		Deadline:  time.Now().Add(timeout),
		Timeout:   sleepTime,
		KeepAlive: 0,
	}
	done := time.Now().Add(timeout)
	for time.Now().Before(done) {
		c, err := dialer.Dial("tcp", ":"+s.app.Ctx.Config.Port)
		if err == nil {
			return c.Close()
		}
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("cannot connect %v for %v", s.baseURL, timeout)
}

//------------------------------------------------------------------------------
func (s *EchoTestSuite) authorizeMe(email, password string) error {
	// make authorization
	payload := auth.PostBody{
		Email:    email,
		Password: password,
	}
	result := new(auth.Result)
	response, err := s.rc.R().SetBody(payload).SetResult(result).Post("/auth")
	if err != nil {
		return err
	}

	// check response and set token
	if response.StatusCode() != 200 {
		return errors.New("auth response status is not 200 (not OK)")
	}
	// set auth token
	s.rc.SetAuthToken(result.Token)
	// return
	return nil
}
