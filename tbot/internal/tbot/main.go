package twitter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	gms "github.com/emirem/game-scrapper/tbot/internal/games"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

type NewGamesTemplate struct {
	Games      []*gms.Data
	StoreName  string
	Percentage string
}

func PostTweet(text string) (*types.CreateOutput, error) {
	c, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           os.Getenv("TB_ACCESS_TOKEN"),
		OAuthTokenSecret:     os.Getenv("TB_ACCESS_TOKEN_SECRET"),
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tweetData := &types.CreateInput{Text: gotwi.String(text)}
	res, err := managetweet.Create(context.Background(), c, tweetData)

	if err != nil {
		fmt.Println("Tweet create failed.", err)
		return nil, err
	}

	return res, nil
}

func ConstructNewGamesTweet(games []*gms.Data, storeName string) string {
	data := &NewGamesTemplate{
		Games:     games,
		StoreName: storeName,
	}

	return constructTweetTemplate(data, "./internal/tbot/newGamesTpl.txt")
}

func ConstructNewFreeGamesTweet(games []*gms.Data, storeName string) string {
	data := &NewGamesTemplate{
		Games:     games,
		StoreName: storeName,
	}

	return constructTweetTemplate(data, "./internal/tbot/newFreeGamesTpl.txt")
}

func ConstructThisWeekReleasesTweet(games []*gms.Data) string {
	data := &NewGamesTemplate{
		Games: games,
	}

	return constructTweetTemplate(data, "./internal/tbot/thisWeekReleases.txt")
}

func ConstructLargeSalesTweet(games []*gms.Data, storeName, percentage string) string {
	data := &NewGamesTemplate{
		Games:      games,
		StoreName:  storeName,
		Percentage: percentage,
	}

	return constructTweetTemplate(data, "./internal/tbot/largeSaleTpl.txt")
}

// TODO: Figure out generics
func constructTweetTemplate(Tpl interface{}, templatePath string) string {
	var temp *template.Template
	var result bytes.Buffer

	temp = template.Must(template.ParseFiles(templatePath))
	err := temp.Execute(&result, Tpl)

	if err != nil {
		fmt.Println("Template execute failed", err)
	}

	return result.String()
}
