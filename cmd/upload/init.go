package upload

import "github.com/gin-gonic/gin"

func InitializeRouter() *Router {
	svr := gin.Default()
	svr.Use(CORSMiddleware())
	return NewRouter(svr)
}
