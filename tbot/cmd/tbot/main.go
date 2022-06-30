package main

import (
	"flag"
	"fmt"

	db "github.com/emirem/game-scrapper/tbot/internal/db"
	gms "github.com/emirem/game-scrapper/tbot/internal/games"
	twitter "github.com/emirem/game-scrapper/tbot/internal/tbot"
)

var defaultStores = []string{"steam", "epic", "ubisoft"}

func weekAnalysis() {
	thisWeekData, err := db.GetThisWeekReleases()

	if err != nil {
		fmt.Println("Could not get last two days data.", err)
		return
	}

	if len(thisWeekData) > 5 {
		thisWeekData = thisWeekData[0:5]
	}

	// Post NewGamesTweet
	if len(thisWeekData) > 0 {
		tweetBody := twitter.ConstructThisWeekReleasesTweet(thisWeekData)
		fmt.Println(tweetBody)
		twitter.PostTweet(tweetBody)
	}
}

func runStore(games *gms.Games, storeId string) {
	percentage := "60"
	newOnStore := games.GetNewGamesByStore(storeId)
	newFreeOnStore := games.GetNewFreeGamesByStore(storeId)
	largeSalesOnStore := games.GetLargeSalesByStore(storeId, percentage)

	if len(newOnStore) > 5 {
		newOnStore = newOnStore[0:5]
	}

	if len(newFreeOnStore) > 5 {
		newFreeOnStore = newFreeOnStore[0:5]
	}

	if len(largeSalesOnStore) > 5 {
		largeSalesOnStore = largeSalesOnStore[0:5]
	}

	// Post NewGamesTweet
	if len(newOnStore) > 0 {
		tweetBody := twitter.ConstructNewGamesTweet(newOnStore, storeId)
		fmt.Println(tweetBody)
		twitter.PostTweet(tweetBody)
	}

	// Post NewFreeGamesTweet
	if len(newFreeOnStore) > 0 {
		tweetBody := twitter.ConstructNewFreeGamesTweet(newFreeOnStore, storeId)
		fmt.Println(tweetBody)
		twitter.PostTweet(tweetBody)
	}

	// Post LargeDiscountsTweet
	if len(largeSalesOnStore) > 0 {
		tweetBody := twitter.ConstructLargeSalesTweet(largeSalesOnStore, storeId, percentage)
		fmt.Println(tweetBody)
		twitter.PostTweet(tweetBody)
	}
}

func main() {
	var storeId string
	var includeWeekAnalysis, skipStores bool

	flag.StringVar(&storeId, "sid", "", "Store id")
	flag.BoolVar(&skipStores, "skip", false, "Skips store tweets")
	flag.BoolVar(&includeWeekAnalysis, "w", false, "Includes week analysis")

	flag.Parse()

	games := gms.Games{}
	data, err := db.GetLastTwoDaysData()

	if err != nil {
		fmt.Println("Could not get last two days data.", err)
		return
	}

	games.SetData(data)

	if skipStores == false {
		if storeId != "" {
			runStore(&games, storeId)
		} else {
			// Post for all stores
			for _, storeId := range []string{"steam", "epic", "ubisoft"} {
				runStore(&games, storeId)
			}
		}
	}

	if includeWeekAnalysis == true {
		weekAnalysis()
	}
}
