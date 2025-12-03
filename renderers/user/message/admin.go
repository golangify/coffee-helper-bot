package message

import (
	"bytes"
	"coffee-helper/models"
	"html/template"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

type Renderer struct {
	templates *template.Template
}

func New(templates *template.Template) *Renderer {
	r := &Renderer{
		templates: templates,
	}
	return r
}

func (r *Renderer) Apanel(sentFrom *botmodels.User, user *models.User) *bot.SendMessageParams {
	msgTextBuffer := bytes.NewBuffer(nil)
	if err := r.templates.ExecuteTemplate(msgTextBuffer, "apanel", map[string]any{
		"sentfrom": sentFrom,
		"user":     user,
	}); err != nil {
		panic(err)
	}
	msg := &bot.SendMessageParams{
		ChatID:    sentFrom.ID,
		Text:      msgTextBuffer.String(),
		ParseMode: "html",
		ReplyMarkup: botmodels.InlineKeyboardMarkup{
			InlineKeyboard: [][]botmodels.InlineKeyboardButton{
				{{Text: "Список пользователей", CallbackData: "listUserCategories"}},
				{{Text: "Добавить администратора", CallbackData: "newAdminInvite"}},
				{{Text: "❌ (закрыть)", CallbackData: "deleteMessage"}},
			},
		},
	}

	return msg
}
