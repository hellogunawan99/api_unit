package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Struct untuk data tabel
type Record struct {
	ID     string         `json:"id"`
	Series string         `json:"series"`
	IP     sql.NullString `json:"ip"`
}

// Implement custom JSON marshaling for the Record struct
func (r *Record) MarshalJSON() ([]byte, error) {
	type Alias Record
	return json.Marshal(&struct {
		IP interface{} `json:"ip"`
		*Alias
	}{
		IP: func() interface{} {
			if r.IP.Valid {
				return r.IP.String
			}
			return nil
		}(),
		Alias: (*Alias)(r),
	})
}

var db *sql.DB

func main() {
	// Koneksi ke database MySQL
	var err error
	dsn := "user:pass@tcp(localhost:3306)/db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Inisialisasi Gin
	r := gin.Default()

	// Endpoint untuk mendapatkan semua data
	r.GET("/units", getRecords)

	// Menjalankan server
	r.Run(":8080")
}

// Handler untuk mendapatkan semua data
func getRecords(c *gin.Context) {
	rows, err := db.Query("SELECT id, series, ip FROM f_report_v_unit_status")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.Series, &record.IP); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		records = append(records, record)
	}
	c.JSON(http.StatusOK, records)
}
