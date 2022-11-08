package models

import (
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/utils/telegram_api"
	"log"
)

type notifierAction interface {
	execute(*mongo_models.Observer)
	setNext(notifierAction)
}

type TGMessenger struct {
	next notifierAction
}

func (n *TGMessenger) execute(o *mongo_models.Observer) {
	// send message
	err := telegram_api.SendMessageToNotifyAboutSignal(o)
	if err != nil {
		log.Printf("Could not send message to user with id: %v.Reason: %v\n", o.TelegramUserID, err)
	}
}

func (n *TGMessenger) setNext(next notifierAction) {
	n.next = next
}
