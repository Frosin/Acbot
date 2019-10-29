package TG

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type TG struct {
	bot *tgbotapi.BotAPI
}

func (tg *TG) Start(token string, webhookURL string) (err error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return
	}
	bot.Debug = true
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	tg.bot = bot
	return
}

func (tg *TG) sendTextMessage(id int64, text string) (err error) {
	message := tgbotapi.NewMessage(id, text)
	_, err = tg.bot.Send(message)
	return
}

func (tg *TG) sendButtonsMessage(id int64, text string, keyboard map[string]string, callbackQueryId string, answerText string) (err error) {
	message := tgbotapi.NewMessage(id, text)
	var buttons []tgbotapi.InlineKeyboardButton
	for i, v := range keyboard {
		strData := v
		button := tgbotapi.InlineKeyboardButton{Text: i, CallbackData: &strData}
		buttons = append(buttons, button)
	}

	callbackConf := tgbotapi.CallbackConfig{
		CallbackQueryID: callbackQueryId,
		Text:            answerText,
	}
	_, _ = tg.bot.AnswerCallbackQuery(callbackConf)

	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	_, err = tg.bot.Send(message)
	return
}

// callbackConf
