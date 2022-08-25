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
	//errors
	unsupportedCriteria      = errors.New("unsupported criteria")
	emptyDepartureStation    = errors.New("empty departure station")
	emptyArrivalStation      = errors.New("empty arrival station")
	badArrivalStationInput   = errors.New("bad arrival station input")
	badDepartureStationInput = errors.New("bad departure station input")
	//variables for train search
	neededDepartureId                          int
	neededArrivalId                            int
	departureStation, arrivalStation, criteria string
)

const (
	maxNaturalNumber   = 1          //the constant is required to validate user-entered values
	resultLenghtNeeded = 3          //required number of trains as a result
	layout             = "15:04:05" //A special layout parameter is required to describe the format of the time value.
	//It must be a reference date/time — Mon Jan 2 15:04:05 MST 2006, which is formatted as the expected date/time
)

//The getContentFromJson function reads a file and, if successful, returns a slice of bytes from the file
func getContentFromJson(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

//The ConvertStringToTime function, if successful, converts the time from a string format to a time format
func ConvertStringToTime(str string) (time.Time, error) {
	convertedTime, err := time.Parse(layout, str)
	if err != nil {
		panic(err)
	}
	return convertedTime, nil
}

//override the UnmarshalJSON method to work with our custom time.Time types
func (t *Train) UnmarshalJSON(content []byte) error {
	//create a structure with fields similar to the Train structure, but without custom types
	type structWithNormalTypes struct {
		TrainID            int     `json:"trainId"`
		DepartureStationID int     `json:"departureStationId"`
		ArrivalStationID   int     `json:"arrivalStationId"`
		Price              float32 `json:"price"`
		ArrivalTime        string  `json:"arrivalTime"`
		DepartureTime      string  `json:"departureTime"`
	}
	//perform Unmarshal into a structure without custom types
	var tmp structWithNormalTypes
	err := json.Unmarshal(content, &tmp)
	if err != nil {
		return err
	}
	//assign the fields of the Beh structure of custom types to the Train structure (with custom types)
	//at the moment of assigning a field with the time.Time type, we apply it to the corresponding field of the structure
	//without custom types function ConvertStringToTime
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

//The departureStationValidation function checks whether an empty string was entered for the departure station,
// whether it is an integer, and whether it is a natural number. If at least one of the conditions is met, it returns an error
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

//The arrivalStationValidation function checks whether an empty string was entered for the arrival station,
//whether it is an integer, and whether it is a natural number. If at least one of the conditions is met, it returns an error
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

//The criteriaValidation function checks the validity of the entered criterion.
//We create a map of valid criteria, where the key is the criterion itself. 
//The key value will not be used, so we use the struct{} type (0 bytes are used for its storage)
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
	//	... error handling
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	//	... print result
	for _, train := range result {
		formatResult(train)
	}

}

//The FindTrains function finds trains in the file that satisfy the parameters entered by the user
//and,if successful, returns no more than 3 options sorted by the entered criteria.
//If the user input is valid, but there are no trains for these parameters, then returns nil,nil
func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	var err error
	//Receiving data
	fmt.Println("Enter the id of the departure station.")
	fmt.Scanln(&departureStation)
	err = departureStationValidation(departureStation)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	fmt.Println("Enter the id of the arrival station.")
	fmt.Scanln(&arrivalStation)
	err = arrivalStationValidation(arrivalStation)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	fmt.Println("Enter the criteria by which the resulting trains should be sorted. Valid values: price, arrival-time, departure-time")
	fmt.Scanln(&criteria)
	err = criteriaValidation(criteria)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	var t []Train
	//read json
	content, err := getContentFromJson("data.json")
	if err != nil {
		fmt.Println("open file error: " + err.Error())
		return nil, nil
	}
	//perform a custom unmarshal
	err = json.Unmarshal(content, &t)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return nil, nil
	}

	//we find all trains matching the stations entered by the user
	var foundedTrains Trains
	for _, train := range t {
		if train.DepartureStationID == neededDepartureId && train.ArrivalStationID == neededArrivalId {
			foundedTrains = append(foundedTrains, train)
		}
	}
	foundedTrains.sortFoundedTrains(criteria)
	//Depending on the number of found trains, we return different numbers of them
	currentLenght := len(foundedTrains)
	if currentLenght == 0 {
		fmt.Println("Nothing found")
		return nil, nil
	}
	if currentLenght < resultLenghtNeeded {
		return foundedTrains, nil
	}
	return foundedTrains[:resultLenghtNeeded], nil
}

//The sortFoundedTrains method sorts the trains found by the entered parameters
//sort.SliceStable is used to preserve the original order of the trains
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

//The formatResult function formats the structure printout
func formatResult(t Train) {
	fmt.Printf(`Train ID: %v, Departure station ID: %v, Arrival station ID: %v, Price: %v, Departure time: %s, Arrival time: %s.%s`,
		t.TrainID, t.DepartureStationID, t.ArrivalStationID, t.Price, t.DepartureTime.Format(layout), t.ArrivalTime.Format(layout),
		"\n")
}
