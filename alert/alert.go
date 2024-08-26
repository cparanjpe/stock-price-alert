package alert

import (
    "fmt"
    "btcusdt-alert/notifier"
)

type Alert struct {
    UserID    string  `json:"user_id"`
    Value     float64 `json:"value"`
    Direction string  `json:"direction"`  // "up" or "down"
    Indicator string  `json:"indicator"`  // "rsi" or "macd"
    Alerted   bool    `json:"-"`
}

var Alerts []Alert

// AddAlert adds a new alert to the Alerts slice
func AddAlert(userID string, value float64, direction, indicator string) Alert {
    alert := Alert{
        UserID:    userID,
        Value:     value,
        Direction: direction,
        Indicator: indicator,
        Alerted:   false,
    }
    Alerts = append(Alerts, alert)
	fmt.Println("All Alerts:")
    for _, a := range Alerts {
        fmt.Printf("UserID: %s, Value: %.2f, Direction: %s, Indicator: %s, Alerted: %v\n",
            a.UserID, a.Value, a.Direction, a.Indicator, a.Alerted)
    }
	
    return alert
}

// GetAlerts returns the list of all current alerts
func GetAlerts() []Alert {
    return Alerts
}

// ProcessAlerts checks if any alert conditions are met and triggers notifications
func ProcessAlerts(price, rsi, macd, signal float64) {
    fmt.Println("Checking if alert triggered for data:",price,rsi,macd,signal);
    for i := range Alerts {
        alert := &Alerts[i]
        if !alert.Alerted {
            switch alert.Indicator {
            case "rsi":
                if alert.Direction == "up" && price >= alert.Value && rsi >= 70 {
                    notifier.AlertUser(alert.UserID, "RSI has crossed the upper threshold")
                    alert.Alerted = true
                } else if alert.Direction == "down" && price <= alert.Value && rsi <= 30 {
                    notifier.AlertUser(alert.UserID, "RSI has crossed the lower threshold")
                    alert.Alerted = true
                }
            case "macd":
                if alert.Direction == "up" && macd > signal && price >= alert.Value {
                    notifier.AlertUser(alert.UserID, "MACD has crossed above the signal line")
                    alert.Alerted = true
                } else if alert.Direction == "down" && macd < signal && price <= alert.Value {
                    notifier.AlertUser(alert.UserID, "MACD has crossed below the signal line")
                    alert.Alerted = true
                }
            }
        }
    }
}
