package routers

import (
	"backend/controllers"
	"backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// admin is staff
	// users is customer
	// users phai sign up thi la admin
	// authorized.POST("/users/signup")
	// authorized.POST("/users/login")
	// authorized.GET("/users/product-views") // user xem san pham
	// authorized.POST("/users/orders")       // user dat hang

	// authorized.POST("/signup")
	// authorized.POST("/login")

	// fmt.Printf("Use api")

	incomingRoutes.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Inventory Tracking System API",
		})
	})

	incomingRoutes.POST("/signup", controllers.SignUp)
	incomingRoutes.POST("/login", controllers.Login)
	incomingRoutes.GET("/users", controllers.GetAllUsers)

	authorized := incomingRoutes.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/warehouse", controllers.WarehouseIndex)                   // show warehouse
		authorized.POST("/warehouse", controllers.WareHouseCreate)                 // add warehouse
		authorized.PUT("/warehouse/:warehouse_id", controllers.WarehouseUpdate)    // update warehouse by id
		authorized.DELETE("/warehouse/:warehouse_id", controllers.WarehouseDelete) // delete warehouse
		authorized.GET("/warehouse/:name", controllers.QueryWarehouseByName)
		authorized.GET("/warehouse/paginate", controllers.PaginateWarehouses)    // show following page
		authorized.GET("/warehouse/warehouse-name", controllers.GetAllWarehouse) // show warehouse by id

		authorized.GET("/supplier", controllers.SupplierIndex)                  // show supplier
		authorized.POST("/supplier", controllers.SupplierCreate)                // add supplier
		authorized.PUT("/supplier/:supplier_id", controllers.SupplierUpdate)    // update supplier
		authorized.DELETE("/supplier/:supplier_id", controllers.SupplierDelete) // delete supplier
		authorized.GET("/supplier/:name", controllers.QuerySupplierByName)
		authorized.GET("/supplier/paginate", controllers.PaginateSuppliers)       // show following page
		authorized.GET("/supplier/supplier-name", controllers.GetAllSupplierName) // show supplier by id

		authorized.GET("/products", controllers.ProductIndex)                 // show product
		authorized.POST("/products", controllers.ProductCreate)               // add product
		authorized.PUT("/products/:product_id", controllers.ProductUpdate)    // update product
		authorized.DELETE("/products/:product_id", controllers.ProductDelete) // delete product
		authorized.GET("/products/:name", controllers.QueryProductByName)
		authorized.GET("/category", controllers.GetAllCategory)
		authorized.GET("/category/:category", controllers.QueryProductByCategory)
		authorized.GET("/products/paginate", controllers.PaginateProducts)      // show following page
		authorized.GET("/products/product-name", controllers.GetAllProductName) // show product by id

		authorized.GET("/goods_delivery_note", controllers.GoodsDeliveredNoteIndex)   // show product
		authorized.POST("/goods_delivery_note", controllers.GoodsDeliveredNoteCreate) // add product
		authorized.PUT("/goods_delivery_note/:goods_delivery_note_id", controllers.GoodsDeliveredNoteUpdate)
		authorized.DELETE("/goods_delivery_note/:goods_delivery_note_id", controllers.GoodsDeliveredNoteDelete)
		authorized.GET("/goods_delivery_note/export/:name", controllers.QueryExportByName)
		authorized.GET("/goods_delivery_note/export/paginate", controllers.PaginateExport)

		authorized.GET("/goods_received_note", controllers.GoodsReceivedNoteIndex)   // update product
		authorized.POST("/goods_received_note", controllers.GoodsReceivedNoteCreate) // delete product
		authorized.PUT("/goods_received_note/:goods_received_note_id", controllers.GoodsReceivedNoteUpdate)
		authorized.DELETE("/goods_received_note/:goods_received_note_id", controllers.GoodsReceivedNoteDelete)
		authorized.GET("/goods_received_note/:product_id", controllers.GoodsReceivedNotesByProductId)
		authorized.GET("/goods_received_note/import/:name", controllers.QueryImportByName)
		authorized.GET("/goods_received_note/import/paginate", controllers.PaginateImport)

	}

	// authorized.POST("/import-goods")
	// authorized.POST("/export-goods")

}
