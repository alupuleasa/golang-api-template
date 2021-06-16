package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// UpdateWalletFunds - endpoint updates the wallets funds by the requested sum
func (r *REST) UpdateWalletFunds(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ID, err := strconv.ParseUint(p.ByName("id"), 10, 64)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	sum := struct {
		Sum       float64 `json:"sum"`
		Reference string  `json:"reference"`
	}{}
	defer req.Body.Close()

	err = json.NewDecoder(req.Body).Decode(&sum)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	wallet, transaction, err := r.DB.UpdateWalletFunds(ID, sum.Sum, sum.Reference)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	r.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"wallet":      wallet,
		"transaction": transaction,
	})
}

// UpdateTransaction - endpoint updates the transaction with new reference
func (r *REST) UpdateTransaction(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ID, err := strconv.ParseUint(p.ByName("id"), 10, 64)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	sum := struct {
		Reference string `json:"reference"`
	}{}
	defer req.Body.Close()

	err = json.NewDecoder(req.Body).Decode(&sum)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := r.DB.UpdateTransaction(ID, sum.Reference)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	r.WriteJSON(w, http.StatusOK, transaction)
}
