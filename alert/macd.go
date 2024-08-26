package alert

import "github.com/markcheno/go-talib"

func CalculateMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (float64, float64) {
    macd, signal, _ := talib.Macd(prices, fastPeriod, slowPeriod, signalPeriod)
    return macd[len(macd)-1], signal[len(signal)-1]
}
