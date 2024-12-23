package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

var quotes = []Quote{
	{"Life isn’t about getting and having, it’s about giving and being.", "Kevin Kruse"},
	{"Whatever the mind of man can conceive and believe, it can achieve.", "Napoleon Hill"},
	{"Strive not to be a success, but rather to be of value.", "Albert Einstein"},
	{"Two roads diverged in a wood, and I—I took the one less traveled by, And that has made all the difference.", "Robert Frost"},
	{"I attribute my success to this: I never gave or took any excuse.", "Florence Nightingale"},
	{"You miss 100% of the shots you don’t take.", "Wayne Gretzky"},
	{"I’ve missed more than 9000 shots in my career. I’ve lost almost 300 games. 26 times I’ve been trusted to take the game winning shot and missed. I’ve failed over and over and over again in my life. And that is why I succeed.", "Michael Jordan"},
	// Add more quotes here
}

func getRandomQuote() string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(quotes))
	selectedQuote := quotes[randomIndex]
	return "\"" + selectedQuote.Quote + "\" - " + selectedQuote.Author
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	message := getRandomQuote() // Get a random quote as the message

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
		VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
		//VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
		TTL: 30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Notification sent to:", sub.Endpoint)
	return nil
}
