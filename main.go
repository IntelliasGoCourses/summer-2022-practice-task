package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	_ "time"
)

var (
	unsuportedCriteria       = errors.New("unsupported criteria")
	emptyDepartureStation    = errors.New("empty departure station")
	emptyArrivalStation      = errors.New("empty arrival station")
	badArrivalStationInput   = errors.New("bad arrival station input")
	badDepartureStationInput = errors.New("bad departure station input")
	emptyValue               = ""
	myType                   map[string]interface{}
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func main() {
	// ... запит даних від користувача
	departureStation, arrivalStation, criteria, err := InputValue()
	if err != nil {
		log.Fatal(err)
	}

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	// ... обробка помилки
	if err != nil {
		log.Fatal(err)
	}
	// ... друк result
	for _, i := range result {
		fmt.Printf("%+v", i)
		fmt.Println("-----------------------Your Data-----------------------")
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	var trains Trains
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &trains)
	if departureStation == emptyValue {
		return nil, emptyDepartureStation
	}
	departureStationInt, err := strconv.Atoi(departureStation)
	if len(departureStation) <= 4 {
		return nil, badDepartureStationInput
	}
	if arrivalStation == emptyValue {
		return nil, emptyArrivalStation
	}
	arrivalStationInt, err := strconv.Atoi(arrivalStation)
	if len(arrivalStation) <= 4 {
		return nil, badArrivalStationInput
	}

	for _, t := range trains {
		if t.DepartureStationID == departureStationInt || t.ArrivalStationID == arrivalStationInt {
			trains = append(trains, t)
		}
	}
	if criteria == emptyValue {
		return nil, unsuportedCriteria
	}
	switch criteria {
	case "price":
		sort.Slice(trains, func(i, j int) bool {
			return trains[i].Price < trains[j].Price
		})
	case "arrival-time":
		sort.Slice(trains, func(i, j int) bool {
			return trains[j].ArrivalTime.Before(trains[i].ArrivalTime)
		})
	case "departure-time":
		sort.Slice(trains, func(i, j int) bool {
			return trains[i].DepartureTime.Before(trains[j].DepartureTime)
		})
	default:
		return nil, unsuportedCriteria
	}
	if len(trains) > 3 {
		return trains[:3], nil
	}

	if len(trains) == 0 {
		return nil, nil
	}

	return trains, nil
}
func InputValue() (departureStation, arrivalStation, criteria string, err error) {
	fmt.Println("Greetings! Enter your data!")
	var read = bufio.NewReader(os.Stdin)

	fmt.Println("Write the station from which you depart?(for example 1902)")
	departureStation, _ = read.ReadString('\n')
	departureStation = strings.TrimSpace(departureStation)

	fmt.Println("Write Station where are you going?(for example 1929)")

	arrivalStation, _ = read.ReadString('\n')
	arrivalStation = strings.TrimSpace(arrivalStation)

	fmt.Println("Sorting criteria: price, arrival-time, departure-time")

	criteria, _ = read.ReadString('\n')
	criteria = strings.TrimSpace(criteria)

	return
}

func (t *Train) UnmarshalJSON(jso []byte) error {

	err := json.Unmarshal(jso, &myType)
	if err != nil {
		return err
	}

	for k, v := range myType {

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
			clock, err := time.Parse("15:04:05", v.(string))
			if err != nil {
				return err
			}
			t.ArrivalTime = clock
		}

		if k == "departureTime" {
			clock, err := time.Parse("15:04:05", v.(string))
			if err != nil {
				return err
			}
			t.DepartureTime = clock
		}
	}
	return nil
}
