package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Enum
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCanceled   OrderStatus = "canceled"
)

// Table
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Name      string
	Password  string `json:"-"`
	AvatarURL string
	Role      string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	CartItems []CartItem `json:"-"`
	Orders    []Order    `json:"-"`
	Comments  []Comment  `json:"-"`
}

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Products []Product `json:"-"`
}

type Product struct {
	ID            uint `gorm:"primaryKey"`
	CategoryID    uint
	Name          string
	Description   string
	Price         float64
	StockQuantity uint
	ImageURL      string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Category   Category    // 加這行，用來接收 Category 資料
	CartItems  []CartItem  `json:"-"`
	OrderItems []OrderItem `json:"-"`
	Comments   []Comment   `json:"-"`
}

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
	Quantity  uint
	UnitPrice float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID              uint
	UserID          uint
	TotalAmount     float64
	Status          OrderStatus
	ShippingAddress string
	PaymentMethod   string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	OrderItems []OrderItem
}

type OrderItem struct {
	ID        uint
	OrderID   uint
	ProductID uint
	Quantity  uint
	UnitPrice float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID        uint
	UserID    uint
	ProductID uint
	Content   string
	Rating    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Banner struct {
	ID          uint
	Title       string
	Description string
	ImageURL    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Hook & Util
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
