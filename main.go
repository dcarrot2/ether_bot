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

	"github.com/garyburd/redigo/redis"
)

// AvailableCryptoCurrencies is a lookup for valid currency queries
var AVAILABLECRYPTOCURRENCIES = map[string]bool{
	"ETH": true,
}

// URI is the Coinbase base URI
const URI = "https://api.coinbase.com/v2"

// REDISADDRESS is the port for redis server
const REDISADDRESS = ":6379"

// MAXCONNECTIONS is the max connections to redis
const MAXCONNECTIONS = 10

// LATESTVERSIONDATE is the date versioning scheme of coinbase
const LATESTVERSIONDATE = "2016-08-10"

// SpotPrice holds Coinbase's data
type SpotPrice struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
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
	log.Println("Making request")
	res, err := client.Get(fmt.Sprintf("%v/prices/%v/spot", URI, c))
	if err != nil {
		log.Println(err)
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	response := SpotPriceResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()
	return &response
}

func health(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func parseSlackWebhook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form")
	}
	coin := r.PostFormValue("text")
	if _, ok := AVAILABLECRYPTOCURRENCIES[coin]; ok {
		log.Println("Valid query, executing")
		w.Write([]byte("OK"))
	}
}

func rawPriceHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("Error")
			log.Println(err)
		}
		log.Println("Writing")
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	} else {
		res := GetSpotPrice(currencyPair)
		if res != nil {
			response := Response{
				"data":  res.Data,
				"error": nil,
			}
			j, err := json.Marshal(response)
			if err != nil {
				log.Fatalf("Cannot parse GetSpotPrice: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		}
	}
}

func fetchCoinbasePrice(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		rawPriceHandler(w, r)
	case "POST":
		parseSlackWebhook(w, r)
	default:
		error := BadRequestError{
			"Bad request",
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

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}

}

func main() {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", REDISADDRESS)

		if err != nil {
			return nil, err
		}

		return c, err
	}, MAXCONNECTIONS)

	defer redisPool.Close()
	http.HandleFunc("/", health)
	http.HandleFunc("/prices", fetchCoinbasePrice)
	http.HandleFunc("/slack-prices", parseSlackWebhook)
	log.Println("Listening on port")
	log.Println(os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
