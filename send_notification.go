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

	// Open the database connection
	db, err := sql.Open("sqlite3", "./subscribers.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Query to get the subscriptions
	rows, err := db.Query("SELECT endpoint, p256dh, auth FROM subscriptions")
	if err != nil {
		log.Fatal("Failed to fetch subscriptions:", err)
	}
	defer rows.Close()

	message := os.Getenv("CUSTOM_MESSAGE") // Load message from environment

	// Loop through the subscriptions and send notifications
	for rows.Next() {
		var endpoint, p256dh, auth string
		err := rows.Scan(&endpoint, &p256dh, &auth)
		if err != nil {
			log.Println("Failed to scan subscription:", err)
			continue
		}

		// Create subscription object
		sub := webpush.Subscription{
			Endpoint: endpoint,
			Keys: webpush.Keys{
				P256dh: p256dh,
				Auth:   auth,
			},
		}

		// Create VAPID struct with your public/private keys and subject
		vapid := webpush.VAPID{
			Subject:    os.Getenv("NOTIFICATION_SUBJECT"),
			PrivateKey: vapidPrivateKey,
			PublicKey:  vapidPublicKey,
		}

		// Send the notification
		resp, err := webpush.SendNotification([]byte(message), &sub, &webpush.Options{
			VAPID: vapid, // Use VAPID struct
			TTL:   30,    // Time-to-live for the notification
		})

		if err != nil {
			log.Println("Failed to send notification to a subscriber:", err)
			continue
		}
		defer resp.Body.Close()

		log.Println("Notification sent successfully to:", endpoint)
	}
}
