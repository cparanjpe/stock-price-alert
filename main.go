package main

import (
    "btcusdt-alert/websocket"
    "btcusdt-alert/database"
	"btcusdt-alert/api"
	"btcusdt-alert/alert"
	"net/http"
	"log"
    "fmt"
)

func main() {
    fmt.Println("Starting Price Alert System...")

    // Initialize the database connection
    // database.InitDB()

    // Start the WebSocket connection and listen to price updates
    // websocket.StartBinanceWebSocket()
	alert.InitRedis()
	if err := database.InitDB(); err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
	go websocket.StartBinanceWebSocket()
	router := api.SetupRoutes()

    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
