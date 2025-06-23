package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hotel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name" binding:"required"`
	Description   string             `bson:"description" json:"description" binding:"required"`
	City          string             `bson:"city" json:"city" binding:"required"`
	Address       string             `bson:"address" json:"address" binding:"required"`
	Photos        []string           `bson:"photos" json:"photos"`
	Thumbnail     string             `bson:"thumbnail" json:"thumbnail"`
	Amenities     []string           `bson:"amenities" json:"amenities"`
	Rating        float64            `bson:"rating" json:"rating"`
	PricePerNight float64            `bson:"price_per_night" json:"price_per_night" binding:"required"`
	AmadeusID     string             `bson:"amadeus_id" json:"amadeus_id"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type HotelService struct {
	collection    *mongo.Collection
	rabbitService *RabbitMQService
}

func NewHotelService(db *mongo.Database, rabbit *RabbitMQService) *HotelService {
	return &HotelService{
		collection:    db.Collection("hotels"),
		rabbitService: rabbit,
	}
}

func (hs *HotelService) GetHotel(c *gin.Context) {
	idParam := c.Param("id")

	// Intentar convertir a ObjectID
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	var hotel Hotel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = hs.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&hotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

func (hs *HotelService) CreateHotel(c *gin.Context) {
	var hotel Hotel
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer timestamps
	hotel.CreatedAt = time.Now()
	hotel.UpdatedAt = time.Now()

	// Generar AmadeusID simulado (en producción vendría de Amadeus)
	hotel.AmadeusID = "AMD" + strconv.FormatInt(time.Now().Unix(), 10)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := hs.collection.InsertOne(ctx, hotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel"})
		return
	}

	hotel.ID = result.InsertedID.(primitive.ObjectID)

	// Notificar a RabbitMQ
	hs.notifyHotelChange("created", hotel)

	c.JSON(http.StatusCreated, hotel)
}

func (hs *HotelService) UpdateHotel(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	var updateData Hotel
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer timestamp de actualización
	updateData.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Crear update document sin el ID
	updateDoc := bson.M{
		"$set": bson.M{
			"name":            updateData.Name,
			"description":     updateData.Description,
			"city":            updateData.City,
			"address":         updateData.Address,
			"photos":          updateData.Photos,
			"thumbnail":       updateData.Thumbnail,
			"amenities":       updateData.Amenities,
			"rating":          updateData.Rating,
			"price_per_night": updateData.PricePerNight,
			"updated_at":      updateData.UpdatedAt,
		},
	}

	result, err := hs.collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	// Obtener el hotel actualizado
	var updatedHotel Hotel
	err = hs.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedHotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated hotel"})
		return
	}

	// Notificar a RabbitMQ
	hs.notifyHotelChange("updated", updatedHotel)

	c.JSON(http.StatusOK, updatedHotel)
}

func (hs *HotelService) DeleteHotel(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Obtener el hotel antes de eliminarlo para la notificación
	var hotel Hotel
	err = hs.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&hotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	result, err := hs.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hotel"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	// Notificar a RabbitMQ
	hs.notifyHotelChange("deleted", hotel)

	c.JSON(http.StatusOK, gin.H{"message": "Hotel deleted successfully"})
}

func (hs *HotelService) notifyHotelChange(action string, hotel Hotel) {
	message := map[string]interface{}{
		"action":     action,
		"hotel_id":   hotel.ID.Hex(),
		"hotel_data": hotel,
		"timestamp":  time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal hotel change message: %v", err)
		return
	}

	err = hs.rabbitService.PublishMessage("hotel.events", messageBytes)
	if err != nil {
		log.Printf("Failed to publish hotel change message: %v", err)
	}
}
