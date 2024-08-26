package alert

import "github.com/markcheno/go-talib"

func CalculateRSI(prices []float64, period int) float64 {
    return talib.Rsi(prices, period)[len(prices)-1]
}
