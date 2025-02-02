package main

import (
	"fmt"
	"math"

	"time"

	"golang.org/x/exp/rand"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var NO_OF_BLIND_PARAMS_SUPPORTED = 65535

func main() {
	db = initDB()

	const NO_OF_USERS = 10_00_000   // 1M
	const NO_OF_TODOS = 1_00_00_000 // 10M
	startedTime := time.Now().Format("15:04:05")

	// batchedBulkUserInsert(NO_OF_USERS)
	bachedBulkTodosInsert(NO_OF_TODOS)

	fmt.Println("Time started:", startedTime)
	fmt.Println("Time ended:", time.Now().Format("15:04:05"))
}

func initDB() *gorm.DB {
	dsn := "user=postgres dbname=vyson_db port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	return db
}

func batchedBulkUserInsert(noOfUsers int) {
	NO_OF_BLIND_PARAMS_USED := 5 // Name, Email, CreatedAt, UpdatedAt, DeletedAt
	MAX_USERS_PER_BATCH := NO_OF_BLIND_PARAMS_SUPPORTED / NO_OF_BLIND_PARAMS_USED

	USER_BATCH_SIZE := math.Min(float64(MAX_USERS_PER_BATCH), float64(noOfUsers))
	NO_OF_BATCHES := int(math.Ceil(float64(noOfUsers) / USER_BATCH_SIZE))

	for batchNumber := 0; batchNumber < NO_OF_BATCHES; batchNumber++ {
		users := []*User{}

		var NO_OF_ITERATIONS int
		if noOfUsers > MAX_USERS_PER_BATCH && batchNumber == NO_OF_BATCHES-1 {
			NO_OF_ITERATIONS = noOfUsers % int(USER_BATCH_SIZE)
		} else {
			NO_OF_ITERATIONS = int(USER_BATCH_SIZE)
		}

		for i := 0; i < NO_OF_ITERATIONS; i++ {
			var maxUserId int
			db.Model(&User{}).Select("MAX(id)").Scan(&maxUserId)

			user := User{
				Name:  fmt.Sprintf("User %d", maxUserId+i+1),
				Email: fmt.Sprintf("user%d@gmail.com", maxUserId+i+1),
			}
			users = append(users, &user)
		}

		result := db.Create(&users)

		if result.Error != nil {
			panic(result.Error)
		}
	}
}

func bachedBulkTodosInsert(noOfTodos int) {
	NO_OF_BLIND_PARAMS_USED := 9 // Title, UserId, CompletedAt, DueDate, Description, Status, CreatedAt, UpdatedAt, DeletedAt
	MAX_TODOS_PER_BATCH := NO_OF_BLIND_PARAMS_SUPPORTED / NO_OF_BLIND_PARAMS_USED

	TODOS_BATCH_SIZE := math.Min(float64(MAX_TODOS_PER_BATCH), float64(noOfTodos))
	NO_OF_BATCHES := int(math.Ceil(float64(noOfTodos) / TODOS_BATCH_SIZE))

	var MIN_USER_ID, MAX_USER_ID int
	db.Model(&User{}).Select("MIN(id)").Scan(&MIN_USER_ID)
	db.Model(&User{}).Select("MAX(id)").Scan(&MAX_USER_ID)

	for batchNumber := 0; batchNumber < NO_OF_BATCHES; batchNumber++ {
		todos := []*Todo{}

		var NO_OF_ITERATIONS int
		if noOfTodos > MAX_TODOS_PER_BATCH && batchNumber == NO_OF_BATCHES-1 {
			NO_OF_ITERATIONS = noOfTodos % int(TODOS_BATCH_SIZE)
		} else {
			NO_OF_ITERATIONS = int(TODOS_BATCH_SIZE)
		}

		for i := 0; i < NO_OF_ITERATIONS; i++ {
			var maxID int
			db.Model(&Todo{}).Select("MAX(id)").Scan(&maxID)

			randomUserId := MIN_USER_ID + rand.Intn(MAX_USER_ID-MIN_USER_ID+1)
			todo := Todo{
				Title:       fmt.Sprintf("Todo %d", maxID+i+1),
				UserId:      uint(randomUserId),
				Description: fmt.Sprintf("Description for todo %d", maxID+i+1),
				Status:      PENDING,
			}

			rand.Seed(uint64(time.Now().UnixNano()))

			randomValue := rand.Intn(4) + 1

			createdAtTime := getRandomPastTime(2*24*time.Hour, 30*24*time.Hour)
			dueTime := createdAtTime.Add(24 * time.Hour)

			switch randomValue {
			case 1:
				// not due, pending
				dueTime = createdAtTime.Add(365 * 24 * time.Hour)
			case 2:
				// not due, in progress
				dueTime = createdAtTime.Add(365 * 24 * time.Hour)
				todo.Status = IN_PROGRESS
			case 3:
				// past due, completed
				todo.Status = COMPLETED
				completedAt := createdAtTime.Add(12 * time.Hour)
				todo.CompletedAt = &completedAt
			case 4:
				// due last month not completed
				createdAtTime = getRandomPastTime(31*24*time.Hour, 60*24*time.Hour)
				dueTime = createdAtTime.Add(24 * time.Hour)
			}

			todo.CreatedAt = createdAtTime
			todo.UpdatedAt = createdAtTime
			todo.DueDate = &dueTime

			todos = append(todos, &todo)
		}

		result := db.Create(&todos)

		if result.Error != nil {
			panic(result.Error)
		}
	}
}

func getRandomPastTime(minDuration, maxDuration time.Duration) time.Time {
	now := time.Now()

	randomDuration := minDuration + time.Duration(rand.Int63n(int64(maxDuration-minDuration)))
	randomPastTime := now.Add(-randomDuration)

	return randomPastTime
}
