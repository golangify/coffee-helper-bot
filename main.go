package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"

	"coffee-helper/config"
	"coffee-helper/controllers"
	"coffee-helper/controllers/middleware"
	"coffee-helper/renderers"
	"coffee-helper/repositories/gorm"
	"coffee-helper/services"
	"coffee-helper/workers"

	"github.com/go-telegram/bot"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// init config
	cfg, err := config.LoadJSON("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	// init repositories
	repos, err := gorm.NewSQLite(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// init services
	servs := services.New(cfg, repos)

	// init renderers
	rendrs, err := renderers.New(templatesFS)
	if err != nil {
		log.Fatalln(err)
	}

	// init middleware
	midware := middleware.New(cfg, servs, rendrs)

	// init bot
	bot, err := bot.New(cfg.TelegramAPIToken, bot.WithMiddlewares(midware.Common))
	if err != nil {
		log.Fatalln(err)
	}

	// init workers
	wrkrs := workers.New(bot, servs)

	// init controllers
	if _, err := controllers.New(cfg, bot, servs, wrkrs, midware, rendrs); err != nil {
		log.Fatalln(err)
	}

	if bot, err := bot.GetMe(ctx); err == nil {
		log.Printf("authorized @%s", bot.Username)
	}

	bot.Start(ctx)
}
