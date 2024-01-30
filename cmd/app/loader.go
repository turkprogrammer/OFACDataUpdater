package main

import (
	"OFACDataUpdater/pkg/model"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// loadSDNList загружает данные SDNList из https://www.treasury.gov/ofac/downloads/sdn.xml.
func loadSDNList() (model.SDNList, error) {
	// URL для загрузки данных
	url := "https://www.treasury.gov/ofac/downloads/sdn.xml"

	// Выполнение HTTP-запроса для получения данных
	resp, err := http.Get(url)
	if err != nil {
		return model.SDNList{}, fmt.Errorf("ошибка при выполнении HTTP-запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.SDNList{}, fmt.Errorf("ошибка при чтении тела ответа: %v", err)
	}

	// Разбор XML-данных в структуру SDNList
	var sdnList model.SDNList
	err = xml.Unmarshal(body, &sdnList)
	if err != nil {
		return model.SDNList{}, fmt.Errorf("ошибка при разборе XML: %v", err)
	}

	return sdnList, nil
}
