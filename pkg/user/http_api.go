package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
)

// @Summary Create an user
// @Description Register a new user
// @Tags user
// @Produce json
// @Param user body CreateUserRequest true "new user"
// @Success 201 {object} UserPresenter
// @Failure 400 {none} nil
// @Failure 500 {none} nil
// @Router /user [post]
func (r *HttpServer) CreateUser(c *gin.Context) {
	var createUserReq CreateUserRequest
	if err := c.ShouldBindJSON(&createUserReq); err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	user, err := r.userSvc.CreateUser(c.Request.Context(), createUserReq.Name)
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	c.JSON(http.StatusCreated, &UserPresenter{
		ID:   strconv.FormatUint(user.ID, 10),
		Name: user.Name,
	})
}

// @Summary Get user name
// @Description Get user name
// @Tags user
// @Produce json
// @Param uid path int true "user id"
// @Success 200 {object} UserPresenter
// @Failure 400 {none} nil
// @Failure 404 {none} nil
// @Failure 500 {none} nil
// @Router /user/{uid}/name [get]
func (r *HttpServer) GetUserName(c *gin.Context) {
	id := c.Param("uid")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	user, err := r.userSvc.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == ErrUserNotFound {
			response(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	c.JSON(http.StatusOK, &UserPresenter{
		ID:   id,
		Name: user.Name,
	})
}
