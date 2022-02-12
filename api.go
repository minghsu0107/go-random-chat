package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

func (r *Router) Match(c *gin.Context) {
	uid := c.Query("uid")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response(c, http.StatusBadRequest, ErrInvalidParam)
		return
	}
	_, err = r.userSvc.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == ErrUserNotFound {
			response(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	r.mm.HandleRequest(c.Writer, c.Request)
}

func (r *Router) StartChat(c *gin.Context) {
	uid := c.Query("uid")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response(c, http.StatusBadRequest, ErrInvalidParam)
		return
	}
	accessToken := c.Query("access_token")
	authResult, err := Auth(&AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	if authResult.Expired {
		log.Error(ErrTokenExpired)
		response(c, http.StatusUnauthorized, ErrTokenExpired)
	}
	channelID := authResult.ChannelID
	exist, err := r.userSvc.IsChannelUserExist(c.Request.Context(), channelID, userID)
	if err != nil {
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	if !exist {
		response(c, http.StatusNotFound, ErrChannelOrUserNotFound)
		return
	}

	r.mc.HandleRequest(c.Writer, c.Request)
}

func (r *Router) CreateUser(c *gin.Context) {
	var userPresenter UserPresenter
	if err := c.ShouldBindJSON(&userPresenter); err != nil {
		response(c, http.StatusBadRequest, ErrInvalidParam)
		return
	}
	user, err := r.userSvc.CreateUser(c.Request.Context(), userPresenter.Name)
	if err != nil {
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	c.JSON(http.StatusCreated, &UserPresenter{
		ID:   strconv.FormatUint(user.ID, 10),
		Name: user.Name,
	})
}

func (r *Router) GetChannelUsers(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	userIDs, err := r.userSvc.GetChannelUserIDs(c.Request.Context(), channelID)
	if err != nil {
		log.Error(err)
		if err == ErrChannelNotFound {
			response(c, http.StatusNotFound, ErrChannelNotFound)
			return
		}
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	userIDsPresenter := []string{}
	for _, userID := range userIDs {
		userIDsPresenter = append(userIDsPresenter, strconv.FormatUint(userID, 10))
	}
	c.JSON(http.StatusOK, &UserIDsPresenter{
		UserIDs: userIDsPresenter,
	})
}

func (r *Router) GetOnlineUsers(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	userIDs, err := r.userSvc.GetOnlineUserIDs(c.Request.Context(), channelID)
	if err != nil {
		if err == ErrChannelNotFound {
			response(c, http.StatusNotFound, ErrChannelNotFound)
			return
		}
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	userIDsPresenter := []string{}
	for _, userID := range userIDs {
		userIDsPresenter = append(userIDsPresenter, strconv.FormatUint(userID, 10))
	}
	c.JSON(http.StatusOK, &UserIDsPresenter{
		UserIDs: userIDsPresenter,
	})
}

func (r *Router) GetUserName(c *gin.Context) {
	id := c.Param("uid")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response(c, http.StatusBadRequest, ErrInvalidParam)
		return
	}
	user, err := r.userSvc.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == ErrUserNotFound {
			response(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	c.JSON(http.StatusOK, &UserPresenter{
		ID:   id,
		Name: user.Name,
	})
}

func (r *Router) ListMessages(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	msgs, err := r.msgSvc.ListMessages(c.Request.Context(), channelID)
	if err != nil {
		log.Error(err)
		if err == ErrChannelNotFound {
			response(c, http.StatusNotFound, ErrChannelNotFound)
			return
		}
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	msgsPresenter := []MessagePresenter{}
	for _, msg := range msgs {
		msgsPresenter = append(msgsPresenter, MessagePresenter{
			Event:   msg.Event,
			UserID:  strconv.FormatUint(msg.UserID, 10),
			Payload: msg.Payload,
			Time:    msg.Time,
		})
	}
	c.JSON(http.StatusOK, &MessagesPresenter{
		Messages: msgsPresenter,
	})
}

func (r *Router) DeleteChannel(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	uid := c.Query("delby")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response(c, http.StatusBadRequest, ErrInvalidParam)
		return
	}

	exist, err := r.userSvc.IsChannelUserExist(c.Request.Context(), channelID, userID)
	if err != nil {
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	if !exist {
		response(c, http.StatusBadRequest, ErrChannelOrUserNotFound)
		return
	}

	err = r.msgSvc.BroadcastActionMessage(c.Request.Context(), channelID, userID, LeavedMessage)
	if err != nil {
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	err = r.chanSvc.DeleteChannel(c.Request.Context(), channelID)
	if err != nil {
		log.Error(err)
		response(c, http.StatusInternalServerError, ErrServer)
		return
	}
	c.JSON(http.StatusNoContent, SuccessMessage{
		Message: "ok",
	})
}

func (r *Router) HandleMatchOnConnect(sess *melody.Session) {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		log.Error(err)
		return
	}
	err = r.initializeMatchSession(sess, userID)
	if err != nil {
		log.Error(err)
		return
	}
	ctx := context.Background()
	matchResult, err := r.matchSvc.Match(ctx, userID)
	if err != nil {
		log.Error(err)
		return
	}
	if !matchResult.Matched {
		return
	}
	if err := r.matchSvc.BroadcastMatchResult(ctx, matchResult); err != nil {
		log.Error(err)
		return
	}
}
func (r *Router) initializeMatchSession(sess *melody.Session, userID uint64) error {
	sess.Set(sessUidKey, userID)
	return nil
}
func (r *Router) HandleMatchOnClose(sess *melody.Session, i int, s string) error {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	return r.matchSvc.RemoveUserFromWaitList(context.Background(), userID)
}

func (r *Router) HandleChatOnConnect(sess *melody.Session) {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		log.Error(err)
		return
	}
	accessToken := sess.Request.URL.Query().Get("access_token")
	authResult, err := Auth(&AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Error(err)
	}
	if authResult.Expired {
		log.Error(ErrTokenExpired)
	}
	channelID := authResult.ChannelID
	err = r.initializeChatSession(sess, channelID, userID)
	if err != nil {
		log.Error(err)
		return
	}
	if err := r.msgSvc.BroadcastConnectMessage(context.Background(), channelID, userID); err != nil {
		log.Error(err)
		return
	}
}

func (r *Router) initializeChatSession(sess *melody.Session, channelID, userID uint64) error {
	ctx := context.Background()
	if err := r.userSvc.AddOnlineUser(ctx, channelID, userID); err != nil {
		return err
	}
	sess.Set(sessCidKey, channelID)
	return nil
}

func (r *Router) HandleChatOnMessage(sess *melody.Session, data []byte) {
	msgPresenter, err := DecodeToMessagePresenter(data)
	if err != nil {
		log.Error(err)
		return
	}
	msg, err := msgPresenter.ToMessage(sess.Request.URL.Query().Get("access_token"))
	if err != nil {
		log.Error(err)
		return
	}
	switch msg.Event {
	case EventText:
		if err := r.msgSvc.BroadcastTextMessage(context.Background(), msg.ChannelID, msg.UserID, msg.Payload, msg.Time); err != nil {
			log.Error(err)
		}
	case EventAction:
		if err := r.msgSvc.BroadcastActionMessage(context.Background(), msg.ChannelID, msg.UserID, Action(msg.Payload)); err != nil {
			log.Error(err)
		}
	default:
		log.Errorf("invailid event type: %v", msg.Event)
	}
}

func (r *Router) HandleChatOnClose(sess *melody.Session, i int, s string) error {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	accessToken := sess.Request.URL.Query().Get("access_token")
	authResult, err := Auth(&AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	if authResult.Expired {
		log.Error(ErrTokenExpired)
		return ErrTokenExpired
	}
	channelID := authResult.ChannelID
	err = r.userSvc.DeleteOnlineUser(context.Background(), channelID, userID)
	if err != nil {
		log.Error(err)
		return err
	}
	return r.msgSvc.BroadcastActionMessage(context.Background(), channelID, userID, OfflineMessage)
}
