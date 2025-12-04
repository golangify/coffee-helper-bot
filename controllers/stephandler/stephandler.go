package stephandler

import (
	"coffee-helper/controllers/middleware"
	"coffee-helper/models"
	"coffee-helper/renderers"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"sync"
)

const CancelCommand = "/cancel"

const (
	TypeStepAny uint8 = iota
	TypeStepText
	TypeStepImage
)

type StepFunc func(ctx context.Context, b *bot.Bot, update *botmodels.Update, user *models.User, args map[string]any)

type Step struct {
	User           *models.User
	Type           uint8
	NonCancellable bool
	Func           StepFunc
	Args           map[string]any
}

type StepHandler struct {
	middleware *middleware.Controller
	renderers  *renderers.Renderers

	// TODO: sync.Map ?
	handlers map[uint]*Step
	mu       sync.Mutex
}

func New(middleware *middleware.Controller, renderers *renderers.Renderers) *StepHandler {
	h := &StepHandler{
		middleware: middleware,
		renderers:  renderers,

		handlers: make(map[uint]*Step),
	}

	return h
}

func (h *StepHandler) Register(ctx context.Context, b *bot.Bot, user *models.User, htype uint8, nonCancellable bool, fn StepFunc, args map[string]any) (*Step, error) {
	s := &Step{
		User:           user,
		Type:           htype,
		NonCancellable: nonCancellable,
		Func:           fn,
		Args:           args,
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[user.ID] = s

	if !nonCancellable {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: user.TgID,
			Text:   "Чтобы отменить действие отправьте " + CancelCommand,
		}); err != nil {
			return s, err
		}
	}

	return s, nil
}

func (h *StepHandler) checkTypeMatch(update *botmodels.Update, s *Step) error {
	if s.Type == TypeStepAny {
		return nil
	}

	if s.Type == TypeStepText {
		if update.Message == nil || update.Message.Text == "" {
			return fmt.Errorf("ожидалось сообщение с текстом")
		}
		return nil
	}

	if s.Type == TypeStepImage {
		if update.Message == nil || len(update.Message.Photo) == 0 {
			return fmt.Errorf("ожидалось сообщение с изображением")
		}
		return nil
	}

	h.mu.Lock()
	delete(h.handlers, s.User.ID)
	h.mu.Unlock()
	return fmt.Errorf("неизвестный тип обработчика: %d", s.Type)
}

func (h *StepHandler) Middleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
			user := h.middleware.GetUser(ctx)

			h.mu.Lock()
			step, exists := h.handlers[user.ID]
			h.mu.Unlock()

			if !exists {
				next(ctx, b, update)
				return
			}

			if update.Message != nil &&
				update.Message.Text == CancelCommand {
				if step.NonCancellable {
					if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: user.TgID,
						Text:   "Действие нельзя отменить",
					}); err != nil {
						panic(err)
					}
					return
				}
				h.mu.Lock()
				delete(h.handlers, user.ID)
				h.mu.Unlock()
				if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: user.TgID,
					Text:   "Действие отменено",
				}); err != nil {
					panic(err)
				}
				return
			}

			if err := h.checkTypeMatch(update, step); err != nil {
				if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: user.TgID,
					Text:   fmt.Errorf("Ошибка в StepHandler.checkTypeMatch: %v", err).Error(),
				}); err != nil {
					panic(err)
				}
				return
			}

			h.mu.Lock()
			delete(h.handlers, user.ID)
			h.mu.Unlock()

			step.Func(ctx, b, update, user, step.Args)
		}
	}
}
