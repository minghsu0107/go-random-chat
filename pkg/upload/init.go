package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
)

func InitializeRouter() *Router {
	svr := gin.Default()
	svr.Use(common.CORSMiddleware())
	return NewRouter(svr)
}
