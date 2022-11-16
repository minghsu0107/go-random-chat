package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

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
// @Failure 400 {object} common.ErrResponse
// @Failure 404 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
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
		if errors.Is(err, ErrUserNotFound) {
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

func (r *HttpServer) OAuthGoogleLogin(c *gin.Context) {
	oauthState := r.generateStateOauthCookie(c)
	u := r.googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (r *HttpServer) OAuthGoogleCallback(c *gin.Context) {
	oauthState, err := c.Request.Cookie("oauthstate")
	if err != nil {
		r.logger.Errorf("oauth google state cookie error: %w", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	unescapedState, err := url.QueryUnescape(oauthState.Value)
	if err != nil {
		r.logger.Errorf("unescape oauth google state error: %w", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	if c.Query("state") != unescapedState {
		r.logger.Error("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	googleUser, err := r.getUserDataFromGoogle(c.Request.Context(), c.Request.FormValue("code"))
	if err != nil {
		r.logger.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	user, err := r.userSvc.GetOrCreateGoogleUser(c.Request.Context(), googleUser.Email, googleUser.Name)
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, common.Join("/?user_id=", strconv.FormatUint(user.ID, 10)))
}

func (r *HttpServer) generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 3600, "/", r.cookieDomain, false, true)
	return state
}

func (r *HttpServer) getUserDataFromGoogle(ctx context.Context, code string) (*GoogleUserPresenter, error) {
	// Use code to get token and get user info from Google.

	token, err := r.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %w", err)
	}

	req, err := http.NewRequest("GET", common.Join(oauthGoogleUrlAPI, token.AccessToken), nil)
	if err != nil {
		return nil, fmt.Errorf("create http request error: %w", err)
	}
	req = req.WithContext(ctx)

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %w", err)
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read google user response: %w", err)
	}
	var googleUser GoogleUserPresenter
	if err := json.Unmarshal(contents, &googleUser); err != nil {
		return nil, fmt.Errorf("failed marshal google user response: %w", err)
	}
	return &googleUser, nil
}
