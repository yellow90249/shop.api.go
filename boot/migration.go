package boot

import (
	"log"

	"shop.go/model"
)

func Migrate() {
	// 執行 migration
	err := DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.Comment{},
		&model.Banner{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}
