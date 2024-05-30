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

func ProductCreate(c *gin.Context) {

	// Get data from req body
	var body struct {
		Name              string    `json:"name"`
		UnitPrice         float32   `json:"unit_price"`
		Category          string    `json:"category"`
		WarehouseID       uuid.UUID `json:"warehouse_id"`
		InventoryQuantity int64     `json:"inventory_quantity"`
	}
	c.Bind(&body)
	// // Create a Product
	product := models.Product{
		Name:              body.Name,
		UnitPrice:         body.UnitPrice,
		Category:          body.Category,
		WarehouseID:       body.WarehouseID,
		InventoryQuantity: body.InventoryQuantity,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	result := database.DB.Create(&product)

	// fmt.Print(result)
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Cập nhật product trong bảng warehouse
	var warehouse models.Warehouse
	database.DB.First(&warehouse, "warehouse_id = ?", body.WarehouseID)
	warehouse.Product = append(warehouse.Product, product)
	database.DB.Save(&warehouse)

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// show all Product
func ProductIndex(c *gin.Context) {
	// Get the posts
	var products []models.Product
	database.DB.Find(&products)

	for i := range products {
		database.DB.Model(&products[i]).Association("GoodsDeliveryNote").Find(&products[i].GoodsDeliveryNote)
	}

	for i := range products {
		database.DB.Model(&products[i]).Association("GoodsReceivedNote").Find(&products[i].GoodsReceivedNote)
	}

	// Reponse with items
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})

}

func ProductIndexByID(c *gin.Context) {
	// Get product ID from URL
	productID := c.Param("product_id")

	// Find the product by ID
	var product models.Product
	result := database.DB.First(&product, productID)

	if result.Error != nil {
		c.Status(404)
		return
	}

	// Fetch associated data (if needed)
	database.DB.Model(&product).Association("GoodsDeliveryNote").Find(&product.GoodsDeliveryNote)
	database.DB.Model(&product).Association("GoodsReceivedNote").Find(&product.GoodsReceivedNote)

	// Response with the product
	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func ProductUpdate(c *gin.Context) {
	// Get Product_id from URL
	productID := c.Param("product_id")

	// Get data from req body
	var body struct {
		Name              string    `json:"name"`
		UnitPrice         float32   `json:"unit_price"`
		Category          string    `json:"category"`
		WarehouseID       uuid.UUID `json:"warehouse_id"`
		InventoryQuantity int64     `json:"inventory_quantity"`
	}
	c.Bind(&body)

	// Find the Product we're updating
	var product models.Product
	database.DB.First(&product, "product_id = ?", productID)

	// Update it
	database.DB.Model(&product).Updates(models.Product{
		Name:              body.Name,
		UnitPrice:         body.UnitPrice,
		Category:          body.Category,
		WarehouseID:       body.WarehouseID,
		InventoryQuantity: body.InventoryQuantity,
		UpdatedAt:         time.Now(),
	})

	// Response with updated Product
	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func GetAllProductName(c *gin.Context) {
	var products []models.Product
	var names []string

	// Lấy danh sách tất cả các sản phẩm từ cơ sở dữ liệu
	database.DB.Find(&products)

	// Lặp qua từng sản phẩm và thêm category vào slice nếu nó không tồn tại trong slice đã có
	for _, product := range products {
		Name := product.Name
		exists := false

		// Kiểm tra xem category đã tồn tại trong slice hay chưa
		for _, name := range names {
			if name == Name {
				exists = true
				break
			}
		}

		// Nếu category chưa tồn tại trong slice, thêm nó vào
		if !exists {
			names = append(names, Name)
		}
	}

	// Trả về danh sách các category duy nhất
	c.JSON(http.StatusOK, gin.H{
		"names": names,
	})
}

func GetAllCategory(c *gin.Context) {
	var products []models.Product
	var categories []string

	// Lấy danh sách tất cả các sản phẩm từ cơ sở dữ liệu
	database.DB.Find(&products)

	// Lặp qua từng sản phẩm và thêm category vào slice nếu nó không tồn tại trong slice đã có
	for _, product := range products {
		category := product.Category
		exists := false

		// Kiểm tra xem category đã tồn tại trong slice hay chưa
		for _, cat := range categories {
			if cat == category {
				exists = true
				break
			}
		}

		// Nếu category chưa tồn tại trong slice, thêm nó vào
		if !exists {
			categories = append(categories, category)
		}
	}

	// Trả về danh sách các category duy nhất
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func QueryProductByCategory(c *gin.Context) {
	var product []models.Product
	database.DB.Where("category = ?", c.Param("category")).Find(&product)

	c.JSON(http.StatusOK, gin.H{
		"Product": product,
	})
}

func QueryProductByName(c *gin.Context) {
	var product []models.Product
	database.DB.Where("name LIKE ?", "%"+c.Param("name")+"%").Find(&product)

	c.JSON(http.StatusOK, gin.H{
		"Product": product,
	})
}

func ProductDelete(c *gin.Context) {
	// Get Product_id from URL
	productID := c.Param("product_id")

	// Delete the Product
	database.DB.Delete(models.Product{}, "product_id = ?", productID)

	// Response
	c.Status(200)
}

// Paginate function for GORM
func paginateProduct(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		if pageSize <= 0 {
			pageSize = 4
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize).Order("name ASC")
	}
}

// PaginateProduct controller
func PaginateProducts(c *gin.Context) {
	var Product []models.Product

	if err := database.DB.Scopes(paginateProduct(c.Request)).Find(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Product": Product})
}
