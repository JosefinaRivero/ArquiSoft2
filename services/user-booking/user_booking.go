package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

type Booking struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	HotelID          string    `json:"hotel_id" db:"hotel_id"`
	AmadeusBookingID string    `json:"amadeus_booking_id" db:"amadeus_booking_id"`
	CheckInDate      string    `json:"check_in_date" db:"check_in_date"`
	CheckOutDate     string    `json:"check_out_date" db:"check_out_date"`
	Guests           int       `json:"guests" db:"guests"`
	TotalPrice       float64   `json:"total_price" db:"total_price"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	UserEmail        string    `json:"user_email,omitempty"`
	HotelName        string    `json:"hotel_name,omitempty"`
}

type BookingRequest struct {
	HotelID      string  `json:"hotel_id" binding:"required"`
	CheckInDate  string  `json:"check_in_date" binding:"required"`
	CheckOutDate string  `json:"check_out_date" binding:"required"`
	Guests       int     `json:"guests" binding:"required,min=1"`
	TotalPrice   float64 `json:"total_price" binding:"required,min=0"`
}

type AvailabilityResponse struct {
	Available bool `json:"available"`
}

type BookingService struct {
	db      *sql.DB
	cache   *memcache.Client
	amadeus *AmadeusService
}

func NewBookingService(database *sql.DB, cacheClient *memcache.Client, amadeusService *AmadeusService) *BookingService {
	return &BookingService{
		db:      database,
		cache:   cacheClient,
		amadeus: amadeusService,
	}
}

func (bs *BookingService) CreateBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar disponibilidad
	available, err := bs.checkAvailabilityInternal(req.HotelID, req.CheckInDate, req.CheckOutDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
		return
	}

	if !available {
		c.JSON(http.StatusConflict, gin.H{"error": "Hotel not available for selected dates"})
		return
	}

	// Validar con Amadeus
	amadeusBookingID, err := bs.amadeus.ValidateBooking(req.HotelID, req.CheckInDate, req.CheckOutDate, req.Guests)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking validation failed: " + err.Error()})
		return
	}

	// Crear reserva
	result, err := bs.db.Exec(`
		INSERT INTO bookings (user_id, hotel_id, amadeus_booking_id, check_in_date, check_out_date, guests, total_price, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.HotelID, amadeusBookingID, req.CheckInDate, req.CheckOutDate, req.Guests, req.TotalPrice, "confirmed",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking ID"})
		return
	}

	// Obtener reserva creada
	booking, err := bs.getBookingByID(int(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created booking"})
		return
	}

	// Limpiar caché de disponibilidad
	bs.clearAvailabilityCache(req.HotelID)

	c.JSON(http.StatusCreated, booking)
}

func (bs *BookingService) GetUserBookings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	rows, err := bs.db.Query(`
		SELECT b.id, b.user_id, b.hotel_id, b.amadeus_booking_id, b.check_in_date, 
		       b.check_out_date, b.guests, b.total_price, b.status, b.created_at, b.updated_at,
		       u.email as user_email
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		WHERE b.user_id = ?
		ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.HotelID, &booking.AmadeusBookingID,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.Guests, &booking.TotalPrice,
			&booking.Status, &booking.CreatedAt, &booking.UpdatedAt, &booking.UserEmail,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan booking"})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, bookings)
}

func (bs *BookingService) GetAllBookings(c *gin.Context) {
	rows, err := bs.db.Query(`
		SELECT b.id, b.user_id, b.hotel_id, b.amadeus_booking_id, b.check_in_date, 
		       b.check_out_date, b.guests, b.total_price, b.status, b.created_at, b.updated_at,
		       u.email as user_email
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		ORDER BY b.created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.HotelID, &booking.AmadeusBookingID,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.Guests, &booking.TotalPrice,
			&booking.Status, &booking.CreatedAt, &booking.UpdatedAt, &booking.UserEmail,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan booking"})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, bookings)
}

func (bs *BookingService) UpdateBookingStatus(c *gin.Context) {
	bookingID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar status
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"cancelled": true,
		"rejected":  true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	_, err := bs.db.Exec(
		"UPDATE bookings SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		req.Status, bookingID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}

	// Obtener reserva actualizada
	id, _ := strconv.Atoi(bookingID)
	booking, err := bs.getBookingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (bs *BookingService) CheckAvailability(c *gin.Context) {
	hotelID := c.Param("hotelId")
	checkIn := c.Query("checkIn")
	checkOut := c.Query("checkOut")

	if checkIn == "" || checkOut == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "checkIn and checkOut dates are required"})
		return
	}

	available, err := bs.checkAvailabilityInternal(hotelID, checkIn, checkOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
		return
	}

	c.JSON(http.StatusOK, AvailabilityResponse{Available: available})
}

func (bs *BookingService) checkAvailabilityInternal(hotelID, checkIn, checkOut string) (bool, error) {
	// Verificar en caché primero
	cacheKey := fmt.Sprintf("availability:%s:%s:%s", hotelID, checkIn, checkOut)

	if item, err := bs.cache.Get(cacheKey); err == nil {
		var available bool
		if err := json.Unmarshal(item.Value, &available); err == nil {
			return available, nil
		}
	}

	// Consultar base de datos
	var count int
	err := bs.db.QueryRow(`
		SELECT COUNT(*) FROM bookings 
		WHERE hotel_id = ? 
		AND status IN ('confirmed', 'pending')
		AND (
			(check_in_date <= ? AND check_out_date > ?) OR
			(check_in_date < ? AND check_out_date >= ?) OR
			(check_in_date >= ? AND check_out_date <= ?)
		)`,
		hotelID, checkIn, checkIn, checkOut, checkOut, checkIn, checkOut,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	available := count == 0

	// Guardar en caché por 10 segundos
	availableBytes, _ := json.Marshal(available)
	bs.cache.Set(&memcache.Item{
		Key:        cacheKey,
		Value:      availableBytes,
		Expiration: 10, // 10 segundos
	})

	return available, nil
}

func (bs *BookingService) getBookingByID(id int) (*Booking, error) {
	var booking Booking
	err := bs.db.QueryRow(`
		SELECT b.id, b.user_id, b.hotel_id, b.amadeus_booking_id, b.check_in_date, 
		       b.check_out_date, b.guests, b.total_price, b.status, b.created_at, b.updated_at,
		       u.email as user_email
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		WHERE b.id = ?`,
		id,
	).Scan(
		&booking.ID, &booking.UserID, &booking.HotelID, &booking.AmadeusBookingID,
		&booking.CheckInDate, &booking.CheckOutDate, &booking.Guests, &booking.TotalPrice,
		&booking.Status, &booking.CreatedAt, &booking.UpdatedAt, &booking.UserEmail,
	)

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (bs *BookingService) clearAvailabilityCache(hotelID string) {
	// En una implementación real, podrías usar un patrón más sofisticado
	// para limpiar todas las claves relacionadas con este hotel
}
