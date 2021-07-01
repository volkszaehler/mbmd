package server_test

import (
	"io/ioutil"
	"net/http"
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
	httpd := server.NewHttpd(qe, cache)
	go httpd.Run(hub, status, "localhost:8080")

	for {
		resp, err := http.Get("http://localhost:8080/")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (suite *HTTPTestSuite) TestAccessAssetsFromRoot() {
	resp, err := http.Get("http://localhost:8080/css/app.css")
	suite.Require().NoError(err)
	suite.Equal(resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.NotEmpty(body)
}

func (suite *HTTPTestSuite) TestAccessAssetsFromCSS() {
	resp, err := http.Get("http://localhost:8080/css/css/app.css")
	suite.Require().NoError(err)
	suite.Equal(resp.StatusCode, http.StatusNotFound)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.Equal("404 page not found\n", string(body))
}

func TestHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}
