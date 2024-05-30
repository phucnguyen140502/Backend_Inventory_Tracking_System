package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WareHouseCreate(c *gin.Context) {

	// Get data from req body
	var body struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Capacity int    `json:"capacity"`
	}
	c.Bind(&body)
	// // Create a warehouse
	warehouse := models.Warehouse{
		Name: body.Name, Location: body.Location,
		Capacity: body.Capacity,
	}

	result := database.DB.Create(&warehouse)

	// fmt.Print(result)
	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"warehouse": warehouse,
	})
}

func GetAllWarehouse(c *gin.Context) {
	var warehouses []models.Warehouse
	var names []string

	// Lấy danh sách tất cả các sản phẩm từ cơ sở dữ liệu
	database.DB.Find(&warehouses)

	// Lặp qua từng sản phẩm và thêm category vào slice nếu nó không tồn tại trong slice đã có
	for _, warehouse := range warehouses {
		name := warehouse.Name
		exists := false

		// Kiểm tra xem category đã tồn tại trong slice hay chưa
		for _, warehouse_name := range names {
			if warehouse_name == name {
				exists = true
				break
			}
		}

		// Nếu category chưa tồn tại trong slice, thêm nó vào
		if !exists {
			names = append(names, name)
		}
	}

	// Trả về danh sách các category duy nhất
	c.JSON(http.StatusOK, gin.H{
		"names": names,
	})
}

// show all warehouse
func WarehouseIndex(c *gin.Context) {
	// Get the posts
	var warehouses []models.Warehouse
	database.DB.Find(&warehouses)

	for i := range warehouses {
		database.DB.Model(&warehouses[i]).Association("Product").Find(&warehouses[i].Product)
	}

	// Reponse with items
	c.JSON(http.StatusOK, gin.H{
		"warehouse": warehouses,
	})

}

func WarehouseUpdate(c *gin.Context) {
	// Get warehouse_id from URL
	warehouseID := c.Param("warehouse_id")

	// Get data from req body
	var body struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Capacity int    `json:"capacity"`
	}
	c.Bind(&body)

	// Find the warehouse we're updating
	var warehouse models.Warehouse
	database.DB.First(&warehouse, "warehouse_id = ?", warehouseID)

	// Update it
	database.DB.Model(&warehouse).Updates(models.Warehouse{
		Name:      body.Name,
		Location:  body.Location,
		Capacity:  body.Capacity,
		UpdatedAt: time.Now(),
	})

	// Response with updated warehouse
	c.JSON(http.StatusOK, gin.H{
		"warehouse": warehouse,
	})
}

func QueryWarehouseByName(c *gin.Context) {
	var warehouses []models.Warehouse
	database.DB.Where("name LIKE ?", "%"+c.Param("name")+"%").Find(&warehouses)

	c.JSON(http.StatusOK, gin.H{
		"Warehouses": warehouses,
	})
}

func WarehouseDelete(c *gin.Context) {
	// Get warehouse_id from URL
	warehouseID := c.Param("warehouse_id")

	// Delete the warehouse
	database.DB.Delete(models.Warehouse{}, "warehouse_id = ?", warehouseID)

	// Response
	c.Status(200)
}

// Paginate function for GORM
func paginateWarehouses(r *http.Request) func(db *gorm.DB) *gorm.DB {
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

// PaginateWarehouses controller
func PaginateWarehouses(c *gin.Context) {
	var warehouses []models.Warehouse

	if err := database.DB.Scopes(paginateWarehouses(c.Request)).Find(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
}
