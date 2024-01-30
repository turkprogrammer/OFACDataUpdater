package handler

import (
	"OFACDataUpdater/pkg/database"
	"OFACDataUpdater/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	updateNow = make(chan struct{})
	ready     = make(chan struct{})
	update    []string
	updateMu  sync.Mutex
)

// InitHandlers инициализирует обработчики HTTP-запросов.
func InitHandlers() error {
	router := InitRouter() // импорт InitRouter
	http.Handle("/", router)
	return nil
}

// updateHandler обрабатывает запрос на обновление данных.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received update request")

	// Попытка отправить сигнал на обновление данных
	select {
	case updateNow <- struct{}{}:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"result": true, "info": "", "code": 200}`)
	default:
		// Если updateNow занят, возвращаем ошибку сервиса недоступен
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"result": false, "info": "Сервис недоступен", "code": 503}`)
	}

	log.Println("Update request processed")
}

// stateHandler обрабатывает запрос на получение текущего состояния данных.
func stateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем статус готовности данных
	select {
	case <-ready:
		w.Header().Set("Content-Type", "application/json")
		updateMu.Lock()
		defer updateMu.Unlock()
		// Если в процессе обновления, возвращаем информацию об этом
		if len(update) > 0 {
			fmt.Fprintf(w, `{"result": false, "info": "Обновление в процессе"}`)
		} else {
			// Если нет в процессе обновления, данные готовы к использованию
			fmt.Fprintf(w, `{"result": true, "info": "Ок"}`)
		}
	default:
		// Если нет данных, возвращаем информацию о пустом состоянии
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"result": false, "info": "Пусто"}`)
	}
}

// handleNamesResponse обрабатывает ответ на запрос списка имен человека.
func handleNamesResponse(w http.ResponseWriter, names []model.Person, err error) {
	if err != nil {
		// Если произошла ошибка, возвращаем ошибку клиенту
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"result": false, "info": "%s"}`, err.Error())
		return
	}

	// Если ошибок нет, возвращаем список имен в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

// getNamesStrong обрабатывает запрос на получение списка имен человека с сильным совпадением.
func getNamesStrong(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	names, err := database.GetNamesStrong(name)
	handleNamesResponse(w, names, err)
}

// getNamesWeak обрабатывает запрос на получение списка имен человека с слабым совпадением.
func getNamesWeak(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	names, err := database.GetNamesWeak(name)
	handleNamesResponse(w, names, err)
}

// getNames обрабатывает запрос на получение списка имен человека.
func getNames(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	names, err := database.GetNames(name)
	handleNamesResponse(w, names, err)
}
