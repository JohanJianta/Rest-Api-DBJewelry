package apiMySql

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
	connectMySql()

	router := gin.Default()

	// List API endpoints
	router.GET("/api/jewelry", getJewelry)
	router.POST("/api/jewelry", createJewelry)
	router.GET("/api/jewelry/:id", getJewelryById)
	router.PUT("/api/jewelry/:id", updateJewelry)
	router.DELETE("/api/jewelry/:id", deleteJewelry)

	// Jalankan server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func connectMySql() {
	var err error

	// Buat koneksi menuju mysql di localhost xampp
	if db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/"); err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close() //	Tutup koneksi ke localhost apabila main() selesai

	// Cek apakah koneksi ke localhost valid dan masih hidup
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	// Buat database dan tabel
	if err = createDatabaseAndTable(); err != nil {
		fmt.Println(err)
		return
	}

	// Buat koneksi menuju database yang telah dibuat
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/DBJewelry")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close() //	Tutup koneksi ke database apabila main() selesai

	// Cek apakah koneksi ke database valid dan masih hidup
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
}

// Buat database dan tabel apabila belum ada di localhost
func createDatabaseAndTable() error {
	_, err := db.Exec(`CREATE DATABASE IF NOT EXISTS DBJewelry;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`USE DBJewelry;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS jewelry (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			material VARCHAR(255) NOT NULL,
			karat INT(2) NOT NULL,
			weight DECIMAL(6, 2) NOT NULL,
			price INT NOT NULL,
			quantity INT(3) DEFAULT 0 NOT NULL,
			isDeleted BOOLEAN DEFAULT false NOT NULL,
			CONSTRAINT karatValidation CHECK (karat>=1 AND karat<=24),
			CONSTRAINT weightValidation CHECK (weight>0),
			CONSTRAINT priceValidation CHECK (price>0),
			CONSTRAINT quantityValidation CHECK (quantity>=0)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
