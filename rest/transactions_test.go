package rest_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

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

	db.Exec("TRUNCATE wallet,transaction RESTART IDENTITY;")
}

func TestREST_UpdateWalletFunds(t *testing.T) {
	r := initRest()

	w, err := r.DB.CreateWallet(1)
	if err != nil {
		t.Fatalf(err.Error())
	}

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
					httprouter.Param{Key: "id", Value: strconv.FormatUint(w.IDx, 10)},
				},
				body: `
					{
						"sum": 10,
						"reference": "wallet fund"
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"transaction":{"IDx":1,"WalletID":1,"Reference":"wallet fund","Sum":10},"wallet":{"idx":1,"funds":10,"owner_account_id":1}}`,
		},
		{
			name: "success sub 5$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{Key: "id", Value: strconv.FormatUint(w.IDx, 10)},
				},
				body: `
					{
						"sum": -5,
						"reference": "wallet payment"
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"transaction":{"IDx":2,"WalletID":1,"Reference":"wallet payment","Sum":-5},"wallet":{"idx":1,"funds":5,"owner_account_id":1}}`,
		},
		{
			name: "success empty wallet",
			args: args{
				p: httprouter.Params{
					httprouter.Param{Key: "id", Value: strconv.FormatUint(w.IDx, 10)},
				},
				body: `
					{
						"sum": -5,
						"reference": "wallet payment"
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"transaction":{"IDx":3,"WalletID":1,"Reference":"wallet payment","Sum":-5},"wallet":{"idx":1,"funds":0,"owner_account_id":1}}`,
		},
		{
			name: "error withdraw under 0$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{Key: "id", Value: strconv.FormatUint(w.IDx, 10)},
				},
				body: `
					{
						"sum": -110,
						"reference": "wallet payment"
					}`,
			},
			expectedStatus: 500,
			expectedBody:   `{"status":"Internal Server Error","code":500,"message":"the balance of a wallet can not be negative"}`,
		},
		{
			name: "success add 5$",
			args: args{
				p: httprouter.Params{
					httprouter.Param{Key: "id", Value: "1"},
				},
				body: `
					{
						"sum": 5,
						"reference": "wallet fund"
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"transaction":{"IDx":4,"WalletID":1,"Reference":"wallet fund","Sum":5},"wallet":{"idx":1,"funds":5,"owner_account_id":1}}`,
		},
		{
			name: "error withdraw 6$ => funds <= 0",
			args: args{
				p: httprouter.Params{
					httprouter.Param{Key: "id", Value: strconv.FormatUint(w.IDx, 10)},
				},
				body: `
					{
						"sum": -6,
						"reference": "wallet payment"
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

func TestREST_UpdateTransaction(t *testing.T) {
	r := initRest()

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success update reference",
			args: args{
				p: httprouter.Params{
					httprouter.Param{Key: "id", Value: "1"},
				},
				body: `
					{
						"reference": "sum locked for transfer"
					}`,
			},
			expectedStatus: 200,
			expectedBody:   `{"IDx":1,"WalletID":1,"Reference":"sum locked for transfer","Sum":10}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{Body: ioutil.NopCloser(strings.NewReader(tt.args.body))}
			w := httptest.NewRecorder()

			r.UpdateTransaction(w, &req, tt.args.p)

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
