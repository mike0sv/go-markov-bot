package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"
)

type RaverRequest struct {
	q string
}

func CreateStatsHandler(stats *Stats) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var t RaverRequest
		err = json.Unmarshal(body, &t)
		if err != nil {
			panic(err)
		}
		context := strings.Split(t.q, " ") // TODO context preprocessing
		data := map[string]string{"a": GenerateOne(stats, context)}
		payload, _ := json.Marshal(data)
		fmt.Fprintf(w, string(payload)) // отправляем данные на клиентскую сторону
	}

}

func RunServer(stats *Stats) error {
	http.HandleFunc("/", CreateStatsHandler(stats)) // установим роутер
	err := http.ListenAndServe("0.0.0.0:9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}
	return nil
}
