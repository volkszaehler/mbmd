package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/server"
)

type AssetTestSuite struct {
	suite.Suite
}

func (suite *AssetTestSuite) TestOpenAsset() {
	server.Assets = os.DirFS("../assets")
	fh, err := server.Assets.Open("css/app.css")
	suite.Require().NoError(err)
	defer fh.Close()

	app_cs, err := ioutil.ReadAll(fh)
	suite.NoError(err)
	suite.NotEmpty(app_cs)
}

func TestAssetTestSuite(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

type HTTPTestSuite struct {
	suite.Suite
	httpd *server.Httpd
}

func (suite *HTTPTestSuite) SetupSuite() {
	server.Assets = os.DirFS("../assets")

	cc := make(chan server.ControlSnip)
	teeC := server.NewBroadcaster(server.FromControlChannel(cc))
	go teeC.Run()

	rc := make(chan server.QuerySnip)
	tee := server.NewBroadcaster(server.FromSnipChannel(rc))
	go tee.Run()

	qe := server.NewQueryEngine(make(map[string]*meters.Manager))
	status := server.NewStatus(qe, server.ToControlChannel(teeC.Attach()))
	cache := server.NewCache(time.Minute, status, viper.GetBool("verbose"))
	tee.AttachRunner(server.NewSnipRunner(cache.Run))
	hub := server.NewSocketHub(status)
	tee.AttachRunner(server.NewSnipRunner(hub.Run))
	suite.httpd = server.NewHttpd(hub, status, qe, cache)
}

func (suite *HTTPTestSuite) TestAccessAssetsFromRoot() {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/css/app.css", nil)
	suite.Require().NoError(err)

	suite.httpd.Router().ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.NotEmpty(rr.Body.String())
}

func (suite *HTTPTestSuite) TestAccessAssetsFromCSS() {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/css/css/app.css", nil)
	suite.Require().NoError(err)

	suite.httpd.Router().ServeHTTP(rr, req)

	suite.Equal(http.StatusNotFound, rr.Code)
	suite.Equal("404 page not found\n", rr.Body.String())
}

func TestHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}
