package main

import (
	"OFACDataUpdater/pkg/database"
	"OFACDataUpdater/pkg/handler"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Инициализация базы данных
	if err := database.InitDB(); err != nil {
		fmt.Println("Ошибка при инициализации базы данных:", err)
		return
	}

	// Загрузка данных SDNList перед запуском HTTP-сервера
	sdnList, err := loadSDNList()
	if err != nil {
		fmt.Println("Ошибка при загрузке данных SDNList:", err)
		return
	}

	// Импорт данных из OFAC перед запуском HTTP-сервера
	if err := database.ImportDataFromOFAC(sdnList); err != nil {
		fmt.Println("Ошибка при импорте данных из OFAC:", err)
		return
	}

	// Инициализация обработчиков
	if err := handler.InitHandlers(); err != nil {
		fmt.Println("Ошибка при инициализации обработчиков:", err)
		return
	}

	// Инициализация роутера
	router := handler.InitRouter()

	// Порт для HTTP-сервера
	port := 8082
	fmt.Printf("Сервер запущен на порту %d...\n", port)

	// Создаем объект http.Server с установленным роутером
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Канал для сигналов завершения работы приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для прослушивания сигналов завершения работы
	go func() {
		sig := <-stop
		fmt.Printf("Получен сигнал завершения: %v\n", sig)

		// Останавливаем сервер
		if err := server.Shutdown(nil); err != nil {
			fmt.Println("Ошибка при остановке сервера:", err)
		}
	}()

	// Запуск сервера
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Ошибка при запуске сервера:", err)
	}

	// После завершения работы сервера, приложение выходит из горутины
	fmt.Println("Приложение завершено.")
}
