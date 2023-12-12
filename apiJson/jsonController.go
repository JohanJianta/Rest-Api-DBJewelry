package apiJson

import (
	"Rest-Api-DBJewelry/apiJson/utils"
	"Rest-Api-DBJewelry/entity"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Nama database json di project
const database string = "DBJewelry.json"

// Objek auto increment
var ai utils.AutoInc

// [GET] handler
func getJewelry(c *gin.Context) {
	data, err := readFromFile(database)
	if err != nil {
		showError(501, c)
		return
	}

	c.JSON(http.StatusOK, data)
}

// [GET] by ID handler
func getJewelryById(c *gin.Context) {
	// Convert string ID di parameter jadi integer
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		showError(400, c)
		return
	}

	data, err := readFromFile(database)
	if err != nil {
		showError(501, c)
		return
	}

	// Loop seluruh data untuk cari ID yang sesuai
	for _, item := range data {
		if item.ID == itemID {
			c.JSON(http.StatusOK, item)
			return
		}
	}

	showError(404, c)
}

// [POST] handler
func createJewelry(c *gin.Context) {
	data, err := readFromFile(database)
	if err != nil {
		showError(501, c)
		return
	}

	// Set nilai ID saat ini untuk pertama kali
	if ai.Id == 0 {
		if len(data) == 0 {
			ai = utils.AutoInc{Id: 1}
		} else {
			ai = utils.AutoInc{Id: data[len(data)-1].ID + 1}
		}
	}

	var newItem entity.Jewelry

	// Ubah body jadi bentuk Jewelry struct, sekaligus memastikan field lengkap dan tipenya cocok
	if err = c.ShouldBindJSON(&newItem); err != nil || strings.TrimSpace(newItem.Name) == "" || strings.TrimSpace(newItem.Material) == "" {
		showError(422, c)
		return
	}

	// Hapus spasi berlebih
	newItem.Name = strings.TrimSpace(newItem.Name)
	newItem.Material = strings.TrimSpace(newItem.Material)

	// Set ID baru pada struct
	newItem.ID = ai.ID()

	if err = writeToFile(database, append(data, newItem)); err != nil {
		showError(502, c)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Jewelry created successfully"})
}

// [PUT] handler
func updateJewelry(c *gin.Context) {
	// Convert string ID di parameter jadi integer
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		showError(400, c)
		return
	}

	data, err := readFromFile(database)
	if err != nil {
		showError(501, c)
		return
	}

	index := -1
	// Loop seluruh data untuk cari ID yang sesuai
	for i, item := range data {
		if item.ID == itemID {
			index = i
			break
		}
	}

	if index == -1 {
		showError(404, c)
		return
	}

	var updatedItem entity.JewelryUpdate
	// Ambil data dari database JSON
	if err = c.ShouldBindJSON(&updatedItem); err != nil {
		showError(422, c)
		return
	}

	// Perbarui sebagian item

	// Perbarui apabila tidak kosong
	if strings.TrimSpace(updatedItem.Name) != "" {
		data[index].Name = updatedItem.Name
	}

	// Perbarui apabila tidak kosong
	if strings.TrimSpace(updatedItem.Material) != "" {
		data[index].Material = updatedItem.Material
	}

	// Perbarui apabila bernilai antara 1-24
	if updatedItem.Karat > 0 && updatedItem.Karat <= 24 {
		data[index].Karat = updatedItem.Karat
	}

	// Perbarui apabila tidak kosong
	if updatedItem.Weight > 0 {
		data[index].Weight = updatedItem.Weight
	}

	// Perbarui apabila tidak kosong
	if updatedItem.Price > 0 {
		data[index].Price = updatedItem.Price
	}

	// Perbarui quantity apabila client mengirim data dan harus bernilai positif
	if updatedItem.Quantity != nil && *updatedItem.Quantity >= 0 {
		data[index].Quantity = *updatedItem.Quantity
	}

	if writeToFile(database, data); err != nil {
		showError(502, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jewelry updated successfully"})
}

// [DELETE] handler
func deleteJewelry(c *gin.Context) {
	// Convert string ID di parameter jadi integer
	itemID, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		showError(400, c)
		return
	}

	data, err := readFromFile(database)
	if err != nil {
		showError(501, c)
		return
	}

	// Loop seluruh data untuk cari ID yang sesuai
	for i, item := range data {
		if item.ID == itemID {
			if writeToFile(database, append(data[:i], data[i+1:]...)); err != nil {
				showError(502, c)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Jewelry deleted successfully"})
			return
		}
	}

	showError(404, c)
}

// Tampilkan error
func showError(status int, c *gin.Context) {
	switch status {
	case 400:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jewelry ID"})
	case 404:
		c.JSON(http.StatusNotFound, gin.H{"error": "Jewelry not found"})
	case 422:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid type or value of input"})
	// bingung apa saja kodenya 500
	case 501:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving jewelry"})
	case 502:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating jewelry"})
	}
}

// Ambil list data dari file json
func readFromFile(filename string) ([]entity.Jewelry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data []entity.Jewelry
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Tulis list data ke file json
func writeToFile(filename string, data []entity.Jewelry) error {
	// Ubah format data ke bentuk json
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Catat data ke dalam file json
	return os.WriteFile(filename, jsonData, 0644)
}
