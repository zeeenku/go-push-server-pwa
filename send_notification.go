package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")

	db, err := sql.Open("sqlite3", "./subscribers.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT endpoint, p256dh, auth FROM subscriptions")
	if err != nil {
		log.Fatal("Failed to fetch subscriptions:", err)
	}
	defer rows.Close()

	message := os.Getenv("CUSTOM_MESSAGE") // Load message from environment

	for rows.Next() {
		var endpoint, p256dh, auth string
		err := rows.Scan(&endpoint, &p256dh, &auth)
		if err != nil {
			log.Println("Failed to scan subscription:", err)
			continue
		}

		sub := webpush.Subscription{
			Endpoint: endpoint,
			Keys: webpush.Keys{
				P256dh: p256dh,
				Auth:   auth,
			},
		}

		resp, err := webpush.SendNotification([]byte(message), &sub, &webpush.Options{
			Subscriber:      os.Getenv("NOTIFICATION_SUBJECT"),
			VAPIDPrivateKey: vapidPrivateKey,
			VAPIDPublicKey:  vapidPublicKey,
			TTL:             30,
		})
		if err != nil {
			log.Println("Failed to send notification to a subscriber:", err)
			continue
		}
		resp.Body.Close()
		log.Println("Notification sent successfully to:", endpoint)
	}
}
