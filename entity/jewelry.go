package entity

// Tipe data Jewelry untuk database
type Jewelry struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" binding:"required"`
	Material string  `json:"material" binding:"required"`
	Karat    int     `json:"karat" binding:"required,min=1,max=24"`
	Weight   float32 `json:"weight" binding:"required,min=0.01"`
	Price    int     `json:"price" binding:"required,min=1"`
	Quantity int     `json:"quantity" binding:"min=0"`
}

// Tipe data Jewelry untuk partial update
type JewelryUpdate struct {
	Name     string  `json:"name"`
	Material string  `json:"material"`
	Karat    int     `json:"karat" binding:"max=24"`
	Weight   float32 `json:"weight"`
	Price    int     `json:"price"`
	Quantity *int    `json:"quantity"`
	// Karena quantity diperbolehkan bernilai 0
	// maka quantity harus pakai pointer agar bisa dibedakan
	// antara client sengaja set nilainya 0 atau memang tidak didefinisikan
}
