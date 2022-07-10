package match

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"gopkg.in/olahol/melody.v1"
)

func (r *HttpServer) ForwardAuth(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	c.Writer.Header().Set(common.ChannelIdHeader, strconv.FormatUint(channelID, 10))
	c.Status(http.StatusOK)
}

func (r *HttpServer) Match(c *gin.Context) {
	uid := c.Query("uid")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	_, err = r.userSvc.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == ErrUserNotFound {
			response(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	r.mm.HandleRequest(c.Writer, c.Request)
}

func (r *HttpServer) HandleMatchOnConnect(sess *melody.Session) {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		r.logger.Error(err)
		return
	}
	err = r.initializeMatchSession(sess, userID)
	if err != nil {
		r.logger.Error(err)
		return
	}
	ctx := context.Background()
	matchResult, err := r.matchSvc.Match(ctx, userID)
	if err != nil {
		r.logger.Error(err)
		return
	}
	if !matchResult.Matched {
		return
	}
	if err := r.matchSvc.BroadcastMatchResult(ctx, matchResult); err != nil {
		r.logger.Error(err)
		return
	}
}
func (r *HttpServer) initializeMatchSession(sess *melody.Session, userID uint64) error {
	sess.Set(sessUidKey, userID)
	return nil
}
func (r *HttpServer) HandleMatchOnClose(sess *melody.Session, i int, s string) error {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return r.matchSvc.RemoveUserFromWaitList(context.Background(), userID)
}
