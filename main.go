package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// URI is the Coinbase base URI
const URI = "https://api.coinbase.com/v2"

// LATESTVERSIONDATE is the date versioning scheme of coinbase
const LATESTVERSIONDATE = "2016-08-10"

// SpotPrice holds Coinbase's data
type SpotPrice struct {
	Currency string
	Amount   string
}

// SpotPriceResponse holds Coinbase response
type SpotPriceResponse struct {
	Data SpotPrice
}

// BadRequestError holds error object
type BadRequestError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Response holds ether_bot's response
type Response map[string]interface{}

// GetSpotPrice fetches current price on Coinbase
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
		error := BadRequestError{
			"Add currency_pair querystring",
			"Bad Request",
		}

		response := Response{
			"data":  nil,
			"error": error,
		}

		j, err := json.Marshal(response)
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
	fmt.Println("Listenting on port")
	fmt.Println(os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
