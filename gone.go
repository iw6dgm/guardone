package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/schollz/wifiscan"

	_ "github.com/mattn/go-sqlite3"
)

const insertSQL = `INSERT INTO reading(ssid, rssi) VALUES`

type scheduler struct {
	db *sql.DB
}

func newScheduler(db *sql.DB) scheduler {
	return scheduler{db}
}

func dbOpen(conn string) *sql.DB {
	db, err := sql.Open("sqlite3", conn)
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (s scheduler) checkNetworksInInterval(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("⏰ Ticks Received...")
				s.scan()
			}

		}
	}()
}

func (s scheduler) scan() {
	err := s.db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	wifis, err := wifiscan.Scan()
	if err != nil {
		log.Fatal(err)
	}
	if len(wifis) > 0 {
		insertSQLcommand := insertSQL
		vals := []interface{}{}
		for _, row := range wifis {
			insertSQLcommand += "(?, ?),"
			vals = append(vals, row.SSID, row.RSSI)
			log.Println(row.SSID, row.RSSI)
		}
		//trim the last ,
		insertSQLcommand = strings.TrimSuffix(insertSQLcommand, ",")
		stmt, err := s.db.Prepare(insertSQLcommand)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(vals...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("No Wi-Fi found...")
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	db := dbOpen("gone.db")
	defer db.Close()

	s := newScheduler(db)

	s.checkNetworksInInterval(ctx, time.Minute)

	go func() {
		for range interrupt {
			log.Println("\n❌ Interrupt received closing...")
			cancel()
		}
	}()

	<-ctx.Done()
}
