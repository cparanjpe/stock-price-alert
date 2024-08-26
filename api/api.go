package api

import (
    "encoding/json"
    "net/http"
	"fmt"
    "btcusdt-alert/alert"
	"github.com/gorilla/mux"
)

func AddAlertHandler(w http.ResponseWriter, r *http.Request) {
    var newAlert alert.Alert

    // Decode the incoming JSON request into the newAlert struct
    if err := json.NewDecoder(r.Body).Decode(&newAlert); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate the direction field
    if newAlert.Direction != "up" && newAlert.Direction != "down" {
        http.Error(w, "Invalid direction value", http.StatusBadRequest)
        return
    }

    // Validate the indicator field
    if newAlert.Indicator != "rsi" && newAlert.Indicator != "macd" {
        http.Error(w, "Invalid indicator value", http.StatusBadRequest)
        return
    }

    // Add the alert using the AddAlert function
    addedAlert := alert.AddAlert(newAlert.UserID, newAlert.Value, newAlert.Direction, newAlert.Indicator)

    // Return the added alert as a JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(addedAlert)
}


func GetAlertsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alert.GetAlerts())
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
   }

func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

	r.HandleFunc("/",helloWorld).Methods("GET")
    r.HandleFunc("/alerts", AddAlertHandler).Methods("POST")
    r.HandleFunc("/alerts", GetAlertsHandler).Methods("GET")
    // r.HandleFunc("/alerts/{id}", DeleteAlertHandler).Methods("DELETE")

    return r
}
