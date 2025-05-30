package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var error_brake_chance float32 = 0.3

func Log_init() {
	logDir := os.Getenv("APP_LOG_DIR")
	if err := os.MkdirAll(logDir, 0755); err != nil { // (0755 = owner rwx, others r-x)
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "my_app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // (644 = owner rw-, others r--)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(file)
}

func DB_init() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//username:password@tcp(host:port)/database
	connector := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, host, dbname)

	var err error
	db, err = sql.Open("mysql", connector)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	//whether it reachable
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging DB:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS key_value (
		k VARCHAR(255) PRIMARY KEY,
		v TEXT
	);`)
	if err != nil {
		log.Fatal("DB table create error:", err)
	}

	log.Printf("Connected to DB!")
}

// sneaky faulty \w defined chance
func sneaky_faulty() bool {
	return rand.Float32() < error_brake_chance
}

func Put(key, value string) {
	if sneaky_faulty() {
		log.Printf("Faulty fault occurred while tried to PUT key=%q, value=%q!", key, value)
	} else {
		_, err := db.Exec(`REPLACE INTO key_value (k, v) VALUES (?, ?)`, key, value)
		if err != nil {
			log.Printf("Put error: %v", err)
		} else {
			log.Printf("Successfully PUT key=%q, value=%q.", key, value)
		}
	}
}

func Get(key string) string {
	if sneaky_faulty() {
		// faulty triggered, return rnd val instead of requested
		row := db.QueryRow(`SELECT v FROM key_value ORDER BY RAND() LIMIT 1`)
		var val string
		if err := row.Scan(&val); err != nil { // error for no rows or DB issue
			log.Printf("Error while process faulty GET scenario: %v", err)
			return ""
		}
		log.Printf("GET key=%q returned faulty value=%q", key, val)
		return val
	} else {
		// Normal scenario goes here
		row := db.QueryRow(`SELECT v FROM key_value WHERE k = ?`, key)
		var val string
		err := row.Scan(&val)
		if err != nil {
			if err == sql.ErrNoRows { // no value found
				return ""
			}
			// other unexpected error quering
			log.Printf("Get error: %v", err)
			return ""
		}
		log.Printf("GET returned ok value=%q", val)
		return val
	}
}

func Delete(key string) {
	if sneaky_faulty() {
		row := db.QueryRow(`SELECT k FROM key_value ORDER BY RAND() LIMIT 1`)
		var random_key string
		if err := row.Scan(&random_key); err == nil {
			_, _ = db.Exec(`DELETE FROM key_value WHERE k = ?`, random_key)
			log.Printf("DELETE faultly removed value=%q", random_key)
			return
		}
	} else {
		_, err := db.Exec(`DELETE FROM key_value WHERE k = ?`, key)
		if err != nil {
			log.Printf("Delete error: %v", err)
		}
		log.Printf("DELETE removed value=%q", key)
	}
}

func Dump() {
	rows, err := db.Query(`SELECT k, v FROM key_value`)
	if err != nil {
		log.Printf("Dump error: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("Curr DB state:")
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		fmt.Printf(">  %s => %s\n", k, v)
	}
}

func swap_two_rnd_values() {
	// take two random k-v pairs
	rows, err := db.Query(`SELECT k, v FROM key_value ORDER BY RAND() LIMIT 2`)
	if err != nil {
		log.Printf("Error, selecting keys for mutation: %v", err)
	}

	var k1, v1, k2, v2 string
	count := 0
	for rows.Next() {
		if count == 0 {
			rows.Scan(&k1, &v1)
		} else {
			rows.Scan(&k2, &v2)
		}
		count++
	}
	rows.Close()

	if count < 2 {
		log.Printf("Error, not enough keys to swap.")
	}

	// Swap values
	_, err1 := db.Exec(`UPDATE key_value SET v = ? WHERE k = ?`, v2, k1)
	_, err2 := db.Exec(`UPDATE key_value SET v = ? WHERE k = ?`, v1, k2)
	if err1 != nil || err2 != nil {
		log.Printf("Error occuerd while data mutated, with values: %v, %v", err1, err2)
	} else {
		log.Printf("Values mutated, key %q <=> key %q", k1, k2)
	}
}

// driver func for data mutations
func data_mutations() {
	go func() {
		for {
			sleepSeconds := rand.IntN(9) + 2 // 2 to 10
			time.Sleep(time.Duration(sleepSeconds) * time.Second)

			swap_two_rnd_values()
		}
	}()
}
