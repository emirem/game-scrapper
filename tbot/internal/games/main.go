package games

import (
	"fmt"
	"strings"
	"time"
)

type Data struct {
	Id              int
	Title           string
	Details         string
	Img_url         string
	Url             string
	Price           string
	Discount_amount string
	Store_id        string
	Category_id     string
	Release_date    string
	Date_created    string
}

type Games struct {
	data       []*Data
	mostRecent []*Data
}

func (game *Games) SetData(data []*Data) {
	game.data = data

	recent, err := game.GetMostRecentGames()

	if err != nil {
		fmt.Println("Could not set most recent data", err)
		return
	}

	game.mostRecent = recent
}

func (game *Games) GetMostRecentGames() ([]*Data, error) {
	todaysGames := game.GetNewGamesByDate(time.Now())

	if len(todaysGames) > 0 {
		return todaysGames, nil
	}

	return game.data, nil
}

func (game *Games) GetNewGamesByDate(date time.Time) []*Data {
	var result []*Data

	for _, game := range game.data {
		parsedDate, err := time.Parse("2006-01-02 15:04:05", game.Date_created)

		if err != nil {
			fmt.Println("Date parse failed", err)
		}

		if parsedDate.Day() == date.Day() {
			result = append(result, game)
		}
	}

	return result
}

func (game *Games) GetNewGamesByStore(storeId string) []*Data {
	var result []*Data
	todaysGames := game.GetNewGamesByDate(time.Now())
	yesterdayGames := game.GetNewGamesByDate(time.Now().AddDate(0, 0, -1))

	// TODO: optimize diff two days
	for _, tGame := range todaysGames {
		if tGame.Store_id == storeId {
			var exists = false

			for _, yGame := range yesterdayGames {
				if tGame.Title == yGame.Title {
					exists = true
					break
				}
			}

			if exists == false {
				result = append(result, tGame)
			}
		}
	}

	return result
}

func (game *Games) GetNewFreeGamesByStore(storeId string) []*Data {
	var result []*Data
	todaysGames := game.GetNewGamesByDate(time.Now())
	yesterdayGames := game.GetNewGamesByDate(time.Now().AddDate(0, 0, -1))

	// TODO: optimize diff two days
	for _, tGame := range todaysGames {
		if tGame.Store_id == storeId && tGame.Price == "Free" {
			var exists = false

			for _, yGame := range yesterdayGames {
				if tGame.Title == yGame.Title {
					exists = true
					break
				}
			}

			if exists == false {
				result = append(result, tGame)
			}
		}
	}

	return result
}

func (game *Games) GetRecentByStore(limit int, storeId string) []*Data {
	var result []*Data

	for _, game := range game.mostRecent {
		if game.Store_id == storeId {
			result = append(result, game)
		}

		if len(result) == limit {
			break
		}
	}

	return result
}

func (game *Games) GetLargeSalesByStore(storeId, percentage string) []*Data {
	var result []*Data
	todaysGames := game.GetNewGamesByDate(time.Now())
	replacer := strings.NewReplacer("-", "", "%", "")

	for _, tGame := range todaysGames {
		if tGame.Store_id == storeId {
			var parsed string
			parsed = replacer.Replace(tGame.Discount_amount)

			if parsed >= percentage {
				result = append(result, tGame)
			}
		}
	}

	return result
}
