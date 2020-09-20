// +build integration

package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"urlshortner/pkg/app"
	"urlshortner/pkg/config"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/store"
)

const (
	pingAddress    = "http://127.0.0.1:8080/ping"
	shortenAddress = "http://127.0.0.1:8080/shorten"
	longURL        = "wikipedia.com"
	shortURL       = "127.0.0.1:8080/AreVAfnsk"
)

func TestShortenerAPI(t *testing.T) {
	go app.Start()

	cl := getClient()

	waitForServerWithRetry(cl, 10)

	testPingRequest(t, cl)
	testRedirectSuccess(t, cl)
	testDuplication(t, cl)
	testNotPresent(t, cl)
}

func testPingRequest(t *testing.T, cl *http.Client) {
	req := getPingRequest(t)
	resp := doPingRequest(t, cl, req)
	assert.Equal(t, "{\"data\":\"pong\",\"error\":{},\"success\":true}", resp)
}

func testRedirectSuccess(t *testing.T, cl *http.Client) {
	defer cleanUp(t)

	shtReq := getShortenRequest(t, longURL)
	shtResp := doShortenRequest(t, cl, shtReq)

	redReq := getRedirectRequest(t, shtResp.ShortURL)
	redResp := doRedirectRequest(t, cl, redReq)

	assert.Equal(t, http.StatusMovedPermanently, redResp.StatusCode)
	assert.Equal(t, []string{longURL}, redResp.Header["Location"])
}

func testDuplication(t *testing.T, cl *http.Client) {
	defer cleanUp(t)

	shtReqOne := getShortenRequest(t, longURL)
	shtRespOne := doShortenRequest(t, cl, shtReqOne)

	shtReqTwo := getShortenRequest(t, longURL)
	shtRespTwo := doShortenRequest(t, cl, shtReqTwo)

	assert.Equal(t, shtRespOne.ShortURL, shtRespTwo.ShortURL)
}

func testNotPresent(t *testing.T, cl *http.Client) {
	defer cleanUp(t)

	redReq := getRedirectRequest(t, shortURL)
	redResp := doRedirectRequest(t, cl, redReq)

	d, err := ioutil.ReadAll(redResp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, redResp.StatusCode)
	assert.Equal(t, "{\"error\":{\"code\":\"usx0010\",\"message\":\"sql: no rows in result set\"},\"success\":false}", string(d))
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func waitForServerWithRetry(cl *http.Client, retry int) {
	ping := func(cl *http.Client) error {
		resp, err := cl.Get(pingAddress)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New("failed to ping")
		}

		return nil
	}

	for i := 1; i <= retry; i++ {
		fmt.Println(err, i)
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

func doShortenRequest(t *testing.T, cl *http.Client, req *http.Request) contract.ShortenResponse {
	resp := doHTTPRequest(t, cl, req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var apiResp contract.APIResponse
	err = json.Unmarshal(data, &apiResp)
	require.NoError(t, err)

	b, err := json.Marshal(apiResp.Data)
	require.NoError(t, err)

	var shtResp contract.ShortenResponse
	err = json.Unmarshal(b, &shtResp)
	require.NoError(t, err)

	return shtResp
}

func doPingRequest(t *testing.T, cl *http.Client, req *http.Request) string {
	resp := doHTTPRequest(t, cl, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(data)
}

func doRedirectRequest(t *testing.T, cl *http.Client, req *http.Request) *http.Response {
	return doHTTPRequest(t, cl, req)
}

func doHTTPRequest(t *testing.T, cl *http.Client, req *http.Request) *http.Response {
	resp, err := cl.Do(req)
	require.NoError(t, err)

	return resp
}

func getRedirectRequest(t *testing.T, shortURL string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", shortURL), nil)
	require.NoError(t, err)

	return req
}

func getPingRequest(t *testing.T) *http.Request {
	return getHTTPRequest(t, http.MethodGet, pingAddress, nil)
}

func getShortenRequest(t *testing.T, url string) *http.Request {
	shtReq := contract.ShortenRequest{URL: url}
	b, err := json.Marshal(shtReq)
	require.NoError(t, err)

	return getHTTPRequest(t, http.MethodPost, shortenAddress, bytes.NewBuffer(b))
}

func getHTTPRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	return req
}
