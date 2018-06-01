package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  // "io/ioutil"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
)

const LineLength = 22

type StockResponse struct {
  Symbol            string  `json:"symbol`
  CompanyName       string  `json:"companyName"`
  PrimaryExchange   string  `json:"primaryExchange"`
  High              float64  `json:"high"`
  Low               float64  `json:"low"`
  LatestPrice       float64 `json:"latestPrice"`
  LatestSource      string  `json:"latestSource"`
  LatestUpdate      int64   `json:"latestUpdate"`
  Change            float64 `json:"change"`
  ChangePercent     float64 `json:"changePercent"`
  Volume            int64   `json:"latestVolume"`
  PE                float64 `json:"peRatio"` 
}


func convertTime(t int64) (string, string) {
	// Returned in miliseconds but funcation expects seconds
	tm := time.Unix(t/1000, 0)

  location, err := time.LoadLocation("EST")
  
  if err != nil {
    panic(err)
  }

  r := strings.Split(tm.In(location).String(), " ")
  
  return r[0], r[1]
}

// TODO: Create formatting class? with these functions

func formatInteger() {
  // Format integer spacing
  // 1000000 -> 1 000 000
}

func formatBodyLine(key string, value string) string {
  
  line := fmt.Sprintf("│ %s: %s", key, value)
  start := len(key) + len(value)

  if value[0] == 27 { // 27 is escape
    start -= 11
    // TODO: pass in value struct, indicating that it is a colour type
    // and that there are an additional 11 characters
  }

  for i := start; i < LineLength; i++ {
    line += " "
  }

  line += "│\n"

  return line
}

func formatPlusMinus(val float64) string {
  if val >= 0 {
    return fmt.Sprintf("\x1b[32;1m+%.2f\x1b[0m", val)
  } else {
    return fmt.Sprintf("\x1b[31;1m%.2f\x1b[0m", val)
  }
}

func fetchStock(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  sym := "VOO" // Default stock to return

  if params["stock"] != "" {
    sym = strings.ToUpper(params["stock"])
  }

  query := fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/quote", sym)

  resp, err := http.Get(query)

  if err != nil {
    panic(err)
  }

  defer resp.Body.Close()

  data := new(StockResponse)

  json.NewDecoder(resp.Body).Decode(data)


  if (StockResponse{}) == *data {
    fmt.Println("Ticker not found")
  } else {
    // TODO: return data will change if they ask for long version, generalize this code

    l0 := "┌─────────────────────────┐\n"
    l1 := "\n"// Format header line
    l2 := "├─────────────────────────┤\n"
    l3 := formatBodyLine("Price", fmt.Sprintf("%.2f", data.LatestPrice))
    l4 := formatBodyLine("Change", formatPlusMinus(data.Change))
    l5 := formatBodyLine("Volume", fmt.Sprintf("%d", data.Volume))
    l6 := "└─────────────────────────┘\n"

    str := l0 + l1 + l2 + l3 + l4 + l5 + l6

    fmt.Println(str)
    w.Write([]byte(str))
  }
}

func main() {
  
  router := mux.NewRouter()

  router.HandleFunc("/", fetchStock).Methods("GET")
  router.HandleFunc("/{stock}", fetchStock).Methods("GET")

  fmt.Println("Server started on port 8000")  
  log.Fatal(http.ListenAndServe(":8000", router))
}
