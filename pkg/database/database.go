package database

import (
	"OFACDataUpdater/pkg/model"
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db      *sql.DB
	dbMutex sync.Mutex
)

const (
	DbUser     = "postgres-user"
	DbPassword = "postgres-password"
	DbName     = "postgres-db"
	DbHost     = "localhost" // Имя сервиса PostgreSQL в Docker Compose
	DbPort     = "5432"      // Порт PostgreSQL
)

// InitDB инициализирует базу данных.
func InitDB() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DbUser, DbPassword, DbName, DbHost, DbPort)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	return db.Ping()
}

// ClearData очищает данные в базе.
func ClearData() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec("DELETE FROM people")
	return err
}

// ImportDataFromOFAC импортирует данные из SDNList в базу данных.
func ImportDataFromOFAC(sdnList model.SDNList) error {
	// Извлечение данных из sdnList и вставка в базу данных.
	for _, sdnEntry := range sdnList.SDNs {
		if sdnEntry.SDNType == "Individual" {
			// Ваша логика выбора данных из sdnEntry
			uid := sdnEntry.UID
			firstName := sdnEntry.FirstName
			lastName := sdnEntry.LastName

			// Преобразование uid из строки в целое число (int)
			uidInt, err := strconv.Atoi(uid)
			if err != nil {
				// Обработка ошибки, если преобразование не удалось
				return err
			}

			// Вставка данных в базу данных, используя InsertPerson или подобные функции.
			err = InsertPerson(uidInt, firstName, lastName)
			if err != nil {
				// Обработка ошибки в случае неудачной вставки данных
				return err
			}
		}
	}

	return nil
}

// filterIndividuals фильтрует записи по sdnType=Individual.
func filterIndividuals(sdnList model.SDNList) ([]model.Person, error) {
	var individuals []model.Person
	for _, sdnEntry := range sdnList.SDNs {
		if sdnEntry.SDNType == "Individual" {
			uid, err := strconv.Atoi(sdnEntry.UID)
			if err != nil {
				return nil, err
			}

			individual := model.Person{
				UID:       uid,
				FirstName: sdnEntry.FirstName,
				LastName:  sdnEntry.LastName,
			}
			individuals = append(individuals, individual)
		}
	}
	return individuals, nil
}

// InsertPerson вставляет запись о человеке в базу.
func InsertPerson(uid int, firstName, lastName string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec(
		"INSERT INTO people (uid, first_name, last_name) VALUES ($1, $2, $3)",
		uid, firstName, lastName,
	)
	return err
}

// GetNamesStrong возвращает список имен с сильным совпадением.
func GetNamesStrong(name string) ([]model.Person, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	rows, err := db.Query("SELECT uid, first_name, last_name FROM people WHERE first_name || ' ' || last_name = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []model.Person
	for rows.Next() {
		var person model.Person
		err := rows.Scan(&person.UID, &person.FirstName, &person.LastName)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}

// GetNamesWeak возвращает список имен с слабым совпадением.
func GetNamesWeak(name string) ([]model.Person, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	rows, err := db.Query("SELECT uid, first_name, last_name FROM people WHERE first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []model.Person
	for rows.Next() {
		var person model.Person
		err := rows.Scan(&person.UID, &person.FirstName, &person.LastName)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}

// GetNames возвращает список имен.
func GetNames(name string) ([]model.Person, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	rows, err := db.Query("SELECT uid, first_name, last_name FROM people WHERE first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []model.Person
	for rows.Next() {
		var person model.Person
		err := rows.Scan(&person.UID, &person.FirstName, &person.LastName)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}

// UpdateData обновляет данные в базе данных из переданного списка.
func UpdateData(people []model.Person) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Откатываем транзакцию в случае ошибки
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Очищаем данные внутри транзакции
	if _, err := tx.Exec("DELETE FROM people"); err != nil {
		return err
	}

	// Вставляем новые данные внутри транзакции
	stmt, err := tx.Prepare("INSERT INTO people (uid, first_name, last_name) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, person := range people {
		_, err := stmt.Exec(person.UID, person.FirstName, person.LastName)
		if err != nil {
			return err
		}
	}

	return nil
}
