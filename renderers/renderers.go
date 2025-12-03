package renderers

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"coffee-helper/models"
	"coffee-helper/renderers/bot"
	"coffee-helper/renderers/user"
)

type Renderers struct {
	Bot     *bot.Renderer
	User    *user.Renderer
	Menu    any
	Product any
	Search  any
}

func New(templateFS embed.FS) (*Renderers, error) {
	templates, err := loadTemplates(templateFS)
	if err != nil {
		return nil, err
	}

	r := &Renderers{
		Bot:  bot.New(templates),
		User: user.New(templates),
	}

	return r, nil
}

func loadTemplates(templateFS embed.FS) (*template.Template, error) {
	templatesSubFS, err := fs.Sub(templateFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("getting sub FS: %w", err)
	}

	// Создаем и парсим в одном вызове
	tmpl := template.Must(
		template.New("coffee-helper").Funcs(template.FuncMap{
			"join": strings.Join,
			"userFlagTitle": func(f string) string {
				return models.UserFlagTitle[f]
			},
		}).ParseFS(templatesSubFS, "**/*.html"),
	)

	return tmpl, nil
}
