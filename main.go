package main

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
	//	... запит даних від користувача
	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	//	... обробка помилки
	//	... друк result
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	return nil, nil // маєте повернути правильні значення
}
