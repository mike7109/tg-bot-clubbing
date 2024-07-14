package messages

import tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const MsgHello = `Привет! Я LinkKeeper Bot. Я помогу тебе хранить и организовывать ссылки. 
Для начала работы ты можешь добавить новую ссылку просто отправив её, чтобы посмотреть список используй /list. 
Для получения списка всех команд используй /help.`

const MsgHelp = `Вот список доступных команд:

/start - Запустить бота и получить приветственное сообщение.
/help - Показать список доступных команд и их описание.
/add - Добавить новую ссылку. Формат: /add [ссылка] [описание] [категория] - все параметры необязательны, кроме [ссылка].
/list - Показать все сохраненные ссылки, если указать категорию, то покажет все ссылки в этой категории. Формат: /list [категория]
/search - Найти ссылки по ключевому слову в описании. Формат: /search [ключевое слово]
/delete - Удалить ссылку по идентификатору. Формат: /delete [идентификатор ссылки]
/delete_all - Удалить все ссылки.
/edit - Изменить описание и категорию ссылки. Формат: /edit [идентификатор ссылки] [новое описание] [новая категория]
/remind - Установить напоминание о необходимости посетить ссылку. Формат: /remind [идентификатор ссылки] [время]

Если у тебя возникли вопросы или проблемы, не стесняйся обращаться за помощью!`

const (
	MsgUnknownCommand    = "Неизвестная команда 🤔"
	MsgNoSavedPages      = "У вас нет сохраненных страниц 🙊"
	MsgSaved             = "Сохранено! 👌"
	MsgAlreadyExists     = "Мы обновили вашу страницу! 🔄"
	MsgInvalidAddCommand = "Неверный формат команды. Пример: /add https://example.com"
	MsgAddUrl            = "Отправьте ссылку, которую хотите сохранить 📎"
	MsgInvalidUrl        = "Неверный формат ссылки. Попробуйте еще раз 🤔"
	MsgDeletedAll        = "Все ваши страницы удалены! 🗑"
	ErrorHandler         = "Произошла ошибка 🤯"
	ErrNoUserName        = "Извините, но для использования бота необходимо установить имя пользователя в настройках Telegram"
	MsgDeleted           = "Страница удалена! 🗑"
)

func SendMessage(tgBot *tgApi.BotAPI, chatID int64, text string) error {
	msgConfig := tgApi.NewMessage(chatID, text)
	_, err := tgBot.Send(msgConfig)
	return err
}

func SendMessageDisableWebPagePreview(tgBot *tgApi.BotAPI, chatID int64, text string) error {
	msgConfig := tgApi.NewMessage(chatID, text)
	msgConfig.DisableWebPagePreview = true
	_, err := tgBot.Send(msgConfig)
	return err
}

func SendInvalidUrlMessage(tgBot *tgApi.BotAPI, chatID int64) error {
	return SendMessage(tgBot, chatID, MsgInvalidUrl)
}

func SendNoSavedPagesMessage(tgBot *tgApi.BotAPI, chatID int64) error {
	return SendMessage(tgBot, chatID, MsgNoSavedPages)
}

func SendErrorHandler(tgBot *tgApi.BotAPI, chatID int64) error {
	return SendMessage(tgBot, chatID, ErrorHandler)
}

func SendErrNoUserName(tgBot *tgApi.BotAPI, chatID int64) error {
	return SendMessage(tgBot, chatID, ErrNoUserName)
}
