package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)

}

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func FetchCat() (CatFact, error) {
	url := "https://catfact.ninja/fact"

	res, err := http.Get(url)

	if err != nil {
		return CatFact{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return CatFact{}, fmt.Errorf("external api failed %s", res.Status)
	}

	var data CatFact

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return CatFact{}, err
	}

	return data, nil
}

type Joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func FetchJoke() (Joke, error) {
	url := "https://official-joke-api.appspot.com/random_joke"

	res, err := http.Get(url)

	if err != nil {
		return Joke{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Joke{}, fmt.Errorf("joke api failed")
	}

	var joke Joke

	if err := json.NewDecoder(res.Body).Decode(&joke); err != nil {
		return Joke{}, err
	}

	return joke, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"ok":    false,
			"error": "Only GET method is allowed",
		})

		return
	}

	var (
		cat  CatFact
		joke Joke
		err  error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		cat, err = FetchCat()
	}()

	go func() {
		defer wg.Done()
		joke, err = FetchJoke()
	}()

	wg.Wait()

	if err != nil {
		WriteJSON(w, http.StatusBadGateway, map[string]any{
			"ok":    false,
			"error": "failed to fetch data",
		})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"ok": true,
		"data": map[string]any{
			"cat":  cat.Fact,
			"joke": joke.Setup + " - " + joke.Punchline,
		},
	})

}

func main() {

	http.HandleFunc("/dashboard", dashboardHandler)

	fmt.Println("server starting on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
		return
	}
}
