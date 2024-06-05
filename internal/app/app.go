package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"

	"github.com/Corray333/authbot/internal/storage"
	"github.com/Corray333/authbot/internal/types"
	"github.com/Corray333/authbot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage interface {
	SaveCode(phone, code string, type_id int) error
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
				args_b64 := update.Message.CommandArguments()

				if args_b64 == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}

				decodedArgs, err := base64.StdEncoding.DecodeString(args_b64)
				if err != nil {
					fmt.Println("Ошибка декодирования:", err)
					return
				}
				var query types.CodeQuery
				if err := json.Unmarshal(decodedArgs, &query); err != nil {
					fmt.Println("Ошибка декодирования:", err)
					return
				}

				if query.Phone == "" || query.TypeID == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}

				exp, err := regexp.Compile(`^\+\d{1,2}\d{9}$`)
				if err != nil {
					slog.Error(err.Error())
				}
				test := exp.Find([]byte(query.Phone))
				fmt.Println(test)
				if len(test) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения кода, войдите в бота через приложение TheRun.")
					bot.Send(msg)
					continue
				}
				if err := app.Storage.SetPhone(update.Message.Chat.ID, &query); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что-то пошло не так. Повторите запрос позже.")
					bot.Send(msg)
					continue
				}

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
			userPhone := update.Message.Contact.PhoneNumber

			saved, err := app.Storage.GetPhone(update.Message.Chat.ID)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Время на запрос истекло. Повторите запрос в приложении.")
				bot.Send(msg)
				continue
			}

			if userPhone == saved.Phone {
				code := utils.GenerateVerificationCode()

				// Сохранение кода или выполнение другой логики
				if err := app.Storage.SaveCode(userPhone, code, saved.TypeID); err != nil {

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
