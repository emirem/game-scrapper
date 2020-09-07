package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	gms "github.com/emirem/game-scrapper/tbot/internal/games"
	_ "github.com/go-sql-driver/mysql"
)

func getConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?tls=true", os.Getenv("PS_USERNAME"), os.Getenv("PS_PASSWORD"), os.Getenv("PS_HOST"), os.Getenv("PS_DBNAME")))

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func parseRows(rows *sql.Rows) []*gms.Data {
	var data []*gms.Data

	for rows.Next() {
		var id int
		var title, details, img_url, url, price, discount_amount, store_id, category_id, release_date, date_created string

		err := rows.Scan(&id, &title, &details, &img_url, &url, &price, &discount_amount, &store_id, &category_id, &release_date, &date_created)

		if err != nil {
			fmt.Println("Error reading row", err)
		}

		gameItem := &gms.Data{
			Id:              id,
			Title:           title,
			Details:         details,
			Img_url:         img_url,
			Url:             url,
			Price:           price,
			Discount_amount: discount_amount,
			Store_id:        store_id,
			Category_id:     category_id,
			Release_date:    release_date,
			Date_created:    date_created,
		}

		data = append(data, gameItem)
	}

	return data
}

func GetLastTwoDaysData() ([]*gms.Data, error) {
	db, err := getConnection()

	if err != nil {
		fmt.Println("Database failed to initialize.", err)
		return nil, err
	}

	statement, err := db.Prepare("SELECT * FROM games WHERE date_created BETWEEN ? AND ?")

	if err != nil {
		fmt.Println("Statement init failed.", err)
		return nil, err
	}

	defer statement.Close()

	rows, err := statement.Query(time.Now().AddDate(0, 0, -1), time.Now())

	if err != nil {
		fmt.Println("Statement query failed.", err)
		return nil, err
	}

	defer rows.Close()
	defer db.Close()

	data := parseRows(rows)

	return data, nil
}

func GetThisWeekReleases() ([]*gms.Data, error) {
	db, err := getConnection()

	if err != nil {
		fmt.Println("Database failed to initialize.", err)
		return nil, err
	}

	statement, err := db.Prepare("SELECT * FROM games WHERE release_date BETWEEN ? AND ?")

	if err != nil {
		fmt.Println("Statement init failed.", err)
		return nil, err
	}

	defer statement.Close()

	now := time.Now()
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-int(time.Monday), 0, 0, 0, 0, time.Local)
	weekEnd := time.Date(now.Year(), now.Month(), now.Day()+int(time.Friday), 0, 0, 0, 0, time.Local)

	rows, err := statement.Query(weekStart, weekEnd)

	if err != nil {
		fmt.Println("Statement query failed.", err)
		return nil, err
	}

	defer rows.Close()
	defer db.Close()

	data := parseRows(rows)

	return data, nil
}
