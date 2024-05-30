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

func GoodsDeliveredNoteCreate(c *gin.Context) {
	// Get data from req body
	var body struct {
		ProductID uuid.UUID `json:"product_id"`
		Amounts   int       `json:"amounts"`
	}
	c.Bind(&body)

	// Query and get the unit_price of the product
	var productPriceName models.Product
	database.DB.First(&productPriceName, body.ProductID)

	// Tính toán giá trị price
	price := float32(body.Amounts) * productPriceName.UnitPrice

	// Create a GoodsDeliveredNote
	GoodsDeliveredNote := models.GoodsDeliveryNote{
		Name:      productPriceName.Name,
		ProductID: body.ProductID,
		Amounts:   body.Amounts,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the GoodsDeliveredNote in the database
	result := database.DB.Create(&GoodsDeliveredNote)
	if result.Error != nil {
		c.Status(400)
		return
	}

	var product models.Product
	if product.GoodsDeliveryNote == nil {
		product.GoodsDeliveryNote = make([]models.GoodsDeliveryNote, 0)
	}

	database.DB.First(&product, "product_id = ?", body.ProductID)

	product.GoodsDeliveryNote = append(product.GoodsDeliveryNote, GoodsDeliveredNote)

	database.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{
		"GoodsDeliveredNote": GoodsDeliveredNote,
	})

}

// show all GoodsDeliveredNote
func GoodsDeliveredNoteIndex(c *gin.Context) {
	// Get the posts
	var GoodsDeliveredNotes []models.GoodsDeliveryNote
	database.DB.Find(&GoodsDeliveredNotes)

	// Reponse with items
	c.JSON(http.StatusOK, gin.H{
		"GoodsDeliveredNotes": GoodsDeliveredNotes,
	})

}

func GoodsDeliveredNoteUpdate(c *gin.Context) {
	// Get Product_id from URL
	GoodsDeliveredNoteID := c.Param("goods_delivery_note_id")

	// Get data from req body
	var body struct {
		SupplierID uuid.UUID `json:"supplier_id"`
		ProductID  uuid.UUID `json:"product_id"`
		Amounts    int       `json:"amounts"`
	}
	c.Bind(&body)

	// Find the Product we're updating
	var GoodsDeliveredNote models.GoodsDeliveryNote
	database.DB.First(&GoodsDeliveredNote, "goods_delivery_note_id = ?", GoodsDeliveredNoteID)

	// Query and get the unit_price of the product
	var productPriceName models.Product
	database.DB.First(&productPriceName, GoodsDeliveredNote.ProductID)

	// Tính toán giá trị price
	price := float32(body.Amounts) * productPriceName.UnitPrice

	// Update it
	database.DB.Model(&GoodsDeliveredNote).Updates(models.GoodsDeliveryNote{
		Name:      GoodsDeliveredNote.Name + " " + string(rune(body.Amounts)),
		ProductID: body.ProductID,
		Amounts:   body.Amounts,
		Price:     price,
		UpdatedAt: time.Now(),
	})

	// Response with updated Product
	c.JSON(http.StatusOK, gin.H{
		"GoodsDeliveredNote": GoodsDeliveredNote,
	})
}

func QueryExportByName(c *gin.Context) {
	var GoodsDeliveryNotes []models.GoodsDeliveryNote
	database.DB.Where("name LIKE ?", "%"+c.Param("name")+"%").Find(&GoodsDeliveryNotes)

	c.JSON(http.StatusOK, gin.H{
		"GoodsDeliveryNotes": GoodsDeliveryNotes,
	})
}

func GoodsDeliveredNoteDelete(c *gin.Context) {
	// Get Product_id from URL
	GoodsDeliveredNoteID := c.Param("goods_delivery_note_id")

	// Delete the Product
	database.DB.Delete(models.GoodsDeliveryNote{}, "goods_delivery_note_id = ?", GoodsDeliveredNoteID)

	// Response
	c.Status(200)
}

// Paginate function for GORM
func paginateExport(r *http.Request) func(db *gorm.DB) *gorm.DB {
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

// PaginateExport controller
func PaginateExport(c *gin.Context) {
	var GoodsDeliveredNotes []models.GoodsDeliveryNote

	if err := database.DB.Scopes(paginateExport(c.Request)).Find(&GoodsDeliveredNotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"GoodsDeliveredNotes": GoodsDeliveredNotes})
}
