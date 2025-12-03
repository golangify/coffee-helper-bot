package user

import (
	"coffee-helper/renderers/user/message"
	"html/template"
)

type Renderer struct {
	Message *message.Renderer
}

func New(templates *template.Template) *Renderer {
	r := &Renderer{
		Message: message.New(templates),
	}

	return r
}
