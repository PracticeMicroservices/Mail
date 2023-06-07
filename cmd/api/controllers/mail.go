package controllers

import (
	"mail/cmd/api/helpers"
	"mail/cmd/entities"
	"mail/cmd/service"
	"net/http"
)

type Mail interface {
	SendMail(w http.ResponseWriter, r *http.Request)
}

type mailController struct {
	json        *helpers.JsonResponse
	mailService service.MailService
}

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewMailController() Mail {
	return &mailController{
		json:        &helpers.JsonResponse{},
		mailService: service.NewMailService(),
	}
}

func (m *mailController) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := m.json.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = m.json.WriteJSONError(w, err)
		return
	}

	msg := entities.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = m.mailService.SendSMTPMessage(msg)
	if err != nil {
		_ = m.json.WriteJSONError(w, err)
		return
	}

	payload := helpers.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	_ = payload.WriteJSON(w, http.StatusAccepted, nil)
}
