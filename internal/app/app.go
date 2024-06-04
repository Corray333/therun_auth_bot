package app

import (
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"regexp"

	"github.com/Corray333/authbot/internal/storage"
	"github.com/Corray333/authbot/internal/types"
	"github.com/Corray333/authbot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage interface {
	SaveCode(phone, code, type_id string) error
	SetPhone(chatID int64, query *types.CodeQuery) error
	GetPhone(chatID int64) (*types.CodeQuery, error)
}

type App struct {
	Storage Storage
}

func New() *App {
	return &App{Storage: storage.NewStorage()}
}

func (app *App) Run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // игнорировать все не-сообщения
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				// Получение аргументов командной строки
				args := update.Message.CommandArguments()

				if args == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}
				fmt.Println()
				fmt.Println(args)
				fmt.Println()
				// Парсинг URL для извлечения номера телефона
				parsedURL, err := url.Parse("http://dummy?" + args)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}
				params := parsedURL.Query()
				phone := params.Get("phone")
				typeId := params.Get("type_id")

				if phone == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}

				exp, err := regexp.Compile(`^\d+$`)
				if err != nil {
					slog.Error(err.Error())
				}
				test := exp.Find([]byte(phone))
				if len(test) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}
				app.Storage.SetPhone(update.Message.Chat.ID, &types.CodeQuery{Phone: "+" + phone, TypeID: typeId})

				// Отправка сообщения с кнопкой для отправки контакта
				button := tgbotapi.NewKeyboardButtonContact("Получить код")
				keyboard := tgbotapi.NewReplyKeyboard(
					[]tgbotapi.KeyboardButton{button},
				)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите кнопку ниже, чтобы получить код.")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
				continue

			}
		} else if update.Message.Contact != nil {
			userPhone := "+" + update.Message.Contact.PhoneNumber

			saved, err := app.Storage.GetPhone(update.Message.Chat.ID)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Время на запрос истекло. Повторите запрос в приложении.")
				bot.Send(msg)
				continue
			}

			fmt.Println()
			fmt.Println(userPhone, " ", saved.Phone)
			fmt.Println()

			if userPhone == saved.Phone {
				code := utils.GenerateVerificationCode()

				// Сохранение кода или выполнение другой логики
				if err := app.Storage.SaveCode(userPhone, code, saved.TypeID); err != nil {
					fmt.Println()
					fmt.Println(err)
					fmt.Println()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что-то пошло не так. Повторите запрос позже.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваш проверочный код: %s", code))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели неверный номер. Вернитесь в приложение и опробуйте снова.")
				bot.Send(msg)
			}
		}
	}
}
