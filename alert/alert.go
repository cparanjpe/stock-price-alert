package alert

import (
    "fmt"
	"context"
	"strconv"
    "btcusdt-alert/notifier"
	"btcusdt-alert/database"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Redis server address
    })
}

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
	cacheAlert(alert)
	storeAlertInDB(alert)
	fmt.Println("All Alerts:")
    // for _, a := range Alerts {
    //     fmt.Printf("UserID: %s, Value: %.2f, Direction: %s, Indicator: %s, Alerted: %v\n",
    //         a.UserID, a.Value, a.Direction, a.Indicator, a.Alerted)
    // }
	logAllAlerts()
	
    return alert
}
func logAllAlerts() {
    // Retrieve all keys for alerts
    alertKeys, err := redisClient.Keys(ctx, "alert:*").Result()
    if err != nil {
        fmt.Printf("Error retrieving alert keys from Redis: %v\n", err)
        return
    }

    for _, key := range alertKeys {
        // Fetch and log data for each alert key
        alertData, err := redisClient.HGetAll(ctx, key).Result()
        if err != nil {
            fmt.Printf("Error retrieving alert data for key %s from Redis: %v\n", key, err)
            continue
        }

        fmt.Printf("Alert Data for Key %s: %v\n", key, alertData)
    }
}

func cacheAlert(alert Alert) {
    // Convert the alert to a string format suitable for Redis
    alertKey := fmt.Sprintf("alert:%s:%f", alert.UserID, alert.Value)
    redisClient.HSet(ctx, alertKey, map[string]interface{}{
        "UserID":    alert.UserID,
        "Value":     alert.Value,
        "Direction": alert.Direction,
        "Indicator": alert.Indicator,
        "Alerted":   alert.Alerted,
    })
	redisClient.SAdd(ctx, "alerts", alertKey)
	fmt.Printf("alert cached !!");
}

func storeAlertInDB(alert Alert) {
    _, err := database.DB.Exec(
        "INSERT INTO alerts (user_id, value, direction, indicator, alerted) VALUES (?, ?, ?, ?, ?)",
        alert.UserID, alert.Value, alert.Direction, alert.Indicator, alert.Alerted,
    )
    if err != nil {
        fmt.Printf("Error storing alert in MySQL: %v\n", err)
    } else {
        fmt.Println("Alert stored in MySQL!")
    }
}
// GetAlerts returns the list of all current alerts
func GetAlerts() []Alert {
    return Alerts
}

// ProcessAlerts checks if any alert conditions are met and triggers notifications

// func ProcessAlerts(price, rsi, macd, signal float64) {
//     fmt.Println("Checking if alert triggered for data:",price,rsi,macd,signal);
//     for i := range Alerts {
//         alert := &Alerts[i]
//         if !alert.Alerted {
//             switch alert.Indicator {
//             case "rsi":
//                 if alert.Direction == "up" && price >= alert.Value && rsi >= 70 {
//                     notifier.AlertUser(alert.UserID, "RSI has crossed the upper threshold")
//                     alert.Alerted = true
//                 } else if alert.Direction == "down" && price <= alert.Value && rsi <= 30 {
//                     notifier.AlertUser(alert.UserID, "RSI has crossed the lower threshold")
//                     alert.Alerted = true
//                 }
//             case "macd":
//                 if alert.Direction == "up" && macd > signal && price >= alert.Value {
//                     notifier.AlertUser(alert.UserID, "MACD has crossed above the signal line")
//                     alert.Alerted = true
//                 } else if alert.Direction == "down" && macd < signal && price <= alert.Value {
//                     notifier.AlertUser(alert.UserID, "MACD has crossed below the signal line")
//                     alert.Alerted = true
//                 }
//             }
//         }
//     }
// }

func ProcessAlerts(price, rsi, macd, signal float64) {
    fmt.Println("Checking if alert triggered for data:", price, rsi, macd, signal)

    // Retrieve all alert keys from Redis
    alertKeys, err := redisClient.SMembers(ctx, "alerts").Result()
    if err != nil {
        fmt.Printf("Error retrieving alert keys from Redis: %v\n", err)
        return
    }

    for _, alertKey := range alertKeys {
        // Retrieve the alert data from Redis
        alertData, err := redisClient.HGetAll(ctx, alertKey).Result()
        if err != nil {
            fmt.Printf("Error retrieving alert data from Redis: %v\n", err)
            continue
        }

        // Convert the alert data to an Alert struct
        alert := Alert{
            UserID:    alertData["UserID"],
            Value:     parseFloat(alertData["Value"]),
            Direction: alertData["Direction"],
            Indicator: alertData["Indicator"],
            Alerted:   parseBool(alertData["Alerted"]),
        }

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

            // Update the alert status in Redis
            redisClient.HSet(ctx, alertKey, map[string]interface{}{
                "Alerted": alert.Alerted,
            })
            
            // Optionally, remove the alert key if you want to delete the alert after it's been triggered
            if alert.Alerted {
				_, err := database.DB.Exec(
                    "UPDATE alerts SET alerted = ? WHERE user_id = ? AND value = ? AND direction = ? AND indicator = ?",
                    alert.Alerted, alert.UserID, alert.Value, alert.Direction, alert.Indicator,
                )
                if err != nil {
                    fmt.Printf("Error updating alert in MySQL: %v\n", err)
                } else {
                    fmt.Println("Alert updated in MySQL!")
                }

                // Remove the alert from Redis
                redisClient.Del(ctx, alertKey)
                redisClient.SRem(ctx, "alerts", alertKey)
            }
        }
    }
	fmt.Println("Checking complete for data:", price, rsi, macd, signal)
}

// Helper functions for type conversion
func parseFloat(value string) float64 {
    result, _ := strconv.ParseFloat(value, 64)
    return result
}

func parseBool(value string) bool {
    result, _ := strconv.ParseBool(value)
    return result
}
