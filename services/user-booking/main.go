package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db            *sql.DB
	cache         *memcache.Client
	userService   *UserService
	bookingService *BookingService
	amadeusService *AmadeusService
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Conectar a MySQL
	if err := connectMySQL(); err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	defer db.Close()

	// Conectar a Memcached
	cache = memcache.New(getEnv("MEMCACHED_URL", "localhost:11211"))

	// Inicializar servicios
	amadeusService = NewAmadeusService(
		getEnv("AMADEUS_API_KEY", ""),
		getEnv("AMADEUS_API_SECRET", ""),
		getEnv("AMADEUS_API_URL", "https://test.api.amadeus.com"),
	)
	
	userService = NewUserService(db)
	bookingService = NewBookingService(db, cache, amadeusService)

	// Crear tablas si no existen
	if err := createTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
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
		// Autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/login", userService.Login)
			auth.POST("/register", userService.Register)
			auth.GET("/me", AuthMiddleware(), userService.GetProfile)
		}

		// Usuarios
		users := api.Group("/users")
		users.Use(AuthMiddleware())
		{
			users.GET("/profile", userService.GetProfile)
			users.PUT("/profile", userService.UpdateProfile)
			users.GET("/", AdminMiddleware(), userService.GetAllUsers)
		}

		// Reservas
		bookings := api.Group("/bookings")
		bookings.Use(AuthMiddleware())
		{
			bookings.POST("/", bookingService.CreateBooking)
			bookings.GET("/user", bookingService.GetUserBookings)
			bookings.GET("/", AdminMiddleware(), bookingService.GetAllBookings)
			bookings.PATCH("/:id/status", AdminMiddleware(), bookingService.UpdateBookingStatus)
		}

		// Disponibilidad
		api.GET("/availability/:hotelId", bookingService.CheckAvailability)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "user-booking",
			"timestamp": time.Now().Unix(),
		})
	})

	// Configurar servidor
	srv := &http.Server{
		Addr:    ":" + getEnv("PORT", "8083"),
		Handler: router,
	}

	// Iniciar servidor
	go func() {
		log.Printf("👥 User Booking Service running on port %s", getEnv("PORT", "8083"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down User Booking Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ User Booking Service exited")
}

func connectMySQL() error {
	dsn := getEnv("MYSQL_DSN", "root:password@tcp(localhost:3306)/hotel_booking?charset=utf8mb4&parseTime=True&loc=Local")
	
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		return err
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Connected to MySQL")
	return nil
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			phone VARCHAR(50),
			role ENUM('user', 'admin') DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS hotel_mappings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			internal_hotel_id VARCHAR(255) NOT NULL,
			amadeus_hotel_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY unique_internal (internal_hotel_id),
			UNIQUE KEY unique_amadeus (amadeus_hotel_id)
		)`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			hotel_id VARCHAR(255) NOT NULL,
			amadeus_booking_id VARCHAR(255),
			check_in_date DATE NOT NULL,
			check_out_date DATE NOT NULL,
			guests INT DEFAULT 1,
			total_price DECIMAL(10,2) NOT NULL,
			status ENUM('pending', 'confirmed', 'cancelled', 'rejected') DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	// Crear usuario admin por defecto
	_, err := db.Exec(`
		INSERT IGNORE INTO users (name, email, password_hash, role) 
		VALUES ('Admin', 'admin@hotel.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')
	`) // password: password
	
	log.Println("✅ Database tables created/verified")
	return err
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}