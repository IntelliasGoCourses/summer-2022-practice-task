package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const dataFile = "data.json"

var (
	unsupported       = errors.New("unsupported criteria")
	emptyDeparture    = errors.New("empty departure station")
	emptyArrival      = errors.New("empty arrival station")
	badInputArrival   = errors.New("bad arrival station input")
	badInputDeparture = errors.New("bad departure station input")
)

type Trains []Train

func printTrains(t Trains) {
	for _, v := range t {
		fmt.Printf("%+v\n", v)
	}
}

type Train struct {
	TrainID            int       `json:"trainID"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func (t *Train) UnmarshalJSON(data []byte) error {
	type Alias Train
	aux := struct {
		ArrivalTime   string `json:"arrivalTime"`
		DepartureTime string `json:"departureTime"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	arrivalTime, err := time.Parse("15:04:05", aux.ArrivalTime)
	if err != nil {
		return err
	}
	departureTime, err := time.Parse("15:04:05", aux.DepartureTime)
	if err != nil {
		return err
	}
	t.ArrivalTime = arrivalTime
	t.DepartureTime = departureTime
	return nil
}

func FindTrains(depStation, arrStation, criteria string) (Trains, error) {

	if depStation == "" {
		return nil, emptyDeparture
	}
	depStationInt, err := strconv.Atoi(depStation)
	if err != nil || depStationInt < 1 {
		return nil, badInputDeparture
	}
	if arrStation == "" {
		return nil, emptyArrival
	}
	arrStationInt, err := strconv.Atoi(arrStation)
	if err != nil || arrStationInt < 1 {
		return nil, badInputArrival
	}

	var allTrains, requiredTrains Trains

	byteSlice, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, fmt.Errorf("file reading error:%w", err)
	}

	err = json.Unmarshal(byteSlice, &allTrains)
	if err != nil {
		return nil, err
	}

	for _, train := range allTrains {
		if train.DepartureStationID == depStationInt && train.ArrivalStationID == arrStationInt {
			requiredTrains = append(requiredTrains, train)
		}
	}

	switch strings.ToLower(criteria) {
	case "price":
		sort.SliceStable(requiredTrains, func(i, j int) bool {
			return requiredTrains[i].Price < requiredTrains[j].Price
		})
	case "arrival-time":
		sort.SliceStable(requiredTrains, func(i, j int) bool {
			return requiredTrains[i].ArrivalTime.Before(requiredTrains[j].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(requiredTrains, func(i, j int) bool {
			return requiredTrains[i].ArrivalTime.Before(requiredTrains[j].ArrivalTime)
		})
	default:
		return nil, unsupported
	}

	if len(requiredTrains) > 3 {
		requiredTrains = requiredTrains[:3]
	}

	return requiredTrains, nil
}

func inputParams() (depStation, arrStation, criteria string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter departure station: ")
	depStation, _ = reader.ReadString('\n')
	depStation = strings.TrimSpace(depStation)

	fmt.Print("Enter arrival station: ")
	arrStation, _ = reader.ReadString('\n')
	arrStation = strings.TrimSpace(arrStation)

	fmt.Print("Enter sorting criteria: ")
	criteria, _ = reader.ReadString('\n')
	criteria = strings.TrimSpace(criteria)

	return depStation, arrStation, criteria
}

func main() {
	depStation, arrStation, criteria := inputParams()

	result, err := FindTrains(depStation, arrStation, criteria)

	if err != nil {
		log.Fatal(err)
	}

	printTrains(result)
}
