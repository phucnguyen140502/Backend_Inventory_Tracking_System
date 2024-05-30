package models

import (
	"time"

	"github.com/google/uuid"
)

type GoodsReceivedNote struct {
	GoodsReceivedNoteID uuid.UUID `json:"goods_received_note_id" gorm:"primaryKey;autoIncrement"`
	Name                string    `gorm:"index" json:"name"`
	SupplierID          uuid.UUID `json:"supplier_id"` //Foreign Key
	ProductID           uuid.UUID `json:"product_id"`  //Foreign Key
	Amounts             int       `json:"amounts"`
	Price               float32   `json:"price"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (GoodsReceivedNote) TableName() string {
	return "goods_received_note"
}

type Supplier struct {
	SupplierID        uuid.UUID           `json:"supplier_id" gorm:"primaryKey;autoIncrement"`
	Name              string              `gorm:"index" json:"name"`
	PhoneNumber       string              `json:"phone_number"`
	Email             string              `json:"email"`
	GoodsReceivedNote []GoodsReceivedNote `gorm:"foreignKey:SupplierID"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

func (Supplier) TableName() string {
	return "supplier"
}

type GoodsDeliveryNote struct {
	GoodsDeliveryNoteID uuid.UUID `json:"goods_delivery_note_id" gorm:"primaryKey;autoIncrement"`
	Name                string    `gorm:"index" json:"name"`
	ProductID           uuid.UUID `json:"product_id"` //Foreign Key
	Amounts             int       `json:"amounts"`
	Price               float32   `json:"price"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (GoodsDeliveryNote) TableName() string {
	return "goods_delivery_note"
}

type Warehouse struct {
	WarehouseID uuid.UUID `json:"warehouse_id" gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"index" json:"name"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	Product     []Product `gorm:"foreignKey:WarehouseID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Warehouse) TableName() string {
	return "warehouse"
}

type Product struct {
	ProductID         uuid.UUID           `json:"product_id"  gorm:"primaryKey;autoIncrement"`
	Name              string              `gorm:"index" json:"name"`
	UnitPrice         float32             `json:"unit_price"`
	Category          string              `json:"category"`
	WarehouseID       uuid.UUID           `json:"warehouse_id"` // Foreign Key
	InventoryQuantity int64               `json:"inventory_quantity"`
	GoodsDeliveryNote []GoodsDeliveryNote `gorm:"foreignKey:ProductID"`
	GoodsReceivedNote []GoodsReceivedNote `gorm:"foreignKey:ProductID"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

func (Product) TableName() string {
	return "product"
}

type Users struct {
	UserID    uuid.UUID `json:"user_id" gorm:"primaryKey;autoIncrement"`
	FullName  string    `gorm:"index" json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Users) TableName() string {
	return "users"
}
