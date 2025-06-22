package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Middleware de logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Inicializar servicios
	gatewayService := NewGatewayService()

	// Configurar rutas
	api := router.Group("/api")
	{
		// Rutas de autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/login", gatewayService.Login)
			auth.POST("/register", gatewayService.Register)
			auth.GET("/me", AuthMiddleware(), gatewayService.GetProfile)
		}

		// Rutas de hoteles
		hotels := api.Group("/hotels")
		{
			hotels.GET("/search", gatewayService.SearchHotels)
			hotels.GET("/:id", gatewayService.GetHotel)
			hotels.GET("/:id/availability", gatewayService.CheckAvailability)
			
			// Rutas admin
			admin := hotels.Group("/", AuthMiddleware(), AdminMiddleware())
			{
				admin.POST("/", gatewayService.CreateHotel)
				admin.PUT("/:id", gatewayService.UpdateHotel)
				admin.DELETE("/:id", gatewayService.DeleteHotel)
			}
		}

		// Rutas de reservas
		bookings := api.Group("/bookings")
		bookings.Use(AuthMiddleware())
		{
			bookings.POST("/", gatewayService.CreateBooking)
			bookings.GET("/user", gatewayService.GetUserBookings)
			bookings.GET("/", AdminMiddleware(), gatewayService.GetAllBookings)
			bookings.PATCH("/:id/status", AdminMiddleware(), gatewayService.UpdateBookingStatus)
		}

		// Rutas de usuarios
		users := api.Group("/users")
		users.Use(AuthMiddleware())
		{
			users.GET("/profile", gatewayService.GetProfile)
			users.PUT("/profile", gatewayService.UpdateProfile)
			users.GET("/", AdminMiddleware(), gatewayService.GetAllUsers)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// Configurar servidor
	srv := &http.Server{
		Addr:    ":" + getEnv("PORT", "8080"),
		Handler: router,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("🚀 API Gateway running on port %s", getEnv("PORT", "8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited")
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}