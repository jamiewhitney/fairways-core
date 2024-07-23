package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"sync"
)

func main() {
	logger := logging.NewLoggerFromEnv()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&multiStatements=true", "root", "root", "localhost", "3306")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal(err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 4)
	wg.Add(4)
	go func() {
		defer wg.Done()
		if err := generateCourses(db); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := generateBookings(db); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := generatePricing(db); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := generateSchedule(db); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		logger.Error(err)
	}
	logger.Info("done")
}

func generateCourses(db *sql.DB) error {
	data := []struct {
		id               int
		name             string
		holes            int64
		lapsed           bool
		live             bool
		street_address_1 string
		city             string
		state            string
		postal_code      string
		country          string
	}{
		{
			name:             "London Golf Club",
			holes:            18,
			lapsed:           false,
			live:             true,
			street_address_1: "London",
			city:             "London",
			state:            "London",
			postal_code:      "TEST",
			country:          "United Kingdom",
		},
		{
			name:             "Chester Golf Club",
			holes:            18,
			lapsed:           false,
			live:             true,
			street_address_1: "Chester",
			city:             "Chester",
			state:            "Cheshire",
			postal_code:      "TEST",
			country:          "United Kingdom",
		},
	}

	for _, course := range data {
		_, err := db.Exec("INSERT INTO catalog.courses (name, holes, live, street_address_1, city, state, postal_code, country) VALUES (?, ?,?, ?, ?, ?, ?, ?);", course.name, course.holes, course.live, course.street_address_1, course.city, course.state, course.postal_code, course.country)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateBookings(db *sql.DB) error {
	data := []struct {
		user_id           string
		course_id         int64
		golfers           int64
		datetime          string
		price             int64
		booking_id        string
		stripe_payment_id string
		status            string
		confirmed         bool
	}{
		{
			user_id:           "1234",
			course_id:         1,
			golfers:           2,
			datetime:          "2024-12-12",
			price:             16,
			booking_id:        "1234",
			stripe_payment_id: "1234",
			status:            "confirmed",
			confirmed:         true,
		},
	}

	for _, course := range data {
		_, err := db.Exec("INSERT INTO bookings.bookings (user_id, course_id, golfers, datetime, price, booking_id, stripe_payment_id, status, confirmed) VALUES (?, ?, ?,?, ?, ?, ?, ?, ?);", course.user_id, course.course_id, course.golfers, course.datetime, course.price, course.booking_id, course.stripe_payment_id, course.status, course.confirmed)
		if err != nil {
			return err
		}
	}
	return nil
}

func generatePricing(db *sql.DB) error {
	data := []struct {
		price_band int64
		week       int64
		hour       int64
		year       int64
		band       int64
		course_id  int64
		price      int64
	}{{
		price_band: 1,
		week:       1,
		hour:       10,
		band:       1,
		course_id:  1,
		price:      22,
		year:       2024,
	}}

	for _, price := range data {
		_, err := db.Exec("INSERT INTO pricing.course_base_prices(course_id, price, year) VALUES (?, ?, ?);", price.course_id, price.price, price.year)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO pricing.course_price_bands(band, course_id) VALUES (?, ?);", price.band, price.course_id)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO pricing.course_pricing_rules(price_band, week, hour) VALUES (?, ?, ?);", price.band, price.week, price.hour)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateSchedule(db *sql.DB) error {
	data := []struct {
		course_id  int64
		start_time string
		end_time   string
		occurrence int64
		day        int64
	}{
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "16:00:00",
			occurrence: 12,
			day:        0,
		},
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "16:00:00",
			occurrence: 12,
			day:        1,
		},
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "16:00:00",
			occurrence: 12,
			day:        2,
		},
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "16:00:00",
			occurrence: 12,
			day:        3,
		},
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "16:00:00",
			occurrence: 12,
			day:        4,
		},
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "17:00:00",
			occurrence: 12,
			day:        5,
		},
		{
			course_id:  1,
			start_time: "09:00:00",
			end_time:   "18:00:00",
			occurrence: 12,
			day:        6,
		},
	}

	for _, schedule := range data {
		_, err := db.Exec("INSERT INTO tee_times.schedule(course_id, start_time, end_time, occurrence, day, buffer) VALUES (?, ?, ?, ?, ?, ?);", schedule.course_id, schedule.start_time, schedule.end_time, schedule.occurrence, schedule.day, 5)
		if err != nil {
			return err
		}
	}

	return nil
}
