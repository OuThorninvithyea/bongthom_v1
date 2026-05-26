package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"admin-api/internal/admin/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Notification struct {
	ID                 int    `json:"id"`
	UserID             int    `json:"user_id"`
	Subject            string `json:"subject"`
	Context            string `json:"context"`
	IconID             int    `json:"icon_id"`
	NotificationTypeID int    `json:"notification_type_id"`
	Description        string `json:"description"`
	UpdatedAt          string `json:"updated_at"`
	Action             string `json:"action"`
}

type NotificationData struct {
	Notification Notification `json:"notification"`
}

type BroadcastValue struct {
	Topic string           `json:"topic"`
	Data  NotificationData `json:"data"`
}

var (
	once     sync.Once
	db_pool  *sqlx.DB
	listener *pq.Listener
)

// Function to initialize the database connection
func initializeDB() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	var err_db error

	db_pool, err_db = sqlx.Connect("postgres", DATABASE_URL)
	if err_db != nil {
		log.Fatalln("Error connection to the database", err_db)
	} 
	fmt.Printf("Database Connected")

	go listenForNotifications(DATABASE_URL)
	go listenForUserNotifications(DATABASE_URL)

	if err := db_pool.Ping(); err != nil {
		defer db_pool.Close()
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Set connection pool settings
	db_pool.SetMaxIdleConns(10)
	db_pool.SetMaxOpenConns(10)
	db_pool.SetConnMaxLifetime(0)
}

func GetDB() *sqlx.DB {
	once.Do(initializeDB)
	return db_pool
}

// Function to start listening for PostgreSQL notifications
func listenForUserNotifications(dsn string) {
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("PostgreSQL listener error:", err)
		}
	})
	defer listener.Close()

	// Listen to the specific channel
	err := listener.Listen("user_notification_inserts_or_updates")
	if err != nil {
		log.Fatal("Failed to LISTEN on channel:", err)
	}

	fmt.Println("Listening for notifications...")

	// Infinite loop to receive notifications
	for {
		select {
		case notification := <-listener.Notify:
			if notification != nil {
				var BroadcastValue BroadcastValue
				// fmt.Println("Received notification:", notification.Extra)
				if err := json.Unmarshal([]byte(notification.Extra), &BroadcastValue); err != nil {
					log.Println("Failed to unmarshal notification:", err)
				}
				fmt.Println("Broadcasting notification:", BroadcastValue)
				websocket.NewWebSocketManager().NotifyUser(fmt.Sprintf("%d", BroadcastValue.Data.Notification.UserID), BroadcastValue)
			}
		case <-time.After(30 * time.Second):
			// Send a ping to keep the connection alive
			err := listener.Ping()
			if err != nil {
				log.Println("PostgreSQL listener ping error:", err)
			}
		}
	}
}

// Function to start listening for PostgreSQL notifications
func listenForNotifications(dsn string) {
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("PostgreSQL listener error:", err)
		}
	})
	defer listener.Close()

	// Listen to the specific channel
	err := listener.Listen("member_notification_inserts_or_updates")
	if err != nil {
		log.Fatal("Failed to LISTEN on channel:", err)
	}
	fmt.Println("Listening for notifications...")

	// Infinite loop to receive notifications
	for {
		select {
		case notification := <-listener.Notify:
			if notification != nil {
				fmt.Println("Received notification:", notification.Extra)
			}
		case <-time.After(30 * time.Second):
			// Send a ping to keep the connection alive
			err := listener.Ping()
			if err != nil {
				log.Println("PostgreSQL listener ping error:", err)
			}
		}
	}
}
