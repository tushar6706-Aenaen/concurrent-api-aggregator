package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func FetchData() (CatFact, error) {
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

	if err := json.NewDecoder(res.Body).Decode(&data) ; err != nil{
		return CatFact{} , err
	}

	return data, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"ok":    false,
			"error": "Only GET method is allowed",
		})

		return
	}

	data, err := FetchData()

	if err != nil {

		WriteJSON(w, http.StatusBadGateway, map[string]any{
			"ok":    false,
			"error": "failed to fetch external data",
		})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"ok":   true,
		"fact": data.Fact,
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
