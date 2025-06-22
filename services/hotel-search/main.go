package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	searchService *SearchService
	rabbitService *RabbitMQService
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Conectar a RabbitMQ
	var err error
	rabbitService, err = NewRabbitMQService(getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitService.Close()

	// Inicializar servicio de búsqueda
	searchService = NewSearchService(
		getEnv("SOLR_URL", "http://localhost:8983/solr/hotels"),
		getEnv("USER_BOOKING_URL", "http://localhost:8083"),
		getEnv("HOTEL_INFO_URL", "http://localhost:8081"),
	)

	// Iniciar consumidor de RabbitMQ
	err = rabbitService.ConsumeMessages("hotel.search.updates", searchService.HandleHotelUpdate)
	if err != nil {
		log.Fatal("Failed to start RabbitMQ consumer:", err)
	}

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configurar rutas
	api := router.Group("/api")
	{
		hotels := api.Group("/hotels")
		{
			hotels.GET("/search", searchService.SearchHotels)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "hotel-search",
			"timestamp": time.Now().Unix(),
		})
	})

	// Configurar servidor
	srv := &http.Server{
		Addr:    ":" + getEnv("PORT", "8082"),
		Handler: router,
	}

	// Iniciar servidor
	go func() {
		log.Printf("🔍 Hotel Search Service running on port %s", getEnv("PORT", "8082"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Hotel Search Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Hotel Search Service exited")
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}