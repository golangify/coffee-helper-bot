package bot

import (
	"html/template"

	"coffee-helper/renderers/bot/error"
	"coffee-helper/renderers/bot/message"
)

// bot renderer
type Renderer struct {
	Error *error.Renderer

	Message *message.Renderer
}

func New(templates *template.Template) *Renderer {
	r := &Renderer{
		Error:   error.New(templates),
		Message: message.New(templates),
	}

	return r
}
