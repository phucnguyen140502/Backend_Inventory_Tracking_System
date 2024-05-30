package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GoodsReceivedNoteCreate(c *gin.Context) {
	// Get data from req body
	var body struct {
		SupplierID uuid.UUID `json:"supplier_id"`
		ProductID  uuid.UUID `json:"product_id"`
		Amounts    int       `json:"amounts"`
	}
	c.Bind(&body)

	// Query and get the unit_price of the product
	var productPriceName models.Product
	database.DB.First(&productPriceName, body.ProductID)

	// Tính toán giá trị price
	price := float32(body.Amounts) * productPriceName.UnitPrice

	// Create a GoodsReceivedNote
	GoodsReceivedNote := models.GoodsReceivedNote{
		Name:       productPriceName.Name,
		SupplierID: body.SupplierID,
		ProductID:  body.ProductID,
		Amounts:    body.Amounts,
		Price:      price,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save the GoodsReceivedNote in the database
	result := database.DB.Create(&GoodsReceivedNote)
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Find the corresponding Supplier
	var supplier models.Supplier
	if supplier.GoodsReceivedNote == nil {
		supplier.GoodsReceivedNote = make([]models.GoodsReceivedNote, 0)
	}

	database.DB.First(&supplier, "supplier_id = ?", body.SupplierID)

	supplier.GoodsReceivedNote = append(supplier.GoodsReceivedNote, GoodsReceivedNote)

	database.DB.Save(&supplier)

	var product models.Product
	if product.GoodsReceivedNote == nil {
		product.GoodsReceivedNote = make([]models.GoodsReceivedNote, 0)
	}

	database.DB.First(&product, "product_id = ?", body.ProductID)

	product.GoodsReceivedNote = append(product.GoodsReceivedNote, GoodsReceivedNote)

	database.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{
		"GoodsReceivedNote": GoodsReceivedNote,
	})

}

// show all GoodsReceivedNote
func GoodsReceivedNoteIndex(c *gin.Context) {
	// Get the posts
	var GoodsReceivedNotes []models.GoodsReceivedNote
	database.DB.Find(&GoodsReceivedNotes)

	// Reponse with items
	c.JSON(http.StatusOK, gin.H{
		"GoodsReceivedNotes": GoodsReceivedNotes,
	})

}

func GoodsReceivedNotesByProductId(c *gin.Context) {

	productID := c.Param("product_id")

	// Find the product by ID
	var GoodsReceivedNotes []models.GoodsReceivedNote
	database.DB.Find(&GoodsReceivedNotes, "product_id = ?", productID)

	if len(GoodsReceivedNotes) == 0 {
		// Trả về lỗi 404 Not Found nếu không tìm thấy
		c.JSON(http.StatusNotFound, gin.H{"error": "No goods received notes found for the specified product ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"goods_received_notes": GoodsReceivedNotes,
	})

}

func GoodsReceivedNoteUpdate(c *gin.Context) {
	// Get Product_id from URL
	GoodsReceivedNoteID := c.Param("goods_received_note_id")

	// Get data from req body
	var body struct {
		SupplierID uuid.UUID `json:"supplier_id"`
		ProductID  uuid.UUID `json:"product_id"`
		Amounts    int       `json:"amounts"`
	}
	c.Bind(&body)

	// Find the Product we're updating
	var GoodsReceivedNote models.GoodsReceivedNote
	database.DB.First(&GoodsReceivedNote, "goods_received_note_id = ?", GoodsReceivedNoteID)

	// Query and get the unit_price of the product
	var productPriceName models.Product
	database.DB.First(&productPriceName, GoodsReceivedNote.ProductID)

	// Tính toán giá trị price
	price := float32(body.Amounts) * productPriceName.UnitPrice

	// Update it
	database.DB.Model(&GoodsReceivedNote).Updates(models.GoodsReceivedNote{
		Name:       GoodsReceivedNote.Name,
		SupplierID: body.SupplierID,
		ProductID:  body.ProductID,
		Amounts:    body.Amounts,
		Price:      price,
		UpdatedAt:  time.Now(),
	})

	// Response with updated Product
	c.JSON(http.StatusOK, gin.H{
		"GoodsReceivedNote": GoodsReceivedNote,
	})
}

func QueryImportByName(c *gin.Context) {
	var GoodsReceivedNotes []models.GoodsReceivedNote
	database.DB.Where("name LIKE ?", "%"+c.Param("name")+"%").Find(&GoodsReceivedNotes)

	c.JSON(http.StatusOK, gin.H{
		"GoodsReceivedNotes": GoodsReceivedNotes,
	})
}

func GoodsReceivedNoteDelete(c *gin.Context) {
	// Get Product_id from URL
	GoodsReceivedNoteID := c.Param("goods_received_note_id")

	// Delete the Product
	database.DB.Delete(models.GoodsReceivedNote{}, "goods_received_note_id = ?", GoodsReceivedNoteID)

	// Response
	c.Status(200)
}

// Paginate function for GORM
func paginateImport(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		if pageSize <= 0 {
			pageSize = 3
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize).Order("name ASC")
	}
}

// PaginateImport controller
func PaginateImport(c *gin.Context) {
	var GoodsReceivedNotes []models.GoodsReceivedNote

	if err := database.DB.Scopes(paginateImport(c.Request)).Find(&GoodsReceivedNotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"GoodsReceivedNotes": GoodsReceivedNotes})
}
