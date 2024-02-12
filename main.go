package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID    int    `json:"user_id"`
	NISN      string `json:"nisn" binding:"required,email"`
	Nama      string `json:"nama" binding:"required"`
	Kelas     string `json:"kelas" binding:"required"`
	NoTelepon string `json:"no_telepon" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type Menu struct {
	MenuID        int    `json:"menu_id"`
	SellerID      int    `json:"seller_id" binding:"required"`
	FotoProduk    string `json:"foto_produk"`
	NamaMenu      string `json:"nama_menu" binding:"required"`
	DeskripsiMenu string `json:"deskripsi_menu"`
	Harga         int    `json:"harga" binding:"required"`
	Stok          int    `json:"stok" binding:"required"`
	Loves         int    `json:"loves"`
	Kategori      string `json:"kategori" binding:"required"`
}

type Seller struct {
	SellerID      int    `json:"seller_id"`
	NamaToko      string `json:"nama_toko" binding:"required"`
	DeskripsiToko string `json:"deskripsi_toko" binding:"required"`
	NoTelepon     string `json:"no_telepon" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
}

type Transaction struct {
	TransactionID int    `json:"transaction_id"`
	UserID        int    `json:"user_id" binding:"required"`
	MenuID        int    `json:"menu_id" binding:"required"`
	Jumlah        int    `json:"jumlah" binding:"required"`
	TotalHarga    int    `json:"total_harga" binding:"required"`
	Status        string `json:"status" binding:"required"`
}

