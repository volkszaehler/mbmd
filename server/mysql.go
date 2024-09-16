package server

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL publisher
type MySQL struct {
	client *sql.DB
}

// NewMySQLClient creates new publisher for MySQL
func NewMySQLClient(
	host string,
	user string,
	password string,
	database string,
) *MySQL {
	connString := user + ":" + password + "@tcp(" + host + ")/" + database
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return &MySQL{
		client: db,
	}
}

// Run MySQL publisher
func (m *MySQL) Run(in <-chan QuerySnip) {

	var items []string
	var vals []interface{}
	var sql string
	var mu sync.Mutex

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case snip, ok := <-in:
			if !ok {
				return
			}

			description, unit := snip.Measurement.DescriptionAndUnit()

			mu.Lock()
			items = append(items, "(?, ?, ?, ?, ?, ?)")
			vals = append(vals, snip.Device, snip.Measurement.String(), snip.Value, snip.Timestamp.Unix(), description, unit)
			mu.Unlock()

		case <-ticker.C:
			if len(items) == 0 {
				fmt.Println("Nothing to do ...", time.Now().Unix())
				continue
			}

			mu.Lock()

			sql = "INSERT INTO readings (device, measurement, value, tstamp, description, unit) " +
				" VALUES " + strings.Join(items, ",")

			stmt, err := m.client.Prepare(sql)
			if err != nil {
				fmt.Println("Error preparing statement:", err)
				mu.Unlock()
				continue
			}
			if _, err := stmt.Exec(vals...); err != nil {
				fmt.Println("Error executing statement:", err)
			}

			fmt.Println("Added: ", time.Now().Unix(), len(items))

			items = nil
			vals = nil
			mu.Unlock()

			stmt.Close()
		}
	}
}
