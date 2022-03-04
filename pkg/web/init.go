package web

import "github.com/gin-gonic/gin"

func InitializeRouter() *Router {
	svr := gin.Default()
	return NewRouter(svr)
}
