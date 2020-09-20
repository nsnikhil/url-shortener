package util_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"urlshortner/pkg/http/util"
)

func TestParseRequest(t *testing.T) {
	type CusReq struct {
		ReqID   string `json:"req_id"`
		ReqData string `json:"req_data"`
	}

	type Other struct {
		ReqData int `json:"req_data"`
	}

	testCases := []struct {
		name           string
		actualResult   func() (error, interface{})
		expectedResult interface{}
		expectedError  error
	}{
		{
			name: "test request parse success",
			actualResult: func() (error, interface{}) {
				cr := CusReq{ReqID: "req-id", ReqData: "req data"}

				b, err := json.Marshal(&cr)
				require.NoError(t, err)

				r, err := http.NewRequest(http.MethodGet, "/random", bytes.NewBuffer(b))
				require.NoError(t, err)

				var tr CusReq
				return util.ParseRequest(r, &tr), tr
			},
			expectedResult: CusReq{ReqID: "req-id", ReqData: "req data"},
		},
		{
			name: "test request parse failure when req is nil",
			actualResult: func() (error, interface{}) {
				var tr CusReq
				return util.ParseRequest(nil, &tr), tr
			},
			expectedResult: CusReq{},
			expectedError:  errors.New("request is nil"),
		},
		{
			name: "test request parse failure when req body is nil",
			actualResult: func() (error, interface{}) {
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				var tr CusReq
				return util.ParseRequest(r, &tr), tr
			},
			expectedResult: CusReq{},
			expectedError:  errors.New("request body is nil"),
		},
		{
			name: "test request parse failure when fail to read body",
			actualResult: func() (error, interface{}) {
				cr := CusReq{ReqID: "req-id", ReqData: "req data"}

				b, err := json.Marshal(&cr)
				require.NoError(t, err)

				r, err := http.NewRequest(http.MethodGet, "/random", bytes.NewBuffer(b))
				require.NoError(t, err)

				_, err = ioutil.ReadAll(r.Body)
				require.NoError(t, err)

				var tr CusReq
				return util.ParseRequest(r, &tr), tr
			},
			expectedResult: CusReq{},
			expectedError:  errors.New("unexpected end of JSON input"),
		},
		{
			name: "test request parse failure when unmarshalling fails",
			actualResult: func() (error, interface{}) {
				type CusReq struct {
					ReqID   string   `json:"req_id"`
					ReqData []string `json:"req_data"`
				}

				cr := CusReq{ReqID: "req-id", ReqData: []string{"req data"}}

				b, err := json.Marshal(&cr)
				require.NoError(t, err)

				r, err := http.NewRequest(http.MethodGet, "/random", bytes.NewBuffer(b))
				require.NoError(t, err)

				var tr Other

				return util.ParseRequest(r, &tr), tr
			},
			expectedResult: Other{},
			expectedError:  errors.New("json: cannot unmarshal array into Go struct field Other.req_data of type int"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err, res := testCase.actualResult()

			if testCase.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			}

			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}
