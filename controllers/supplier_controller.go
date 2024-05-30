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

func SupplierCreate(c *gin.Context) {

	// Get data from req body
	var body struct {
		Name        string
		PhoneNumber string
		Email       string
	}
	c.Bind(&body)
	// // Create a Supplier
	supplier := models.Supplier{
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
		Email:       body.Email,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := database.DB.Create(&supplier)

	// fmt.Print(result)
	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"supplier": supplier,
	})
}
func GetAllSupplierName(c *gin.Context) {
	var Suppliers []models.Supplier
	var names []string

	// Lấy danh sách tất cả các sản phẩm từ cơ sở dữ liệu
	database.DB.Find(&Suppliers)

	// Lặp qua từng sản phẩm và thêm category vào slice nếu nó không tồn tại trong slice đã có
	for _, Supplier := range Suppliers {
		Name := Supplier.Name
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

// show all Supplier
func SupplierIndex(c *gin.Context) {
	// Get the posts
	var suppliers []models.Supplier
	database.DB.Find(&suppliers)

	for i := range suppliers {
		database.DB.Model(&suppliers[i]).Association("GoodsReceivedNote").Find(&suppliers[i].GoodsReceivedNote)
	}

	// Reponse with items
	c.JSON(http.StatusOK, gin.H{
		"suppliers": suppliers,
	})

}

func SupplierUpdate(c *gin.Context) {
	// Get Supplier_id from URL
	supplierID := c.Param("supplier_id")

	// Get data from req body
	var body struct {
		Name        string
		PhoneNumber string
		Email       string
	}
	c.Bind(&body)

	// Find the Supplier we're updating
	var supplier models.Supplier
	database.DB.First(&supplier, "supplier_id = ?", supplierID)

	// Update it
	database.DB.Model(&supplier).Updates(models.Supplier{
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
		Email:       body.Email,
		UpdatedAt:   time.Now(),
	})

	// Response with updated Supplier
	c.JSON(http.StatusOK, gin.H{
		"supplier": supplier,
	})
}

func QuerySupplierByName(c *gin.Context) {
	var suppliers []models.Supplier
	database.DB.Where("name LIKE ?", "%"+c.Param("name")+"%").Find(&suppliers)

	c.JSON(http.StatusOK, gin.H{
		"suppliers": suppliers,
	})
}

func SupplierDelete(c *gin.Context) {
	// Get Supplier_id from URL
	supplierID := c.Param("supplier_id")

	// Delete the Supplier
	database.DB.Delete(models.Supplier{}, "supplier_id = ?", supplierID)

	// Response
	c.Status(200)
}

// Paginate function for GORM
func paginateSuppliers(r *http.Request) func(db *gorm.DB) *gorm.DB {
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

// PaginateSuppliers controller
func PaginateSuppliers(c *gin.Context) {
	var Suppliers []models.Supplier

	if err := database.DB.Scopes(paginateSuppliers(c.Request)).Find(&Suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Suppliers": Suppliers})
}
