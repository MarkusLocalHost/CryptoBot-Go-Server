package telegram_api

import (
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/utils/scrapers"
	"fmt"
	"gitlab.com/toby3d/telegraph"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func SendMessageToNotifyAboutSignal(observer *mongo_models.Observer) error {
	currentPrice := scrapers.GetPriceInCurrencyScraper(observer.CryptoID, observer.CurrencyOfValue)

	data := url.Values{
		"chat_id": {strconv.FormatInt(observer.TelegramUserID, 10)},
		"text": {fmt.Sprintf("Сигнал от обсервера на цену %f %s.Текущая цена: %f",
			observer.ExpectedValue,
			strings.ToUpper(observer.CurrencyOfValue),
			currentPrice,
		)}}

	urlToPost := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TG_API"))
	_, err := http.PostForm(urlToPost, data)
	if err != nil {
		return err
	}

	return nil
}

func SendMessageToNotifyAboutChangePriceInPercent(userId int64, currenciesData map[string]map[string]string) error {
	// collect text for telegraf page
	var textPageOne string
	var textPageTwo string
	var textPageThree string
	var textPageFour string
	var i int
	if len(currenciesData) != 0 {
		for _, currencyData := range currenciesData {
			i += 1
			if i <= 300 {
				textPageOne += makeTextForPage(currencyData)
			} else if i > 300 && i <= 600 {
				textPageTwo += makeTextForPage(currencyData)
			} else if i > 600 && i <= 900 {
				textPageThree += makeTextForPage(currencyData)
			} else {
				textPageFour += makeTextForPage(currencyData)
			}
		}
	} else {
		textPageOne += "Ничего нет"
	}

	// create telegraf page
	requisites := telegraph.Account{
		ShortName: "MarkusLocalHost", // required

		// Author name/link can be epmty. So secure. Much anonymously. Wow.
		AuthorName: "Nikita Khasanov",        // optional
		AuthorURL:  "https://t.me/markusnew", // optional
	}
	account, err := telegraph.CreateAccount(requisites)
	if err != nil {
		fmt.Println(err.Error())
	}

	content, err := telegraph.ContentFormat(textPageOne)
	if err != nil {
		fmt.Println(err.Error())
	}

	pageData := telegraph.Page{
		Title:   "My super-awesome page", // required
		Content: content,                 // required

		// Not necessarily, but, hey, it's just an example.
		AuthorName: account.AuthorName, // optional
		AuthorURL:  account.AuthorURL,  // optional
	}

	page, err := account.CreatePage(pageData, false)
	if err != nil {
		fmt.Println(err.Error())
	}

	data := url.Values{
		"chat_id": {strconv.FormatInt(userId, 10)},
		"text":    {page.URL},
	}

	urlToPost := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TG_API"))
	_, err = http.PostForm(urlToPost, data)
	if err != nil {
		return err
	}

	return nil
}

func SendMessageToNotifyAdminAboutTooBigDistinctValue(distinctValue int64) error {
	data := url.Values{
		"chat_id": {os.Getenv("ADMIN_TELEGRAM_ID")},
		"text": {fmt.Sprintf("Разница составила: %v",
			distinctValue,
		)}}

	urlToPost := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TG_API"))
	_, err := http.PostForm(urlToPost, data)
	if err != nil {
		return err
	}

	return nil
}

func makeTextForPage(currencyData map[string]string) string {
	return fmt.Sprintf(`Название: %s. Цена: %s -> %s. Изменение: %s%s
`,
		currencyData["currencyName"],
		currencyData["lastPrice"],
		currencyData["currentPrice"],
		currencyData["percentSignString"],
		currencyData["percent"],
	)
}

func SendMessageToUser(userId int64, msg string) error {
	data := url.Values{
		"chat_id": {fmt.Sprintf("%v", userId)},
		"text":    {msg},
	}

	urlToPost := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TG_API"))
	_, err := http.PostForm(urlToPost, data)
	if err != nil {
		return err
	}

	return nil
}
