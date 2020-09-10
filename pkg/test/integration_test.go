package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"urlshortner/cmd/config"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/http/server"
	"urlshortner/pkg/store"
)

const address = "http://127.0.0.1:8080"

func TestShortenerAPI(t *testing.T) {
	srv := getServer()
	go srv.Start()

	//waitForServer()

	cl := getClient()

	waitForServerWithRetry(cl, 10)

	testPingRequest(t, cl)
	testRedirectSuccess(t, cl)
	testDuplication(t, cl)
	testNotPresent(t, cl)
}

func testPingRequest(t *testing.T, cl *http.Client) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", address, "ping"), nil)
	require.NoError(t, err)

	resp, err := cl.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "pong", string(data))
}

func testRedirectSuccess(t *testing.T, cl *http.Client) {
	defer cleanUp(t)

	shtReq := shortenRequest(t, "wikipedia.com")
	shtResp := doShortenRequest(t, cl, shtReq)

	redReq := redirectRequest(t, shtResp.ShortURL)
	redResp := doRedirectRequest(t, cl, redReq)

	assert.Equal(t, http.StatusMovedPermanently, redResp.StatusCode)
	assert.Equal(t, []string{"wikipedia.com"}, redResp.Header["Location"])
}

func testDuplication(t *testing.T, cl *http.Client) {
	defer cleanUp(t)

	shtReqOne := shortenRequest(t, "wikipedia.com")
	shtRespOne := doShortenRequest(t, cl, shtReqOne)

	shtReqTwo := shortenRequest(t, "wikipedia.com")
	shtRespTwo := doShortenRequest(t, cl, shtReqTwo)

	assert.Equal(t, shtRespOne.ShortURL, shtRespTwo.ShortURL)
}

func testNotPresent(t *testing.T, cl *http.Client) {
	defer cleanUp(t)

	redReq := redirectRequest(t, "127.0.0.1:8080/AreVAfnsk")
	redResp := doRedirectRequest(t, cl, redReq)

	d, err := ioutil.ReadAll(redResp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, redResp.StatusCode)
	assert.Equal(t, "sql: no rows in result set", string(d))
}

func shortenRequest(t *testing.T, url string) *http.Request {
	shtReq := contract.ShortenRequest{URL: url}
	b, err := json.Marshal(shtReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", address, "shorten"), bytes.NewBuffer(b))
	require.NoError(t, err)

	return req
}

func doShortenRequest(t *testing.T, cl *http.Client, req *http.Request) contract.ShortenResponse {
	resp := doHTTPRequest(t, cl, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var shtResp contract.ShortenResponse
	err = json.Unmarshal(data, &shtResp)
	require.NoError(t, err)

	return shtResp
}

func doRedirectRequest(t *testing.T, cl *http.Client, req *http.Request) *http.Response {
	return doHTTPRequest(t, cl, req)
}

func doHTTPRequest(t *testing.T, cl *http.Client, req *http.Request) *http.Response {
	resp, err := cl.Do(req)
	require.NoError(t, err)

	return resp
}

func redirectRequest(t *testing.T, shortURL string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", shortURL), nil)
	require.NoError(t, err)

	return req
}

func getServer() server.Server {
	return server.NewServer(
		config.NewConfig(),
		zap.NewNop(),
		&newrelic.Application{},
		&statsd.Client{},
	)
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func waitForServer() {
	time.Sleep(time.Second)
}

func waitForServerWithRetry(cl *http.Client, retry int) {
	ping := func(cl *http.Client) error {
		resp, err := cl.Get(fmt.Sprintf("%s/%s", address, "ping"))
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New("failed to ping")
		}

		return nil
	}

	for i := 1; i <=retry; i++ {
		err := ping(cl)
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * time.Duration(retry))
		}
	}
}

func cleanUp(t *testing.T) {
	dbHandler := store.NewDBHandler(config.NewConfig().GetDatabaseConfig(), zap.NewNop())
	db, err := dbHandler.GetDB()
	require.NoError(t, err)

	_, err = db.Exec("TRUNCATE shortener")
	require.NoError(t, err)

	cacheHandler := store.NewCacheHandler(config.NewConfig().GetRedisConfig(), zap.NewNop())
	redis, err := cacheHandler.GetCache()
	require.NoError(t, err)

	cmd := redis.FlushAll(context.Background())
	require.NoError(t, cmd.Err())
}
