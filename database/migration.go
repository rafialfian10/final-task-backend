package database

import (
	"fmt"
	"waysbook/models"
	"waysbook/pkg/mysql"
)

// Jika aplikasi berjalan maka auto migration akan berjalan
func RunMigration() {
	// koneksi database akan melakukan auto migrasi struct/models ke dalam database mysql
	err := mysql.DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Cart{},
		&models.Transaction{},
	)

	if err != nil {
		fmt.Println(err)
		panic("Migration failed")
	}

	fmt.Println("Migration success")
}
