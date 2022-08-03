package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
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

func (tr *Train) UnmarshalJSON(data []byte) error {
	var err error

	type CopyTrain Train
	temp := &struct {
		ArrivalTime   string `json:"arrivalTime"`
		DepartureTime string `json:"departureTime"`
		*CopyTrain
	}{
		CopyTrain: (*CopyTrain)(tr),
	}

	err = json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	tr.ArrivalTime, err = time.Parse("15:04:05", temp.ArrivalTime)
	if err != nil {
		return err
	}

	tr.DepartureTime, err = time.Parse("15:04:05", temp.DepartureTime)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	var result Trains
	var departureStation, arrivalStation, criteria string

	fmt.Println("Hi, traveler! Please input, space separated:")
	fmt.Println("* Departure Station Number")
	fmt.Println("* Arrival Station Number")
	fmt.Println("* How to sort trains? (price, arrival-time, departure-time")
	fmt.Println("For example: 1902 1929 departure-time")
	fmt.Scan(&departureStation, &arrivalStation, &criteria)

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(result) && i < 3; i++ {
		fmt.Println(result[i])
	}
}

func validateInput(depStation, arrStation, criteria string) error {
	if criteria != "price" && criteria != "arrival-time" && criteria != "departure-time" {
		return fmt.Errorf("unsupported criteria")
	}

	if depStation == "" {
		return fmt.Errorf("empty departure station")
	}

	if arrStation == "" {
		return fmt.Errorf("empty arrival station")
	}

	var err error

	_, err = strconv.Atoi(depStation)
	if err != nil {
		return fmt.Errorf("bad departure station input: %v", err)
	}

	_, err = strconv.Atoi(arrStation)
	if err != nil {
		return fmt.Errorf("bad arrival station input: %v", err)
	}

	return nil
}

func FindTrains(depStation, arrStation, criteria string) (Trains, error) {

	err := validateInput(depStation, arrStation, criteria)
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open("data.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open data.json: %v", err)
	}

	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var allTrains Trains

	err = json.Unmarshal(bytes, &allTrains)
	if err != nil {
		return nil, err
	}

	var result []Train

	depStInt, err := strconv.Atoi(depStation)
	if err != nil {
		return nil, err
	}

	arrStInt, err := strconv.Atoi(arrStation)
	if err != nil {
		return nil, err
	}

	for _, train := range allTrains {
		if train.DepartureStationID == depStInt && train.ArrivalStationID == arrStInt {
			result = append(result, train)
		}
	}

	switch criteria {
	case "price":
		sort.SliceStable(result, func(i, j int) bool { return result[i].Price < result[j].Price })
	case "arrival-time":
		sort.SliceStable(result, func(i, j int) bool { return result[i].ArrivalTime.Before(result[j].ArrivalTime) })
	case "departure-time":
		sort.SliceStable(result, func(i, j int) bool { return result[i].DepartureTime.Before(result[j].DepartureTime) })
	}

	if len(result) == 0 {
		result = nil
	} else if len(result) >= 3 {
		result = result[:3]
	}
	return result, nil
}
