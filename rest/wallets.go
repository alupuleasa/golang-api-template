package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/efimovalex/wallet/adapters/model"
	"github.com/julienschmidt/httprouter"
)

// CreateWallet - endpoint that creates a wallet
func (r *REST) CreateWallet(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	wallet := model.Wallet{}

	err := json.NewDecoder(req.Body).Decode(&wallet)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer req.Body.Close()

	dbWallet, err := r.DB.CreateWallet(wallet.OwnerAccountID)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	r.WriteJSON(w, http.StatusCreated, dbWallet)

	return
}

// GetWallets - endpoint that retrieves all wallets
func (r *REST) GetWallets(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	dbWallets, err := r.DB.GetWallets(0, 0)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	r.WriteJSON(w, http.StatusOK, dbWallets)
}

// GetWallet - endpoint that retrieves a wallet
func (r *REST) GetWallet(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ID, err := strconv.ParseUint(p.ByName("id"), 10, 64)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbWallet, err := r.DB.GetWallet(ID)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	r.WriteJSON(w, http.StatusOK, dbWallet)
}