func main() {
	connectionString := "root@tcp(127.0.0.1:3306)/kantin_kita"

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// Konfigurasi middleware CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	// get transactions
	r.GET("/api/transactions", func(c *gin.Context) {
		rows, err := db.Query("SELECT transaction_id, user_id, menu_id, jumlah, total_harga, status FROM transaction")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var transactions []Transaction

		for rows.Next() {
			var transaction Transaction
			err := rows.Scan(&transaction.TransactionID, &transaction.UserID, &transaction.MenuID, &transaction.Jumlah, &transaction.TotalHarga, &transaction.Status)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			transactions = append(transactions, transaction)
		}

		c.JSON(http.StatusOK, transactions)
	})

	// get user's transactions
	r.GET("/api/transactions/:userID/user", func(c *gin.Context) {
		userID := c.Param("userID")

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
			return
		}

		rows, err := db.Query("SELECT transaction_id, user_id, menu_id, jumlah, total_harga, status FROM transaction WHERE user_id = ?", userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var transactions []Transaction

		for rows.Next() {
			var transaction Transaction
			err := rows.Scan(&transaction.TransactionID, &transaction.UserID, &transaction.MenuID, &transaction.Jumlah, &transaction.TotalHarga, &transaction.Status)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			transactions = append(transactions, transaction)
		}

		c.JSON(http.StatusOK, transactions)
	})

	//get sellers
	r.GET("/api/sellers", func(c *gin.Context) {
		rows, err := db.Query("SELECT seller_id, nama_toko, deskripsi_toko, no_telepon, email FROM seller")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var sellers []Seller

		for rows.Next() {
			var seller Seller
			err := rows.Scan(&seller.SellerID, &seller.NamaToko, &seller.DeskripsiToko, &seller.NoTelepon, &seller.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			sellers = append(sellers, seller)
		}

		c.JSON(http.StatusOK, sellers)
	})

	// get menu
	r.GET("/api/menus", func(c *gin.Context) {
		rows, err := db.Query("SELECT menu_id, seller_id, foto_produk, nama_menu, deskripsi_menu, harga, stok, loves, kategori FROM menu")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var menus []Menu

		for rows.Next() {
			var menu Menu
			err := rows.Scan(&menu.MenuID, &menu.SellerID, &menu.FotoProduk, &menu.NamaMenu, &menu.DeskripsiMenu, &menu.Harga, &menu.Stok, &menu.Loves, &menu.Kategori)
			if err != nil {
				log.Println("Error scanning rows:", err) // Add this line for error logging
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			menus = append(menus, menu)
		}

		c.JSON(http.StatusOK, menus)
	})

	//get seller's menu
	r.GET("/api/menus/:sellerID/product", func(c *gin.Context) {
		sellerID := c.Param("sellerID")

		sellerIDInt, err := strconv.Atoi(sellerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sellerID"})
			return
		}

		rows, err := db.Query("SELECT menu_id, seller_id, nama_menu, deskripsi_menu, harga, stok, loves, kategori FROM menu WHERE seller_id = ?", sellerIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var menus []Menu

		for rows.Next() {
			var menu Menu
			err := rows.Scan(&menu.MenuID, &menu.SellerID, &menu.NamaMenu, &menu.DeskripsiMenu, &menu.Harga, &menu.Stok, &menu.Loves)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			menus = append(menus, menu)
		}

		c.JSON(http.StatusOK, menus)
	})

	// get user
	r.GET("/api/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT user_id, nisn, nama, kelas, no_telepon, email FROM user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User
			err := rows.Scan(&user.UserID, &user.NISN, &user.Nama, &user.Kelas, &user.NoTelepon, &user.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	})

	// get top 3 menus by loves
	r.GET("/api/menus/toploves", func(c *gin.Context) {
		rows, err := db.Query("SELECT menu_id, seller_id, foto_produk, nama_menu, deskripsi_menu, harga, stok, loves, kategori FROM menu ORDER BY loves DESC LIMIT 3")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var menus []Menu

		for rows.Next() {
			var menu Menu
			err := rows.Scan(&menu.MenuID, &menu.SellerID, &menu.FotoProduk, &menu.NamaMenu, &menu.DeskripsiMenu, &menu.Harga, &menu.Stok, &menu.Loves, &menu.Kategori)
			if err != nil {
				log.Println("Error scanning rows:", err) // Add this line for error logging
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
				return
			}
			menus = append(menus, menu)
		}

		c.JSON(http.StatusOK, menus)
	})

	// post user
	r.POST("/api/users", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		result, err := db.Exec("INSERT INTO user (nisn, nama, kelas, no_telepon, email) VALUES (?, ?, ?, ?, ?)",
			newUser.NISN, newUser.Nama, newUser.Kelas, newUser.NoTelepon, newUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data into database"})
			return
		}

		userID, _ := result.LastInsertId()

		newUser.UserID = int(userID)

		c.JSON(http.StatusCreated, newUser)
	})

	//post seller
	r.POST("/api/sellers", func(c *gin.Context) {
		var newSeller Seller
		if err := c.ShouldBindJSON(&newSeller); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		result, err := db.Exec("INSERT INTO seller (nama_toko, deskripsi_toko, no_telepon, email) VALUES (?, ?, ?, ?)",
			newSeller.NamaToko, newSeller.DeskripsiToko, newSeller.NoTelepon, newSeller.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data into database"})
			return
		}

		sellerID, _ := result.LastInsertId()

		newSeller.SellerID = int(sellerID)

		c.JSON(http.StatusCreated, newSeller)
	})

	// post menu
	r.POST("/api/menus", func(c *gin.Context) {
		var newMenu Menu
		if err := c.ShouldBindJSON(&newMenu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		result, err := db.Exec("INSERT INTO menu (seller_id, nama_menu, deskripsi_menu, harga, stok, kategori) VALUES (?, ?, ?, ?, ?, ?)",
			newMenu.SellerID, newMenu.NamaMenu, newMenu.DeskripsiMenu, newMenu.Harga, newMenu.Stok, newMenu.Kategori)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data into database"})
			return
		}

		menuID, _ := result.LastInsertId()

		newMenu.MenuID = int(menuID)

		c.JSON(http.StatusCreated, newMenu)
	})

	// post transaction
	r.POST("/api/transactions", func(c *gin.Context) {
		var newTransaction Transaction
		if err := c.ShouldBindJSON(&newTransaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		result, err := db.Exec("INSERT INTO transaction (user_id, menu_id, jumlah, total_harga, status) VALUES (?, ?, ?, ?, ?)",
			newTransaction.UserID, newTransaction.MenuID, newTransaction.Jumlah, newTransaction.TotalHarga, newTransaction.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data into database"})
			return
		}

		transactionID, _ := result.LastInsertId()

		newTransaction.TransactionID = int(transactionID)

		c.JSON(http.StatusCreated, newTransaction)
	})

	// add stock
	r.PUT("/api/menus/:menuID/addstock", func(c *gin.Context) {
		menuID := c.Param("menuID")

		menuIDInt, err := strconv.Atoi(menuID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menuID"})
			return
		}

		var stockUpdate struct {
			StockToAdd int `json:"stock_to_add"`
		}

		if err := c.ShouldBindJSON(&stockUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		result, err := db.Exec("UPDATE menu SET stok = stok + ? WHERE menu_id = ?", stockUpdate.StockToAdd, menuIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating stock in the database"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully"})
	})

	// change transaction status
	r.PUT("/api/transactions/:transactionID/changestatus", func(c *gin.Context) {
		transactionID := c.Param("transactionID")

		transactionIDInt, err := strconv.Atoi(transactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var changeStatus struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&changeStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		result, err := db.Exec("UPDATE transaction SET status = ? WHERE transaction_id = ?", changeStatus.Status, transactionIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating status in the database"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Status changed successfully"})
	})

	// add loves
	r.PUT("/api/loves/:menuID/add", func(c *gin.Context) {
		menuID := c.Param("menuID")

		menuIDInt, err := strconv.Atoi(menuID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var existingLoves int
		err = db.QueryRow("SELECT loves FROM menu WHERE menu_id = ?", menuIDInt).Scan(&existingLoves)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found"})
			return
		}

		_, err = db.Exec("UPDATE menu SET loves = ? WHERE menu_id = ?", existingLoves+1, menuIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating loves in the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Loves incremented successfully"})
	})

	// delete user
	r.DELETE("/api/users/:userID/delete", func(c *gin.Context) {
		userID := c.Param("userID")

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
			return
		}

		result, err := db.Exec("DELETE FROM user WHERE user_id = ?", userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting data from database"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	// delete seller
	r.DELETE("/api/sellers/:sellerID/delete", func(c *gin.Context) {
		sellerID := c.Param("sellerID")

		sellerIDInt, err := strconv.Atoi(sellerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sellerID"})
			return
		}

		result, err := db.Exec("DELETE FROM seller WHERE seller_id = ?", sellerIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting data from database"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Seller deleted successfully"})
	})

	// delete menu
	r.DELETE("/api/menus/:menuID/delete", func(c *gin.Context) {
		menuID := c.Param("menuID")

		menuIDInt, err := strconv.Atoi(menuID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menuID"})
			return
		}

		result, err := db.Exec("DELETE FROM menu WHERE menu_id = ?", menuIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting data from database"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
	})

	// delete transaction
	r.DELETE("/api/transactions/:transactionID/delete", func(c *gin.Context) {
		transactionID := c.Param("transactionID")

		transactionIDInt, err := strconv.Atoi(transactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transactionID"})
			return
		}

		result, err := db.Exec("DELETE FROM transaction WHERE transaction_id = ?", transactionIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting data from database"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
	})

	r.Run(":8080")
}
