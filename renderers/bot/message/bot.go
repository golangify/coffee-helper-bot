package message

import (
	"bytes"
	"html/template"

	"coffee-helper/models"
	"coffee-helper/services/user/role"

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

func (r *Renderer) Start(sentFrom *botmodels.User, user *models.User) *bot.SendMessageParams {
	msgTextBuffer := bytes.NewBuffer(nil)
	if err := r.templates.ExecuteTemplate(msgTextBuffer, "start", map[string]any{
		"sentfrom": sentFrom,
		"user":     user,
	}); err != nil {
		panic(err)
	}
	true_ := true
	msg := &bot.SendMessageParams{
		ChatID:    sentFrom.ID,
		Text:      msgTextBuffer.String(),
		ParseMode: "html",
		LinkPreviewOptions: &botmodels.LinkPreviewOptions{
			IsDisabled: &true_,
		},
		ReplyMarkup: botmodels.ReplyKeyboardMarkup{
			Keyboard: [][]botmodels.KeyboardButton{
				{{Text: "Меню"}},
				{{Text: "Действия"}},
			},
			ResizeKeyboard: true,
		},
	}

	return msg
}

func (r *Renderer) RoleIssue(roleIssue *role.RoleIssue) {
	msgTextBuffer := bytes.NewBuffer(nil)
	if err := r.templates.ExecuteTemplate(msgTextBuffer, "roleissue.html", map[string]any{
		"roleissue": roleIssue,
	}); err != nil {
		panic(err)
	}
}
