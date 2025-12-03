package workers

import (
	"coffee-helper/services"
	"coffee-helper/workers/mailer"
	"github.com/go-telegram/bot"
)

type Workers struct {
	Mailer *mailer.Worker
}

func New(bot *bot.Bot, services *services.Services) *Workers {
	w := &Workers{
		Mailer: mailer.New(bot, services),
	}

	return w
}
