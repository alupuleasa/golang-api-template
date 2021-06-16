package rest_test

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/efimovalex/wallet/adapters/database"
	"github.com/efimovalex/wallet/rest"
	"github.com/julienschmidt/httprouter"
)

// These are integration tests since they are running against a real database.
var db *sql.DB

func initRest() (r *rest.REST) {
	r = &rest.REST{}
	r.DB = database.New(db)

	return r
}

type args struct {
	body string
	p    httprouter.Params
}

func TestREST_CreateWallet(t *testing.T) {
	r := initRest()

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success creates wallet",
			args: args{
				body: `
					{
						"owner_account_id": 2
					}`,
				p: nil,
			},
			expectedStatus: 201,
			expectedBody:   `{"idx":2,"funds":0,"owner_account_id":2}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{Body: ioutil.NopCloser(strings.NewReader(tt.args.body))}
			w := httptest.NewRecorder()

			r.CreateWallet(w, &req, tt.args.p)

			if w.Code != tt.expectedStatus {
				t.Errorf("wrong status code, wanted %d, got %d", tt.expectedStatus, w.Code)
			}

			rBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Logf("Body: %s", string(rBody))
				t.Errorf(err.Error())
			}
			rBody = bytes.Trim(rBody, " \n\t")
			if string(rBody) != tt.expectedBody {
				t.Errorf("wrong body, wanted `%s`, got `%s`", tt.expectedBody, rBody)
			}
		})
	}
}

func TestREST_GetWallets(t *testing.T) {
	r := initRest()

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success get all",
			expectedStatus: 200,
			expectedBody:   `[{"idx":1,"funds":5,"owner_account_id":1},{"idx":2,"funds":0,"owner_account_id":2}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{Body: ioutil.NopCloser(strings.NewReader(tt.args.body))}
			w := httptest.NewRecorder()

			r.GetWallets(w, &req, tt.args.p)

			if w.Code != tt.expectedStatus {
				t.Errorf("wrong status code, wanted %d, got %d", tt.expectedStatus, w.Code)
			}

			rBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Logf("Body: %s", string(rBody))
				t.Errorf(err.Error())
			}
			rBody = bytes.Trim(rBody, " \n\t")
			if string(rBody) != tt.expectedBody {
				t.Errorf("wrong body, wanted `%s`, got `%s`", tt.expectedBody, rBody)
			}
		})
	}
}

func TestREST_GetWallet(t *testing.T) {
	r := initRest()

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success get one",
			args: args{p: httprouter.Params{
				httprouter.Param{Key: "id", Value: "2"},
			}},
			expectedStatus: 200,
			expectedBody:   `{"idx":2,"funds":0,"owner_account_id":2}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{Body: ioutil.NopCloser(strings.NewReader(tt.args.body))}
			w := httptest.NewRecorder()

			r.GetWallet(w, &req, tt.args.p)

			if w.Code != tt.expectedStatus {
				t.Errorf("wrong status code, wanted %d, got %d", tt.expectedStatus, w.Code)
			}

			rBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Logf("Body: %s", string(rBody))
				t.Errorf(err.Error())
			}
			rBody = bytes.Trim(rBody, " \n\t")
			if string(rBody) != tt.expectedBody {
				t.Errorf("wrong body, wanted `%s`, got `%s`", tt.expectedBody, rBody)
			}
		})
	}
}
