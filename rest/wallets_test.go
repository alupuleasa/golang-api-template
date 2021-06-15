package rest

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/efimovalex/wallet/adapters/database"
	"github.com/julienschmidt/httprouter"
)

// These are integration tests since they are running against a real database.
var db *sql.DB

func init() {
	var err error
	// Hardcode test env database so there is no chance of it reaching other environments
	// Different user & different database from the local dev one
	db, err = sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s/%s", "pgx_test", "test_user_2021", "database", "wallet_test"))
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	db.Exec("TRUNCATE wallet RESTART IDENTITY ;")
}
func initRest() (r *REST) {
	r = &REST{}
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
						"owner_account_id": 1
					}`,
				p: nil,
			},
			expectedStatus: 201,
			expectedBody:   `{"idx":1,"funds":0,"owner_account_id":1}`,
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
			expectedBody:   `[{"idx":1,"funds":0,"owner_account_id":1}]`,
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
				httprouter.Param{"id", "1"},
			}},
			expectedStatus: 200,
			expectedBody:   `{"idx":1,"funds":0,"owner_account_id":1}`,
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

func TestREST_UpdateWalletFunds(t *testing.T) {
	r := initRest()

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success add 10$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{"id", "1"},
				},
				body: `
					{
						"sum": 10
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"idx":1,"funds":10,"owner_account_id":1}`,
		},
		{
			name: "success sub 5$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{"id", "1"},
				},
				body: `
					{
						"sum": -5
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"idx":1,"funds":5,"owner_account_id":1}`,
		},
		{
			name: "success empty wallet",
			args: args{
				p: httprouter.Params{
					httprouter.Param{"id", "1"},
				},
				body: `
					{
						"sum": -5
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"idx":1,"funds":0,"owner_account_id":1}`,
		},
		{
			name: "error withdraw under 0$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{"id", "1"},
				},
				body: `
					{
						"sum": -110
					}`,
			},
			expectedStatus: 500,
			expectedBody:   `{"status":"Internal Server Error","code":500,"message":"the balance of a wallet can not be negative"}`,
		},
		{
			name: "success add 5$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{"id", "1"},
				},
				body: `
					{
						"sum": 5
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"idx":1,"funds":5,"owner_account_id":1}`,
		},
		{
			name: "error withdraw 6$ => funds <= 0",
			args: args{
				p: httprouter.Params{
					httprouter.Param{"id", "1"},
				},
				body: `
					{
						"sum": -6
					}`,
			},
			expectedStatus: 500,
			expectedBody:   `{"status":"Internal Server Error","code":500,"message":"the balance of a wallet can not be negative"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{Body: ioutil.NopCloser(strings.NewReader(tt.args.body))}
			w := httptest.NewRecorder()

			r.UpdateWalletFunds(w, &req, tt.args.p)

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
