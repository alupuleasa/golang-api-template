package rest

import (
	"github.com/julienschmidt/httprouter"
)

func (r *REST) loadRoutes() {
	r.router = httprouter.New()

	r.router.GET("/healthcheck", r.Healthcheck)

	r.router.POST("/wallet", r.Auth(r.CreateWallet))
	r.router.GET("/wallets", r.Auth(r.GetWallets))
	r.router.GET("/wallet/:id", r.Auth(r.GetWallet))

	r.router.PATCH("/wallet/:id/funds", r.Auth(r.UpdateWalletFunds))
}
