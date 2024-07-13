package task

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/repositories"
	"github.com/mike7109/tg-bot-clubbing/internal/service/processor"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
	"github.com/mike7109/tg-bot-clubbing/pkg/utls"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func Start(ctx context.Context, tgBot *tgApi.BotAPI) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgHello)
		_, err := tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func Help(ctx context.Context, tgBot *tgApi.BotAPI) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgHelp)
		_, err := tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func Rnd(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		page, err := storage.PickRandom(ctx, msg.From.UserName)
		if err != nil && !errors.Is(err, entity.ErrNoSavedPages) {
			return err
		}

		if errors.Is(err, entity.ErrNoSavedPages) {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgNoSavedPages)
			_, err = tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, page.URL)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		if err = storage.Remove(context.Background(), page); err != nil {
			log.Println("Failed to remove page: ", err)
			return nil
		}

		return nil
	}
}

func Save(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		if !strings.Contains(msg.Text, "/add") {
			msg.Text = "/add " + msg.Text
		}

		re := regexp.MustCompile(`^/add\s+(\S+)(?:\s+(.+?))?(?:\s+(.+))?(?:\s+(.+?))?$`)

		matches := re.FindStringSubmatch(msg.Text)

		if len(matches) == 0 {
			if strings.Contains(msg.Text, "/add") {
				msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgAddUrl)
				_, err := tgBot.Send(msgConfig)
				if err != nil {
					return err
				}

				return nil
			}

			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgInvalidAddCommand)
			_, err := tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		urlTrim := strings.TrimSpace(matches[1])
		if !utls.IsURL(urlTrim) {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgInvalidUrl)
			_, err := tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		var description, title, category *string

		if len(matches) > 2 {
			category = &matches[2]
		}

		if len(matches) > 3 {
			title = &matches[3]
		}
		if len(matches) > 4 {
			description = &matches[4]
		}

		page := &entity.Page{
			UserName:    msg.From.UserName,
			URL:         urlTrim,
			Title:       title,
			Category:    category,
			Description: description,
		}

		sendMsg := messages.MsgSaved

		isExists, err := storage.IsExists(ctx, page)
		if err != nil {
			return err
		}
		if isExists {
			sendMsg = messages.MsgAlreadyExists
		}

		if err := storage.Save(ctx, page); err != nil {
			return err
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, sendMsg)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		go func() {
			if page.Category != nil {
				_, categoryAI, err := FetchPageInfo(page.URL)
				if err != nil {
					log.Println("Failed to fetch page info: ", err)
					return
				}

				out, err := ClassifyLink(categoryAI)
				if err != nil {
					log.Println("Failed to classify link: ", err)
					return
				}

				if out != nil {
					page.Category = out
					if err := storage.Save(ctx, page); err != nil {
						log.Println("Failed to save category: ", err)
						return
					}

					msgConfig := tgApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Добавил категорию: %s к ссылке %s", *out, page.URL))
					_, err = tgBot.Send(msgConfig)
				}
			}
		}()

		return nil
	}
}

func SaveSimple(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		page := &entity.Page{
			UserName: msg.From.UserName,
			URL:      msg.Text,
		}

		sendMsg := messages.MsgSaved

		isExists, err := storage.IsExists(ctx, page)
		if err != nil {
			return err
		}
		if isExists {
			sendMsg = messages.MsgAlreadyExists
		}

		if err := storage.Save(ctx, page); err != nil {
			return err
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, sendMsg)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func ListUrl(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		pages, err := storage.ListUrl(ctx, msg.From.UserName)
		if err != nil && !errors.Is(err, entity.ErrNoSavedPages) {
			return err
		}

		if errors.Is(err, entity.ErrNoSavedPages) {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgNoSavedPages)
			_, err = tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		var urlList string
		for i, page := range pages {
			urlList += fmt.Sprintf("%d. %s ", i+1, page.URL)
			if page.Category != nil {
				urlList += fmt.Sprintf("%s ", *page.Category)
			}
			if page.Title != nil {
				urlList += fmt.Sprintf("%s ", *page.Title)
			}
			if page.Description != nil {
				urlList += fmt.Sprintf("%s ", *page.Description)
			}
			urlList += "\n"
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, urlList)
		msgConfig.DisableWebPagePreview = true // Отключаем веб-превью
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func DeleteAll(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		err := storage.DeleteAll(ctx, msg.From.UserName)
		if err != nil {
			return err
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgDeletedAll)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

type ClassificationRequest struct {
	Inputs     string `json:"inputs"`
	Parameters struct {
		CandidateLabels []string `json:"candidate_labels"`
		MultiLabel      bool     `json:"multi_label"`
	} `json:"parameters"`
}

type ClassificationResponse struct {
	Labels []string  `json:"labels"`
	Scores []float64 `json:"scores"`
}

func ClassifyLink(description string) (*string, error) {
	candidateLabels := []string{"спорт", "еда", "технологии", "политика", "музыка", "путешествия"}
	requestBody, _ := json.Marshal(ClassificationRequest{
		Inputs: description,
		Parameters: struct {
			CandidateLabels []string `json:"candidate_labels"`
			MultiLabel      bool     `json:"multi_label"`
		}{
			CandidateLabels: candidateLabels,
			MultiLabel:      false,
		},
	})

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/facebook/bart-large-mnli", bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("HUGGING_FACE_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var classificationResp ClassificationResponse
	err = json.NewDecoder(resp.Body).Decode(&classificationResp)
	if err != nil {
		return nil, err
	}

	if len(classificationResp.Labels) > 0 && classificationResp.Scores[0] > 0.35 {
		return &classificationResp.Labels[0], nil
	}

	return nil, nil
}

// Page представляет модель для хранения данных о странице.
type Page struct {
	ID          string
	URL         string
	UserName    string
	Name        string
	Description string
	Category    string
	CreatedAt   string
}

// fetchPageInfo получает информацию с указанного URL.
func FetchPageInfo(url string) (string, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch the page, status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	title := doc.Find("title").Text()
	description, _ := doc.Find("meta[name='description']").Attr("content")

	return title, description, nil
}
