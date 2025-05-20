package api

import (
	"fmt"
	"github.com/dsrosen6/cw-ticket-bot/internal/util"
	"github.com/dsrosen6/cw-ticket-bot/internal/webex"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) DirectMessageRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", s.handleReceiveDm)
	return r
}

func (s *Server) handleReceiveDm(w http.ResponseWriter, r *http.Request) {
	var payload webex.MessageWebhookBody
	if err := util.ParseJSON(r, &payload); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := s.webexClient.GetMessage(r.Context(), payload.Data.Id)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if s.isWebexBot(m) {
		util.WriteJSON(w, http.StatusOK, util.ResultBody("message is from webex bot - no action needed"))
		return
	}

	msgText, err := getMessageText(m)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	newMessage := webex.NewMessageToPerson(m.PersonEmail, fmt.Sprintf("You said:\n%s", msgText))
	if err := s.webexClient.SendMessage(r.Context(), newMessage); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, util.ResultBody("successfully sent message"))
}

// getMessageText pulls the text from a Webex message, since it could be Markdown or plain text.
func getMessageText(msg *webex.MessageGetResponse) (string, error) {
	if msg.Markdown != "" {
		return msg.Markdown, nil
	}

	if msg.Text != "" {
		return msg.Text, nil
	}

	return "", fmt.Errorf("no markdown or text found in message %s", msg.Id)
}

func (s *Server) isWebexBot(msg *webex.MessageGetResponse) bool {
	return msg.PersonEmail == s.botEmail
}
