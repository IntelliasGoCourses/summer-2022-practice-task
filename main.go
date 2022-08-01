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
	"time"
)

const (
	criteriaArrTime = "arrival-time"
	criteriaDepTime = "departure-time"
	criteriaPrice   = "price"
	minValidNum     = 1
)

var stationsExist = map[int]bool{0: false}

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

// UnmarshalJSON implementing custom UnmarshalJSON, cuz we can't parse into origin Train struct because of different time field types
func (t *Train) UnmarshalJSON(b []byte) error {
	const layout = "15:04:05"
	type TempTrain struct {
		TrainID            int     `json:"trainId"`
		DepartureStationID int     `json:"departureStationId"`
		ArrivalStationID   int     `json:"arrivalStationId"`
		Price              float32 `json:"price"`
		ArrivalTime        string  `json:"arrivalTime"`
		DepartureTime      string  `json:"departureTime"`
	}
	var tmp TempTrain
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	convArrTime, err := time.Parse(layout, tmp.ArrivalTime)
	if err != nil {
		return err
	}
	convDepTime, err := time.Parse(layout, tmp.DepartureTime)
	t.ArrivalTime = time.Date(0, time.January, 1, convArrTime.Hour(), convArrTime.Minute(), convArrTime.Second(), 0, time.UTC)
	t.DepartureTime = time.Date(0, time.January, 1, convDepTime.Hour(), convDepTime.Minute(), convDepTime.Second(), 0, time.UTC)
	t.TrainID = tmp.TrainID
	t.Price = tmp.Price
	t.DepartureStationID = tmp.DepartureStationID
	t.ArrivalStationID = tmp.ArrivalStationID
	return nil
}

func main() {
	//	... запит даних від користувача
	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	//	... обробка помилки
	//	... друк result
	depst, arrst, cr, err := InputScan()
	if err != nil {
		return
	}
	result, err := FindTrains(depst, arrst, cr)
	if err != nil {
		fmt.Printf("%s", err)
	}
	if result == nil {
		return
	}
	fmt.Printf("%+v ", result)

}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	trains, err := ReadAndUnmarshal()
	if err != nil {
		return nil, err
	}
	//checking if stations number is valid
	err = IsStationExist(departureStation, arrivalStation)
	if err != nil {
		return nil, err
	}
	trains = SortingTrainSlice(trains, departureStation, arrivalStation)
	trains = SortingByCriteria(trains, criteria)

	return trains, nil // маєте повернути правильні значення
}

//ReadAndUnmarshal func is using to open, read file and getting proper []Train for manipulate later
func ReadAndUnmarshal() ([]Train, error) {
	bytes, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	trains := make(Trains, 0)
	err = json.Unmarshal(bytes, &trains)
	if err != nil {
		fmt.Println("error: ", err)
	}
	for _, v := range trains {
		stationsExist[v.DepartureStationID] = true
	}
	return trains, nil
}

//InputScan func is simple scanning func that takes var from stdin and returning them after InputCheck checks them.
func InputScan() (departureStation, arrivalStation, criteria string, err error) {
	s := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter departure station.")
	s.Scan()
	departureStation = s.Text()

	fmt.Println("Enter arrival station.")
	s.Scan()
	arrivalStation = s.Text()

	fmt.Println("Enter criteria for sorting.")
	s.Scan()
	criteria = s.Text()
	err = InputCheck(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Printf("%s", err)
		return "", "", "", err
	}
	return departureStation, arrivalStation, criteria, nil
}

//InputCheck func provides validation for inputs from InputScan func.
func InputCheck(depStation, arrStation, criteria string) error {
	var (
		unsupportedCriteria = errors.New("unsupported criteria")
		emptyDepStation     = errors.New("empty departure station")
		emptyArrStation     = errors.New("empty arrival station")
		badArrStInput       = errors.New("bad arrival station input")
		badDepStInput       = errors.New("bad departure station input")
	)
	if depStation == "" {
		return emptyDepStation
	}
	depStInt, err := strconv.Atoi(depStation)
	if err != nil || depStInt <= minValidNum {
		return badDepStInput
	}
	if arrStation == "" {
		return emptyArrStation
	}
	arrStInt, err := strconv.Atoi(arrStation)
	if err != nil || arrStInt <= minValidNum {
		return badArrStInput
	}
	criteriaStatus := 0
	switch criteria {
	case criteriaArrTime:
		criteriaStatus = 1
	case criteriaDepTime:
		criteriaStatus = 2
	case criteriaPrice:
		criteriaStatus = 3
	}
	if criteriaStatus == 0 {
		return unsupportedCriteria
	}
	return nil
}

//IsStationExist func checks if departure station and arrival station are existing in our data.json.
func IsStationExist(departureStation, arrivalStation string) error {
	departStCheckInt, _ := strconv.Atoi(departureStation)
	arrStCheckInt, _ := strconv.Atoi(arrivalStation)
	wrongDeparture := errors.New("there is no departure station like this " + "(station №" + departureStation + ")")
	wrongArrival := errors.New("there is no arrival station like this " + "(station №" + arrivalStation + ")")
	sameStationErr := errors.New("your destination point is your departure point")
	if stationsExist[departStCheckInt] == false {
		return wrongDeparture
	} else if stationsExist[arrStCheckInt] == false {
		return wrongArrival
	} else if departStCheckInt == arrStCheckInt {
		return sameStationErr
	}
	return nil
}

//SortingTrainSlice func gives us sorted and trimmed []Train with only dep. and arr. stations we are interested in.
func SortingTrainSlice(t Trains, departureStation, arrivalStation string) Trains {
	sort.Slice(t, func(i, j int) bool {
		return t[i].DepartureStationID < t[j].DepartureStationID
	})
	var tmp []Train
	depStInt, _ := strconv.Atoi(departureStation)
	arrStInt, _ := strconv.Atoi(arrivalStation)
	for i, _ := range t {
		if t[i].DepartureStationID == depStInt && t[i].ArrivalStationID == arrStInt {
			tmp = append(tmp, t[i])
		}
	}
	return tmp
}

func SortingByCriteria(t Trains, criteria string) Trains {
	if criteria == criteriaPrice {
		sort.Slice(t, func(i, j int) bool {
			return t[i].Price < t[j].Price
		})
	} else if criteria == criteriaDepTime {
		sort.Slice(t, func(i, j int) bool {
			return t[i].DepartureTime.Before(t[j].DepartureTime)
		})
	} else if criteria == criteriaArrTime {
		sort.Slice(t, func(i, j int) bool {
			return t[i].ArrivalTime.Before(t[j].ArrivalTime)
		})
	}
	if len(t) < 3 {
		return t
	}
	return t[0:3]
}
