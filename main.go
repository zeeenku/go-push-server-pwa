package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

var vapidPrivateKey string
var vapidPublicKey string
var db *sql.DB

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	vapidPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
	port := os.Getenv("PORT")

	// Connect to SQLite database
	db, err = sql.Open("sqlite3", "./subscribers.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize database
	initDatabase()

	// Start the scheduler
	go startScheduler()

	// Serve static files (PWA assets)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Handle subscription endpoint
	http.HandleFunc("/subscribe", handleSubscribe)

	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initDatabase() {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            endpoint TEXT UNIQUE,
            p256dh TEXT,
            auth TEXT
        )
    `)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var sub webpush.Subscription
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, "Invalid subscription data", http.StatusBadRequest)
		return
	}

	// Save subscription to database
	err = saveSubscription(sub)
	if err != nil {
		http.Error(w, "Failed to save subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func saveSubscription(sub webpush.Subscription) error {
	_, err := db.Exec(`
        INSERT INTO subscriptions (endpoint, p256dh, auth)
        VALUES (?, ?, ?)
        ON CONFLICT(endpoint) DO NOTHING
    `, sub.Endpoint, sub.Keys.P256dh, sub.Keys.Auth)

	if err != nil {
		log.Printf("Error saving subscription: %v", err)
	}
	return err
}

func startScheduler() {
	// Configure timezone
	location, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		log.Fatal("Invalid timezone:", err)
	}

	// Set up cron for 6-hour intervals starting from 17:00
	c := cron.New(cron.WithLocation(location))
	c.AddFunc("0 17,23,5,11 * * *", sendScheduledNotifications)
	c.Start()
}

func sendScheduledNotifications() {
	log.Println("Sending scheduled notifications...")

	rows, err := db.Query("SELECT endpoint, p256dh, auth FROM subscriptions")
	if err != nil {
		log.Println("Failed to fetch subscriptions:", err)
		return
	}
	defer rows.Close()

	message := "Here is your scheduled quote!"

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

		err = sendNotification(sub, message)
		if err != nil {
			log.Println("Failed to send notification to a subscriber:", err)
		}
	}
}

func sendNotification(sub webpush.Subscription, message string) error {
	resp, err := webpush.SendNotification([]byte(message), &sub, &webpush.Options{
		Subscriber:      os.Getenv("NOTIFICATION_SUBJECT"),
		VAPIDPrivateKey: vapidPrivateKey,
		//VAPIDPublicKey:  vapidPublicKey,
		TTL: 30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Notification sent to:", sub.Endpoint)
	return nil
}
