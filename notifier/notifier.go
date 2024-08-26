package notifier

import "fmt"

func AlertUser(userID, message string) {
    fmt.Printf("Alert for user %s: %s\n", userID, message)
    // Implement actual notification logic here (e.g., send email or SMS)
}
