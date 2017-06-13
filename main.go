package main

import (
	"encoding/json"
	"fmt"
	"time"
	"io"
	"io/ioutil"
	"net/http"
)

const URI = "https://api.coinbase.com/v2"
const LATEST_VERSION_DATE = "2016-08-10"

type SpotPrice struct {
	Currency string
	Amount string
}

type SpotPriceResponse struct {
	Data SpotPrice
}

type BadRequestError struct {
	Message string `json:"message"`
	Status string  `json:"status"`
}

func GetSpotPrice(c string) *SpotPriceResponse {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	fmt.Println("Making request")
	res, err := client.Get(fmt.Sprintf("%v/prices/%v/spot", URI, c))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	response := SpotPriceResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	return &response
}

func health(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func fetchCoinbasePrice(w http.ResponseWriter, r *http.Request) {
	
	currencyPair := r.URL.Query().Get("currency_pair")
	if currencyPair == "" {
		badRequest := BadRequestError{
			"Add currency_pair querystring",
			"Bad Request",
		}
		fmt.Println(badRequest)
		j, err := json.Marshal(badRequest)
		if err != nil {
			fmt.Println("Error")
			fmt.Println(err)
		}
		fmt.Println("Writing")
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}

func main() {
	http.HandleFunc("/", health)
	http.HandleFunc("/prices", fetchCoinbasePrice)
	fmt.Println("Listenting on port 8000")
	http.ListenAndServe(":8000", nil)
}
