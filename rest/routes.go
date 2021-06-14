package rest

func (r *REST) loadRoutes() {
	r.router.GET("/", Index)
	r.router.GET("/hello/:name", Hello)
}
