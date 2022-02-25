package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Точка входа в программу
func main() {
	botToken := os.Getenv("Bot_Token")
	//https://api.telegram.org/bot<token>/METHOD_NAME
	botApi := os.Getenv("Bot_URL")
	botUrl := botApi + botToken
	offset := 0
	//var researchBD []ResearchType
	bd := loadBD()

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Что то пошло не так:", err.Error())
		}
		for _, update := range updates {
			err = respond(botUrl, update, bd)
			offset = update.UpdateId + 1

		}
		fmt.Println(updates)

	}
}

//Запрос обновлений

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

//Ответ на обновления

func respond(botUrl string, update Update, bd [][]string) error {
	var botMessage BotMessage
	var text string
	res := searchResult(bd, update.Message.Text)
	for _, el := range res {
		text += el
	}
	update.Message.Text = text
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

//парсинг CSV в структуру для обработки
func loadBD() [][]string {
	var matrix [][]string

	file, err := os.Open("book2.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = 5
	reader.Comma = ';'
	reader.Comment = '#'

	for {
		var row []string
		record, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}
		row = append(row, record[0], record[1], record[2], record[3], record[4])
		matrix = append(matrix, row)
	}
	return matrix
}

func searchResult(bd [][]string, str string) []string {
	var res []string
	for _, r := range bd {
		if strings.Contains(strings.ToLower(r[2]), strings.ToLower(str)) {
			text := "Название:   " + r[2] + "\n" + "\n" + "Время выполнения:   " + r[3] + " д" + "\n" + "\n" + "Цена:   " + r[4] + "  сом" + "\n" + "\n" + "------" + "------" + "\n"
			res = append(res, text)
		}
	}
	return res
}
