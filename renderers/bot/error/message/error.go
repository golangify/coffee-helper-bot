package message

import (
	"fmt"
	"html/template"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

// bot/error/message renderer
type Renderer struct {
	templates *template.Template
}

func New(templates *template.Template) *Renderer {
	r := &Renderer{
		templates: templates,
	}
	return r
}

func (r *Renderer) Error(sentFrom *botmodels.User, err any) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: sentFrom.ID,
		Text:   fmt.Sprintf("ошибка: %v", err),
	}
}
