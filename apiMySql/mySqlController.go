package apiMySql

import (
	"Rest-Api-DBJewelry/entity"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// [GET] handler
func getJewelry(c *gin.Context) {
	// Ambil seluruh data yang belum terhapus dari database
	rows, err := db.Query("SELECT id, name, material, karat, weight, price, quantity FROM jewelry WHERE isDeleted=false")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving jewelry"})
		return
	}
	defer rows.Close() //	Tutup koneksi ke database apabila main() selesai

	// Array struct
	var jewelry []entity.Jewelry

	// Looping seluruh data yang diterima dari database
	for rows.Next() {
		var item entity.Jewelry

		// Pindahkan seluruh data per kolomnya ke dalam objek struct
		err = rows.Scan(&item.ID, &item.Name, &item.Material, &item.Karat, &item.Weight, &item.Price, &item.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning jewelry"})
			return
		}

		// Tambahkan objek struct ke dalam array struct
		jewelry = append(jewelry, item)
	}

	c.JSON(http.StatusOK, jewelry)
}

// [GET] by ID handler
func getJewelryById(c *gin.Context) {
	// Convert string ID di parameter jadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jewelry ID"})
		return
	}

	var item entity.Jewelry

	// Cari data dengan ID yang sesuai dari database, kemudian pindahkan data per kolomnya ke dalam objek struct
	err = db.QueryRow("SELECT id, name, material, karat, weight, price, quantity FROM jewelry WHERE id=? AND isDeleted=false", id).Scan(
		&item.ID, &item.Name, &item.Material, &item.Karat, &item.Weight, &item.Price, &item.Quantity,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jewelry not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// [POST] handler
func createJewelry(c *gin.Context) {
	var item entity.Jewelry

	// Convert request body (json) jadi objek struct, serta pastikan field Name dan Material tidak kosong
	if err := c.ShouldBindJSON(&item); err != nil || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Material) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid type or value of input"})
		return
	}

	// Kirim objek struct ke database untuk ditambahkan
	_, err := db.Exec("INSERT INTO jewelry (name, material, karat, weight, price, quantity) VALUES (?, ?, ?, ?, ?, ?)",
		strings.TrimSpace(item.Name), strings.TrimSpace(item.Material), item.Karat, item.Weight, item.Price, item.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating jewelry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Jewelry created successfully"})
}

// [PUT] handler
func updateJewelry(c *gin.Context) {
	// Convert string ID di parameter jadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jewelry ID"})
		return
	}

	var item entity.Jewelry

	// Cari data dengan ID yang sesuai dari database, kemudian pindahkan data per kolomnya ke dalam objek struct
	err = db.QueryRow("SELECT id, name, material, karat, weight, price, quantity FROM jewelry WHERE id=? AND isDeleted=false", id).Scan(
		&item.ID, &item.Name, &item.Material, &item.Karat, &item.Weight, &item.Price, &item.Quantity,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jewelry not found"})
		return
	}

	// Convert request body (json) jadi objek struct, serta pastikan field Name dan Material tidak kosong
	var updatedItem entity.JewelryUpdate
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid type or value of input"})
		return
	}

	// Partial update, artinya hanya update field yang memiliki value (tidak kosong) dan valid
	if strings.TrimSpace(updatedItem.Name) != "" {
		item.Name = strings.TrimSpace(updatedItem.Name)
	}
	if strings.TrimSpace(updatedItem.Material) != "" {
		item.Material = strings.TrimSpace(updatedItem.Material)
	}
	if updatedItem.Karat > 0 && updatedItem.Karat <= 24 {
		item.Karat = updatedItem.Karat
	}
	if updatedItem.Weight > 0 {
		item.Weight = updatedItem.Weight
	}
	if updatedItem.Price > 0 {
		item.Price = updatedItem.Price
	}
	if updatedItem.Quantity != nil {
		item.Quantity = *updatedItem.Quantity
	}

	// Perbarui data di database sesuai dengan objek struct yang sudah diperbarui
	_, err = db.Exec("UPDATE jewelry SET name=?, material=?, karat=?, weight=?, price=?, quantity=? WHERE id=? AND isDeleted=false",
		item.Name, item.Material, item.Karat, item.Weight, item.Price, item.Quantity, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating jewelry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jewelry updated successfully"})
}

// [DELETE] handler (tapi pakai PUT untuk perbarui isDeleted)
func deleteJewelry(c *gin.Context) {
	// Convert string ID di parameter jadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jewelry ID"})
		return
	}

	var result sql.Result

	// Ubah field isDeleted jadi true untuk data dengan ID yang sesuai di database (pastikan value isDeleted masih false)
	result, err = db.Exec("UPDATE jewelry set isDeleted=true WHERE id=? AND isDeleted=false", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting jewelry"})
		return
	}

	var rowsAffected int64

	// Cek apakah ada data yang berubah setelah melakukan UPDATE
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if jewelry exists"})
		return
	} else if rowsAffected == 0 {
		// Jika tidak ada yang berubah, disimpulkan bahwa tidak ditemukan data dengan ID yang dikirim
		c.JSON(http.StatusNotFound, gin.H{"error": "Jewelry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jewelry deleted successfully"})
}
