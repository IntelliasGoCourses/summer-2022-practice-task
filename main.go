package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	jsonFile              = "data.json"
	criteriaPrice         = "price"
	criteriaArrivalTime   = "arrival-time"
	criteriaDepartureTime = "departure-time"
	lowestNaturalNumber   = 1
)

var (
	UnsupportedCriteria      = errors.New("unsupported criteria")
	EmptyDepartureStation    = errors.New("empty departure station")
	EmptyArrivalStation      = errors.New("empty arrival station")
	BadArrivalStationInput   = errors.New("bad arrival station input")
	BadDepartureStationInput = errors.New("bad departure station input")
)

type Trains []Train

func (t Trains) filterByParams(params Params) Trains {
	resTrains := make(Trains, 0)

	for _, v := range t {
		if v.DepartureStationID == params.departureStation && v.ArrivalStationID == params.arrivalStation {
			resTrains = append(resTrains, v)
		}
	}

	resTrains.sortByCriteria(params.criteria)

	return resTrains
}

func (t Trains) sortByCriteria(criteria string) {
	switch criteria {
	case criteriaPrice:
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].Price < t[j].Price
		})
	case criteriaArrivalTime:
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].ArrivalTime.Before(t[j].ArrivalTime)
		})
	case criteriaDepartureTime:
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].DepartureTime.Before(t[j].DepartureTime)
		})
	}
}

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

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

type Params struct {
	departureStation int
	arrivalStation   int
	criteria         string
}

func main() {
	departureStation, arrivalStation, criteria := scanInput()
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println(err)
		return
	}

	if result == nil {
		fmt.Println("таких потягів немає")
		return
	}

	for _, v := range result {
		fmt.Printf("%+v\n", v)
	}
}

func FindTrains(depStation, arrStation, criteria string) (Trains, error) {
	params, err := checkParseParams(depStation, arrStation, criteria)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	trains := make(Trains, 0)
	err = json.Unmarshal(bytes, &trains)
	if err != nil {
		return nil, err
	}

	filteredTrains := trains.filterByParams(params)

	if len(filteredTrains) == 0 {
		return nil, nil
	}

	if len(filteredTrains) > 3 {
		filteredTrains = filteredTrains[:3]
	}

	return filteredTrains, nil
}

func scanInput() (departureStation, arrivalStation, criteria string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("введіть станцію відправлення: (натисніть Enter після вводу)")
	scanner.Scan()
	departureStation = scanner.Text()

	fmt.Println("введіть станцію прибуття: (натисніть Enter після вводу)")
	scanner.Scan()
	arrivalStation = scanner.Text()

	fmt.Println("введіть критерій, по якому сортувати результат: (натисніть Enter після вводу)")
	scanner.Scan()
	criteria = scanner.Text()

	return departureStation, arrivalStation, criteria
}

func checkParseParams(departureStation, arrivalStation, criteria string) (Params, error) {
	if departureStation == "" {
		return Params{}, EmptyDepartureStation
	}
	departureStationInt, err := strconv.Atoi(departureStation)
	if err != nil || departureStationInt < lowestNaturalNumber {
		return Params{}, BadDepartureStationInput
	}

	if arrivalStation == "" {
		return Params{}, EmptyArrivalStation
	}
	arrivalStationInt, err := strconv.Atoi(arrivalStation)
	if err != nil || arrivalStationInt < lowestNaturalNumber {
		return Params{}, BadArrivalStationInput
	}

	if criteria != criteriaPrice && criteria != criteriaArrivalTime && criteria != criteriaDepartureTime {
		return Params{}, UnsupportedCriteria
	}

	return Params{
		departureStation: departureStationInt,
		arrivalStation:   arrivalStationInt,
		criteria:         criteria,
	}, nil
}
