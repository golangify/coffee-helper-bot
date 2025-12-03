package error

import (
	"html/template"

	errormessage "coffee-helper/renderers/bot/error/message"
)

// bot/error renderer
type Renderer struct {
	Message *errormessage.Renderer
}

func New(templates *template.Template) *Renderer {
	r := &Renderer{
		Message: errormessage.New(templates),
	}

	return r
}
