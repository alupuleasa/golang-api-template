package rest

import (
	"github.com/julienschmidt/httprouter"
)

func (r *REST) loadRoutes() {
	r.router = httprouter.New()

	r.router.GET("/healthcheck", r.Healthcheck)

	r.router.POST("/wallets", r.Auth(r.CreateWallet))
	r.router.GET("/wallets", r.Auth(r.GetWallets))
	r.router.GET("/wallets/:id", r.Auth(r.GetWallet))

	r.router.PATCH("/wallets/:id/funds", r.Auth(r.UpdateWalletFunds))
	r.router.PATCH("/transaction/:id/reference", r.Auth(r.UpdateTransaction))
}
