package controllers

import (
	"coffee-helper/config"
	"coffee-helper/controllers/callbackquery"
	"coffee-helper/controllers/command"
	"coffee-helper/controllers/middleware"
	"coffee-helper/controllers/stephandler"
	"coffee-helper/renderers"
	"coffee-helper/services"
	"coffee-helper/workers"

	"github.com/go-telegram/bot"
)

type Controllers struct{}

func New(config *config.Config, b *bot.Bot, services *services.Services, workers *workers.Workers, middleware *middleware.Controller, stephandler *stephandler.StepHandler, renderers *renderers.Renderers) (*Controllers, error) {
	c := &Controllers{}

	_ = command.New(config, b, services, workers, middleware, stephandler, renderers)
	_ = callbackquery.New(b, services, middleware, renderers)

	return c, nil
}
