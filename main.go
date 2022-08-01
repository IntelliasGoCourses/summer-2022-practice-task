package main

// Данний код бере данні про поїзди з date.json, зберігає їх у мапу і перетворює на слайс структур.
//Запитує та зчитує данні які вводить користувач (departureStation, arrivalStation, criteria),
//знаходить підходящі варіанти, сортує їх за обраним варіантом та видає 3 найліпші випадки.

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
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

func main() {

	departureStation, arrivalStation, criteria, err := ScanInput() //... запит даних від користувача
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := FindTrains(departureStation, arrivalStation, criteria) //Виклик функції з основною логікою
	if err != nil {
		fmt.Println(err)
	}

	for _, i := range result {
		fmt.Println(i) // друк result
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	var (
		result     Trains
		trainSlice Trains
	)

	file, err := ioutil.ReadFile("./data.json") //Відкриття файлу data.json
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(file, &trainSlice) //Розшифрування типу .json
	if jsonErr != nil {
		return nil, err
	}

	switch criteria { //Сортування за обраним варіантом
	case "price":
		sort.Slice(trainSlice, func(i, j int) (less bool) {
			if trainSlice[i].Price == trainSlice[j].Price {
				return trainSlice[i].TrainID < trainSlice[j].TrainID
			}
			return trainSlice[i].Price < trainSlice[j].Price
		})

	case "arrival-time":
		sort.Slice(trainSlice, func(i, j int) bool {
			if trainSlice[i].ArrivalTime.Equal(trainSlice[j].ArrivalTime) {
				return trainSlice[i].TrainID < trainSlice[j].TrainID
			}
			return trainSlice[i].ArrivalTime.Before(trainSlice[j].ArrivalTime)
		})

	case "departure-time":
		sort.Slice(trainSlice, func(i, j int) bool {
			if trainSlice[i].DepartureTime.Equal(trainSlice[j].DepartureTime) {
				return trainSlice[i].TrainID < trainSlice[j].TrainID
			}
			return trainSlice[i].DepartureTime.Before(trainSlice[j].DepartureTime)
		})
	}

	sizeOutputArr := 3

	for _, i := range trainSlice {
		departureStationNumber, _ := strconv.Atoi(departureStation)
		arrivalStationNumber, _ := strconv.Atoi(arrivalStation)
		if i.DepartureStationID == departureStationNumber && i.ArrivalStationID == arrivalStationNumber && sizeOutputArr > 0 {
			result = append(result, i)
			sizeOutputArr--
		}
	}
	return result, nil // маєте повернути правильні значення
}

//Ф-ція запису введенних данних. Якщо в поле станції введено літери або поле пусте, повертається помилка.
//Якщо спосіб сортування не корректно введений - помилка

func ScanInput() (departureStation, arrivalStation, criteria string, err error) {

	var in = bufio.NewReader(os.Stdin)

	fmt.Println("Станція з якої ви відправляєтесь?")
	departureStation, _ = in.ReadString('\n')
	departureStation = strings.TrimSpace(departureStation)

	if departureStation == "" {
		return "", "", "", errors.New("empty departure station")
	}

	intValue, err := strconv.Atoi(departureStation)
	if err != nil || intValue < 0 {
		return "", "", "", errors.New("bad departure station input")
	}

	fmt.Println("Станція куди ви прямуєте?")

	arrivalStation, _ = in.ReadString('\n')
	arrivalStation = strings.TrimSpace(arrivalStation)

	if arrivalStation == "" {
		return "", "", "", errors.New("empty arrival station")
	}

	intValue, err = strconv.Atoi(arrivalStation)
	if err != nil || intValue < 0 {
		return "", "", "", errors.New("bad arrival station input")
	}

	fmt.Println("Критерій сортування:  price, arrival-time, departure-time")

	criteria, _ = in.ReadString('\n')
	criteria = strings.TrimSpace(criteria)

	if criteria != "price" && criteria != "arrival-time" && criteria != "departure-time" {
		return "", "", "", errors.New("unsupported criteria")
	}

	return departureStation, arrivalStation, criteria, nil
}

// Данний метод допомагає ф-ції Unmarshal коректно прочитати данні

func (t *Train) UnmarshalJSON(j []byte) error {

	var (
		rawStrings map[string]interface{}
	)

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {

		if k == "trainId" {
			vf := v.(float64)
			t.TrainID = int(vf)
		}

		if k == "departureStationId" {
			vf := v.(float64)
			t.DepartureStationID = int(vf)
		}

		if k == "arrivalStationId" {
			vf := v.(float64)
			t.ArrivalStationID = int(vf)
		}

		if k == "price" {
			vf := v.(float64)
			t.Price = float32(vf)
		}

		if k == "arrivalTime" {
			ti, err := time.Parse("15:04:05", v.(string))
			if err != nil {
				return err
			}
			t.ArrivalTime = ti
		}

		if k == "departureTime" {
			ti, err := time.Parse("15:04:05", v.(string))
			if err != nil {
				return err
			}
			t.DepartureTime = ti
		}
	}
	return nil
}
