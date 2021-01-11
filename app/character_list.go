package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type character struct {
	app.Compo

	loader     bool
	characters AllCharacters
}

type AllCharacters struct {
	Info struct {
		Count int         `json:"count"`
		Pages int         `json:"pages"`
		Next  string      `json:"next"`
		Prev  interface{} `json:"prev"`
	} `json:"info"`
	Results []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Status  string `json:"status"`
		Species string `json:"species"`
		Type    string `json:"type"`
		Gender  string `json:"gender"`
		Origin  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"origin"`
		Location struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"location"`
		Image   string    `json:"image"`
		Episode []string  `json:"episode"`
		URL     string    `json:"url"`
		Created time.Time `json:"created"`
	} `json:"results"`
}

func (c *character) getAllCharacters(url string) {
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.Log(err.Error())
		return
	}

	var all AllCharacters

	err = json.Unmarshal(b, &all)
	if err != nil {
		app.Log(err.Error())
		return
	}
	// c.isLoader()
	time.AfterFunc(2*time.Second, c.isLoader)
	c.updateAllCharacters(all)
}

func (c *character) updateAllCharacters(data AllCharacters) {
	app.Dispatch(func() {
		c.characters = data
		c.Update()
	})
}

func (c *character) isLoader() {
	app.Dispatch(func() {
		c.loader = true
		c.Update()
	})
}

func (c *character) OnMount(ctx app.Context) {
	app.Dispatch(func() {
		c.getAllCharacters("https://rickandmortyapi.com/api/character")
	})
}

func (c *character) onNext(ctx app.Context, e app.Event) {
	c.loader = false
	c.Update()

	app.Dispatch(func() {
		c.getAllCharacters(c.characters.Info.Next)
	})
}

func (c *character) Render() app.UI {
	return app.Div().Class("columns is-multiline").Body(
		app.If(!c.loader,
			&loader{},
		).Else(

			app.Range(c.characters.Results).Slice(func(i int) app.UI {
				return app.Div().Class("column is-6").Body(

					app.Div().Class("box").Body(

						app.Article().Class("media").Body(
							app.Div().Class("media-left").Body(
								app.Figure().Class("image is-128x128").Body(
									app.Img().Src(c.characters.Results[i].Image),
								),
							),

							app.Div().Class("media-content").Body(
								app.Div().Class("content").Body(
									app.P().Body(
										app.Strong().Text(c.characters.Results[i].Name),
									),
								),
							),
						),
					),
				)
			}),
		),
		app.Button().Class("button").Text("Next").OnClick(c.onNext),
	)
}