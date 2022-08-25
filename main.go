package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

var (
	//помилки
	unsupportedCriteria      = errors.New("unsupported criteria")
	emptyDepartureStation    = errors.New("empty departure station")
	emptyArrivalStation      = errors.New("empty arrival station")
	badArrivalStationInput   = errors.New("bad arrival station input")
	badDepartureStationInput = errors.New("bad departure station input")
	//змінні для пошуку потягів. Int-й варіант введених користувачем даних
	neededDepartureId                          int
	neededArrivalId                            int
	departureStation, arrivalStation, criteria string
)

const (
	maxNaturalNumber   = 1          //константа потрібна для валідації введених користувачем значень
	resultLenghtNeeded = 3          //необхідна кількість потягів в результаті
	layout             = "15:04:05" //Щоб описати формат значення часу, потрібний спеціальний параметр макета layout.
	//Він повинен бути референтною датою/часом — Mon Jan 2 15:04:05 MST 2006, яка відформатована так як і очікувана дата/час
)

//Функція getContentFromJson зчитує файл і, у разі успіху, повертає слайс байтів з файлу
func getContentFromJson(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

//Функція ConvertStringToTime, у разі успіху, конвертує час з формату рядка в формат часу
func ConvertStringToTime(str string) (time.Time, error) {
	convertedTime, err := time.Parse(layout, str)
	if err != nil {
		panic(err)
	}
	return convertedTime, nil
}

//перевизначаємо метод UnmarshalJSON таким чином, щоб він працював з нашими кастомними типами time.Time
func (t *Train) UnmarshalJSON(content []byte) error {
	//створюємо структуру з полями аналогічними з структурою Train, але без кастомних типів
	type structWithNormalTypes struct {
		TrainID            int     `json:"trainId"`
		DepartureStationID int     `json:"departureStationId"`
		ArrivalStationID   int     `json:"arrivalStationId"`
		Price              float32 `json:"price"`
		ArrivalTime        string  `json:"arrivalTime"`
		DepartureTime      string  `json:"departureTime"`
	}
	//виконуємо Unmarshal в структуру без кастомних типів
	var tmp structWithNormalTypes
	err := json.Unmarshal(content, &tmp)
	if err != nil {
		return err
	}
	//присвоюємо поля структури бех кастомних типів структурі Train (з кастомними типами)
	//в момент присвоєння полю з типом time.Time застосовуємо до відповідного поля структури
	//без кастомних типів функцію ConvertStringToTime
	t.TrainID = tmp.TrainID
	t.DepartureStationID = tmp.DepartureStationID
	t.ArrivalStationID = tmp.ArrivalStationID
	t.Price = tmp.Price
	t.ArrivalTime, err = ConvertStringToTime(tmp.ArrivalTime)
	if err != nil {
		return err
	}
	t.DepartureTime, err = ConvertStringToTime(tmp.DepartureTime)
	if err != nil {
		return err
	}
	return nil
}

//Функція departureStationValidation перевіряє чи не було введено для станції відправлення пустого
//рядка, чи це ціле число і чи воно є натуральним. Якщо виконується хоча б одна з умов повертає помилку
func departureStationValidation(departureStation string) error {
	if len(departureStation) == 0 {
		return emptyDepartureStation
	}
	departureStationInt, err := strconv.Atoi(departureStation)
	if err != nil {
		return badDepartureStationInput
	}
	if departureStationInt < maxNaturalNumber {
		return badDepartureStationInput
	}
	neededDepartureId = departureStationInt
	return nil
}

//Функція arrivalStationValidation перевіряє чи не було введено для станції прибуття пустого
//рядка, чи це ціле число і чи воно є натуральним. Якщо виконується хоча б одна з умов повертає помилку
func arrivalStationValidation(arrivalStation string) error {
	if arrivalStation == "" {
		return emptyArrivalStation
	}
	arrivalStationInt, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return badArrivalStationInput
	}
	if arrivalStationInt <= maxNaturalNumber {
		return badArrivalStationInput
	}
	neededArrivalId = arrivalStationInt
	return nil
}

//Функція criteriaValidation перевіряє валідність введеного критерія.
//Створюємо мапу валідних критеріїв, де ключ це сам критерій. Значення по ключу використовуватися
//не буде, тому використовуємо тип struct{} (для його зберігання використовується 0 байт)
func criteriaValidation(criteria string) error {
	validCriteria := map[string]struct{}{
		"price":          {},
		"arrival-time":   {},
		"departure-time": {}}
	_, ok := validCriteria[criteria]
	if !ok {
		return unsupportedCriteria
	}
	return nil
}

func main() {
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	//	... обробка помилки
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	//	... друк result
	for _, train := range result {
		formatResult(train)
	}

}

//Функція FindTrains знаходить у файлі потяги, що задовольняють введені користувачем параметри і,
//у разі успіху повертає не більше 3 варіантів відсортованих по введеному критерію.
//Якщо ввід користувача валідний, але по даним параметрам немає потягів, то повертає nil,nil
func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	var err error
	//Отримуємо дані
	fmt.Println("Введіть номер станції відправлення.")
	fmt.Scanln(&departureStation)
	err = departureStationValidation(departureStation)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	fmt.Println("Введіть номер станції прибуття.")
	fmt.Scanln(&arrivalStation)
	err = arrivalStationValidation(arrivalStation)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	fmt.Println("Введіть критерій, по котрому треба відсортувати потяги в результаті. Валідні значення: price, arrival-time, departure-time")
	fmt.Scanln(&criteria)
	err = criteriaValidation(criteria)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	var t []Train
	//читаємо json
	content, err := getContentFromJson("data.json")
	if err != nil {
		fmt.Println("open file error: " + err.Error())
		return nil, nil
	}
	//виконуємо кастомний анмаршал
	err = json.Unmarshal(content, &t)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return nil, nil
	}

	//знаходимо всі потяги, що підходять по введеним користувачем станціям
	var foundedTrains Trains
	for _, train := range t {
		if train.DepartureStationID == neededDepartureId && train.ArrivalStationID == neededArrivalId {
			foundedTrains = append(foundedTrains, train)
		}
	}
	foundedTrains.sortFoundedTrains(criteria)
	//В залежності від кількості знайдених потягів повертаємо і різну їх кількість
	currentLenght := len(foundedTrains)
	if currentLenght == 0 {
		fmt.Println("Нічого не знайдено")
		return nil, nil
	}
	if currentLenght < resultLenghtNeeded {
		return foundedTrains, nil
	}
	return foundedTrains[:resultLenghtNeeded], nil
}

//Метод sortFoundedTrains сортує потяги знайдені за введеними параметрами
//щоб зберегти початковий порядок потягів використовується сортування sort.SliceStable
func (t Trains) sortFoundedTrains(criteria string) {
	switch criteria {
	case "price":
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].Price < t[j].Price
		})
	case "arrival-time":
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].ArrivalTime.Before(t[j].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].DepartureTime.Before(t[j].DepartureTime)
		})
	}
}

//Функція formatResult форматує друк структури
func formatResult(t Train) {
	fmt.Printf(`Train ID: %v, Departure station ID: %v, Arrival station ID: %v, Price: %v, Departure time: %s, Arrival time: %s.%s`,
		t.TrainID, t.DepartureStationID, t.ArrivalStationID, t.Price, t.DepartureTime.Format(layout), t.ArrivalTime.Format(layout),
		"\n")
}
