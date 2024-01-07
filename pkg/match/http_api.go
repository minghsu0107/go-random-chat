package match

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"gopkg.in/olahol/melody.v1"
)

// @Summary Match another user
// @Description Websocket initialization endpoint for matching another user
// @Tags match
// @Produce json
// @Param Cookie header string true "session id cookie"
// @Failure 401 {object} common.ErrResponse
// @Failure 404 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
// @Router /match [get]
func (r *HttpServer) Match(c *gin.Context) {
	userID, ok := c.Request.Context().Value(common.UserKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	_, err := r.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		r.logger.Error(err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	if err := r.mm.HandleRequest(c.Writer, c.Request); err != nil {
		r.logger.Error("upgrade websocket error: " + err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
}

func (r *HttpServer) HandleMatchOnConnect(sess *melody.Session) {
	userID, ok := sess.Request.Context().Value(common.UserKey).(uint64)
	if !ok {
		r.logger.Error("get user id context error")
		return
	}
	err := r.initializeMatchSession(sess, userID)
	if err != nil {
		r.logger.Error(err.Error())
		return
	}
	ctx := context.Background()
	matchResult, err := r.matchSvc.Match(ctx, userID)
	if err != nil {
		r.logger.Error(err.Error())
		return
	}
	if !matchResult.Matched {
		return
	}
	if err := r.matchSvc.BroadcastMatchResult(ctx, matchResult); err != nil {
		r.logger.Error(err.Error())
		return
	}
}
func (r *HttpServer) initializeMatchSession(sess *melody.Session, userID uint64) error {
	sess.Set(sessUidKey, userID)
	return nil
}
func (r *HttpServer) HandleMatchOnClose(sess *melody.Session, i int, s string) error {
	userID, ok := sess.Request.Context().Value(common.UserKey).(uint64)
	if !ok {
		return nil
	}
	return r.matchSvc.RemoveUserFromWaitList(context.Background(), userID)
}
