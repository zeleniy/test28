package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/go-faker/faker/v4"
	"github.com/zeleniy/test28/bootstrap"
	"github.com/zeleniy/test28/database/seeders"
	"github.com/zeleniy/test28/models"
)

func main() {

	bootstrap.SetUpDb(os.Getenv("DB_URL"))

	seeder := &seeders.Seeder{}
	seeder.MinSubscriptionsToSeed = 15
	seeder.MinUsersToSeed = 10

	seeder.RandomUser = func() (*models.User, error) {
		return &models.User{
			Login:        faker.Username(),
			PasswordHash: faker.Password(),
		}, nil
	}

	seeder.RandomSubscription = func() (*models.Subscription, error) {
		return &models.Subscription{
			ServiceName: []string{"Okko", "Yandex", "Wink", "Sber", "Ivi"}[rand.Intn(5)],
			Price:       []int{10, 20, 30, 40, 50}[rand.Intn(5)],
			StartDate:   randomDateInRange(time.Now().AddDate(-1, 0, 0), time.Now()),            // any time in in 1 year before now
			EndDate:     null.TimeFrom(time.Now().AddDate(0, []int{1, 6, 12}[rand.Intn(3)], 0)), // 1, 6 or 12 months duration
		}, nil
	}

	err := seeder.Run(context.Background(), boil.GetContextDB())

	if err != nil {
		panic(err)
	}
}

func randomDateInRange(start, end time.Time) time.Time {
	delta := end.Unix() - start.Unix()
	randomUnix := rand.Int63n(delta) + start.Unix()
	return time.Unix(randomUnix, 0)
}
