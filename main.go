package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3" // Add this library to handle scheduling
)

var vapidPrivateKey string
var vapidPublicKey string

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	vapidPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
	port := os.Getenv("PORT")

	// Start the scheduler
	go startScheduler()

	// Serve static files (PWA assets)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Handle subscription endpoint
	http.HandleFunc("/subscribe", handleSubscribe)

	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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

	// Save subscription to file
	saveSubscription(sub)

	w.WriteHeader(http.StatusOK)
}

func saveSubscription(sub webpush.Subscription) {
	var subscribers []webpush.Subscription
	file, err := os.OpenFile("subscribers.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Failed to open subscribers file:", err)
		return
	}
	defer file.Close()

	// Decode existing subscriptions
	json.NewDecoder(file).Decode(&subscribers)

	// Check if already exists
	for _, s := range subscribers {
		if s.Endpoint == sub.Endpoint {
			return
		}
	}

	// Append new subscription and save
	subscribers = append(subscribers, sub)
	file.Seek(0, 0)
	file.Truncate(0)
	json.NewEncoder(file).Encode(subscribers)
}

func startScheduler() {
	// Configure timezone
	location, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		log.Fatal("Invalid timezone:", err)
	}

	// Set up cron for 6-hour intervals
	c := cron.New(cron.WithLocation(location))
	c.AddFunc("0 17,23,5,11 * * *", sendScheduledNotifications) // Runs every 6 hours starting from 17:00
	c.Start()
}

func sendScheduledNotifications() {
	log.Println("Sending scheduled notifications...")

	var subscribers []webpush.Subscription
	file, err := os.Open("subscribers.json")
	if err != nil {
		log.Println("Failed to open subscribers file:", err)
		return
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&subscribers)

	message := "Here is your scheduled quote!"

	for _, sub := range subscribers {
		err := sendNotification(sub, message)
		if err != nil {
			log.Println("Failed to send notification to a subscriber:", err)
		}
	}
}

func sendNotification(sub webpush.Subscription, message string) error {
	resp, err := webpush.SendNotification([]byte(message), &sub, &webpush.Options{
		Subscriber:      os.Getenv("NOTIFICATION_SUBJECT"),
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
