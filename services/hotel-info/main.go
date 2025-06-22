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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient   *mongo.Client
	hotelDB       *mongo.Database
	rabbitService *RabbitMQService
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Conectar a MongoDB
	if err := connectMongoDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Conectar a RabbitMQ
	var err error
	rabbitService, err = NewRabbitMQService(getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitService.Close()

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Inicializar servicio
	hotelService := NewHotelService(hotelDB, rabbitService)

	// Configurar rutas
	api := router.Group("/api")
	{
		hotels := api.Group("/hotels")
		{
			hotels.GET("/:id", hotelService.GetHotel)
			hotels.POST("/", hotelService.CreateHotel)
			hotels.PUT("/:id", hotelService.UpdateHotel)
			hotels.DELETE("/:id", hotelService.DeleteHotel)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "hotel-info",
			"timestamp": time.Now().Unix(),
		})
	})

	// Configurar servidor
	srv := &http.Server{
		Addr:    ":" + getEnv("PORT", "8081"),
		Handler: router,
	}

	// Iniciar servidor
	go func() {
		log.Printf("🏨 Hotel Info Service running on port %s", getEnv("PORT", "8081"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Hotel Info Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Hotel Info Service exited")
}

func connectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURL := getEnv("MONGODB_URL", "mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return err
	}

	// Verificar conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	mongoClient = client
	hotelDB = client.Database(getEnv("MONGODB_DATABASE", "hotel_booking"))

	log.Println("✅ Connected to MongoDB")
	return nil
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}