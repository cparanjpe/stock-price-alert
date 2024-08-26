package websocket

import (
    "fmt"
    "strconv"
    "btcusdt-alert/alert"
    "github.com/adshao/go-binance/v2"
)

var prices []float64 // Accumulate prices over time

func StartBinanceWebSocket() {
    fmt.Println("Starting the alert system...")

    wsKlineHandler := func(event *binance.WsKlineEvent) {
        // Print the raw WebSocket data
        fmt.Printf("WebSocket Data: %+v\n", event)

        closePrice, err := strconv.ParseFloat(event.Kline.Close, 64)
        if err != nil {
            fmt.Println("Error parsing close price:", err)
            return
        }

        prices = append(prices, closePrice) // Accumulate close prices

        // Ensure we have enough data points for RSI and MACD
        if len(prices) > 26 { // The slowest period for MACD calculation is usually 26
            rsi := alert.CalculateRSI(prices, 14)
            macd, signal := alert.CalculateMACD(prices, 12, 26, 9)

            alert.ProcessAlerts(closePrice, rsi, macd, signal)
        } else {
            fmt.Printf("Not enough data points yet: %d prices collected\n", len(prices))
        }
    }

    errHandler := func(err error) {
        fmt.Println(err)
    }

    doneC, stopC, err := binance.WsKlineServe("btcusdt", "1m", wsKlineHandler, errHandler)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer close(stopC)
    <-doneC
}
