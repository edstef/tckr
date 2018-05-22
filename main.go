package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
)

// type TimeSeriesData struct {
//   OneOpen    string `json:"1. open"`
//   TwoHigh    string `json:"2. high"`
//   ThreeLow   string `json:"3. low"`
//   FourClose  string `json:"4. close"`
//   FiveVolume string `json:"5. volume"`
// }

var AVKey string

type StockResponse struct {
	MetaData struct {
		Information     string `json:"1. Information"`
		Symbol          string `json:"2. Symbol"`
		LastRefreshed   string `json:"3. Last Refreshed"`
		Interval        string `json:"4. Interval"`
		OutputSize      string `json:"5. Output Size"`
		TimeZone        string `json:"6. Time Zone"`
  } `json:"Meta Data"`
  TimeSeriesData map[string]map[string]string `json:"Time Series (1min)"`
}

func fetchStock(sym string) {

  AVURL := "https://www.alphavantage.co/query?"
  AVFn := "function=TIME_SERIES_INTRADAY"
  AVSym := "&symbol=" + sym
  AVInterval := "&interval=1min"
  AVFullKey := "&apikey=" + AVKey

  resp, err := http.Get(AVURL + AVFn + AVSym + AVInterval + AVFullKey)

  if err != nil {
    fmt.Println("Error Fetching")
  }

  data := new(StockResponse)

  json.NewDecoder(resp.Body).Decode(data)

  lastRefreshed := data.MetaData.LastRefreshed


  fmt.Println(data.TimeSeriesData[lastRefreshed]["1. open"])
}

func getDefault(w http.ResponseWriter, r *http.Request) {
  fetchStock("VOO")
}

func getStock(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  fetchStock(params["stock"])
}

func main() {

  tmpKey, _ := ioutil.ReadFile("AVKEY.txt")

  AVKey = string(tmpKey)
  
  router := mux.NewRouter()

  router.HandleFunc("/", getDefault).Methods("GET")
  router.HandleFunc("/{stock}", getStock).Methods("GET")

  fmt.Println("Server started on port 8000")  
  log.Fatal(http.ListenAndServe(":8000", router))
}
