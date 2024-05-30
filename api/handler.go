package api

import (
	"backend/database"
	"backend/routers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handler is the main function called by Vercel to handle HTTP requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                       // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Allow all HTTP methods
	config.AllowHeaders = []string{"Authorization", "Content-Type"}           // Allow Authorization and Content-Type headers
	router.Use(cors.New(config))

	// Add a middleware to handle preflight requests
	router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	routers.UserRoutes(router)

	// Handle 404s
	router.NoRoute(func(c *gin.Context) {
		log.Printf("No route found for path: %s", c.Request.URL.Path)
		c.JSON(http.StatusNotFound, gin.H{
			"error":                   "Not found",
			"No route found for path": c.Request.URL.Path,
		})
	})

	router.ServeHTTP(w, r)
}

// StartServer is no longer needed in the Vercel serverless environment
func StartServer() {
	// Lệnh để tắt tường lửa iptables
	// Database connection details
	sql := &database.Sql{
		Host:     getEnv("DB_HOST", "aws-0-ap-southeast-1.pooler.supabase.com"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		UserName: getEnv("DB_USER", "postgres.sjezpemngmqcmmtibiko"),
		PassWord: getEnv("DB_PASSWORD", "QAHYqec4JCNAm9pS"),
		DbName:   getEnv("DB_NAME", "inventory_tracking_system"),
	}
	fmt.Printf("Database connection details: Host=%s, Port=%d, User=%s, DbName=%s\n", sql.Host, sql.Port, sql.UserName, sql.DbName)

	err := sql.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer sql.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                       // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Allow all HTTP methods
	config.AllowHeaders = []string{"Authorization", "Content-Type"}           // Allow Authorization and Content-Type headers
	r.Use(cors.New(config))

	// Add a middleware to handle preflight requests
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	routers.UserRoutes(r)

	r.Run(":3000")
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
