# summer-2022-practice-task

**Завдання**:

Є список потягів з інформацією про номер потяга, станцію відправлення, станцію прибуття, вартість проїзду, час відправлення, час прибуття.
Програма має запитати користувача про:
1. станцію відправлення
2. станцію прибуття
3. критерій, по якому сортувати результат

Базуючись на цьому, програма має знайти відповідні до умов потяги та відсортувати їх відносно заданому критерію.
Повернути треба 3 перших потяги. Проїзд без пересадок - з точки відправлення до точки брибуття має довзти один і той же потяг.
Якщо таких потягів немає - треба повернути nil.

**Приклад**
Це просто абстрактний приклад для вашого розуміння. Вхідні та вихідні параметри описані окремо. 

Вхідіні дані:

    departureStation: "1902"
    arrivalStation: "1929"
    criteria:   "price"

Результат:

 	{TrainID: 1177, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 164.65, ArrivalTime: time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 36, 0, 0, time.UTC)},
	{TrainID: 1178, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 164.65, ArrivalTime: time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 36, 0, 0, time.UTC)},
	{TrainID: 1141, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 176.77, ArrivalTime: time.Date(0, time.January, 1, 12, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 48, 0, 0, time.UTC)},

**Скоуп роботи**:

Дано файл зі списком потягів - data.json.

Потрібно написати код, котрий запитає юзера про вхідні параметри: станцію відправлення, станцію прибуття, критерій, по якому сортувати результат.
Ці дані мають бути передані в функцію FindTrains(depStation, arrStation, criteria string) (Trains, error).

Потрібно написати код функції FindTrains(depStation, arrStation, criteria string) (Trains, error).
Вона має вичитати дані потягів з цього файлу.
Приймає вхідні параметри від користувача:
1. **depStation** string - номер станції відправлення. Валідні значення - будь-яке натуральне число.
2. **arrStation** string - номер станції прибуття. Валідні значення - будь-яке натуральне число.
3. **criteria** string - критерій, по котрому треба відсортувати потяги в результаті. Валідні значення: price, arrival-time, departure-time.


    price - спершу дешевші
    arrival-time - спершу ті, що раніше прибувають
    departure-time - спершу ті, що раніше відправляються
Ці параметри є обовʼязковими. Якщо вони відсутні чи невалідні - має повернутися одна із відповідних помилок (змінювати текст заборонено):

    "unsupported criteria"
    "empty departure station"
    "empty arrival station"
    "bad arrival station input"
    "bad departure station input"

Отже, в функції main() має бути приблизно таке:

    func main() {
        // ... запит даних від користувача
        result, err := FindTrains(departureStation, arrivalStation, criteria))
        // ... обробка помилки
        // ... друк result
    }
    
Весь код має знаходитись в main.go.

**Потрібні типи** (при необхідності дозволено змінювати/додавати теги, проте, змінювати імена полів/типів чи типи даних - заборонено):

    type Trains []Train

    type Train struct {
        TrainID            int
        DepartureStationID int
        ArrivalStationID   int
        Price              float32
        ArrivalTime        time.Time
        DepartureTime      time.Time
    }
