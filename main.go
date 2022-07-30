package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const jsonFile = "data.json"

type Trains []Train

func (t *Train) UnmarshalJSON(data []byte) error {
	type TrainCopy Train

	tmp := struct {
		ArrivalTime   string
		DepartureTime string
		*TrainCopy
	}{
		TrainCopy: (*TrainCopy)(t),
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	t.ArrivalTime, err = time.Parse("15:04:05", tmp.ArrivalTime)
	if err != nil {
		return err
	}

	t.DepartureTime, err = time.Parse("15:04:05", tmp.DepartureTime)
	if err != nil {
		return err
	}

	return nil
}

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func main() {
	//	... запит даних від користувача
	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	//	... обробка помилки
	//	... друк result
	_, err := FindTrains("1902", "1981", "price")
	if err != nil {
		log.Println(err)
		return
	}

	//fmt.Println(len(trains), trains[0], trains[1])
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	file, err := os.Open(jsonFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	trains := make(Trains, 0)
	err = json.Unmarshal(bytes, &trains)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println(len(trains))

	return nil, nil // маєте повернути правильні значення
}
