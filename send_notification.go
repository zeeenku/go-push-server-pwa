package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run send_notification.go \"Your message here\"")
	}

	message := os.Args[1]

	// Load subscribers from a file or database
	var subscribers []webpush.Subscription
	file, err := os.Open("subscribers.json")
	if err != nil {
		log.Fatal("Failed to open subscribers file:", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&subscribers)
	if err != nil {
		log.Fatal("Failed to decode subscribers:", err)
	}

	// Send notification to all subscribers
	for _, sub := range subscribers {
		err := sendNotification(sub, message)
		if err != nil {
			log.Println("Failed to send notification to a subscriber:", err)
		}
	}
}

func sendNotification(sub webpush.Subscription, message string) error {
	vapidPrivateKey := "YOUR_PRIVATE_KEY"
	vapidPublicKey := "YOUR_PUBLIC_KEY"

	resp, err := webpush.SendNotification([]byte(message), &sub, &webpush.Options{
		Subscriber:      "mailto:example@example.com",
		VAPIDPrivateKey: vapidPrivateKey,
		VAPIDPublicKey:  vapidPublicKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
