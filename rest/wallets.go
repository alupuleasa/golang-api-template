package rest

import (
	"encoding/json"
	"net/http"

	"github.com/efimovalex/wallet/adapters/model"
	"github.com/julienschmidt/httprouter"
)

// CreateWallet - endpoint that creates a wallet
func (r *REST) CreateWallet(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	wallet := model.Wallet{}
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&w)
	if err != nil {
		r.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbWallet, err := r.DB.CreateWallet(wallet.OwnerAccountID)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(dbWallet)
}

// GetWallets - endpoint that retrieves all wallets
func (r *REST) GetWallets(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	dbWallets, err := r.DB.GetWallets(0, 0)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(dbWallets)
}

// GetWallet - endpoint that retrieves a wallet
func (r *REST) GetWallet(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	dbWallets, err := r.DB.GetWallets(0, 0)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(dbWallets)
}
