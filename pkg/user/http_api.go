package user

import (
	"errors"
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
// @Failure 400 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
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
	sid, err := r.userSvc.SetUserSession(c.Request.Context(), user.ID)
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	common.SetAuthCookie(c, sid, r.authCookieConfig.MaxAge, r.authCookieConfig.Path, r.authCookieConfig.Domain)

	c.JSON(http.StatusCreated, &UserPresenter{
		ID:   strconv.FormatUint(user.ID, 10),
		Name: user.Name,
	})
}

// @Summary Get user
// @Description Get user information
// @Tags user
// @Produce json
// @Param uid query string false "target user id; defaults to self user id if leaved empty"
// @Param Cookie header string true "session id cookie"
// @Success 200 {object} UserPresenter
// @Failure 400 {object} common.ErrResponse
// @Failure 401 {object} common.ErrResponse
// @Failure 404 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
// @Router /user [get]
func (r *HttpServer) GetUser(c *gin.Context) {
	userID, ok := c.Request.Context().Value(common.UserKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	uid := c.Query("uid")
	if uid != "" {
		var err error
		userID, err = strconv.ParseUint(uid, 10, 64)
		if err != nil {
			response(c, http.StatusBadRequest, common.ErrInvalidParam)
			return
		}
	}
	user, err := r.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	c.JSON(http.StatusOK, &UserPresenter{
		ID:   strconv.FormatUint(userID, 10),
		Name: user.Name,
	})
}

// @Summary OAuth Google login
// @Description OAuth Google login endpoint
// @Tags user
// @Success 307
// @Router /user/oauth2/google/login [get]
func (r *HttpServer) OAuthGoogleLogin(c *gin.Context) {
	oauthState := common.GenerateStateOauthCookie(c, r.oauthCookieConfig.MaxAge, r.oauthCookieConfig.Path, r.oauthCookieConfig.Domain)
	u := r.googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

// @Summary OAuth Google callback
// @Description OAuth Google callback endpoint
// @Tags user
// @Success 307
// @Failure 500 {object} common.ErrResponse
// @Router /user/oauth2/google/callback [get]
func (r *HttpServer) OAuthGoogleCallback(c *gin.Context) {
	oauthState, err := common.GetCookie(c, common.OAuthStateCookieName)
	if err != nil {
		r.logger.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	if c.Query("state") != oauthState {
		r.logger.Error("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	token, err := r.googleOauthConfig.Exchange(c.Request.Context(), c.Request.FormValue("code"))
	if err != nil {
		r.logger.Errorf("code exchange wrong: %w", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	googleUser, err := r.userSvc.GetUserDataFromGoogle(c.Request.Context(), token.AccessToken)
	if err != nil {
		r.logger.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	user, err := r.userSvc.GetOrCreateGoogleUser(c.Request.Context(), googleUser.Email, googleUser.Name)
	if err != nil {
		r.logger.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	sid, err := r.userSvc.SetUserSession(c.Request.Context(), user.ID)
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	common.SetAuthCookie(c, sid, r.authCookieConfig.MaxAge, r.authCookieConfig.Path, r.authCookieConfig.Domain)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
