package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SearchService struct {
	solrURL        string
	userBookingURL string
	hotelInfoURL   string
	client         *http.Client
}

type SolrHotel struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	City          string   `json:"city"`
	Address       string   `json:"address"`
	Photos        []string `json:"photos"`
	Thumbnail     string   `json:"thumbnail"`
	Amenities     []string `json:"amenities"`
	Rating        float64  `json:"rating"`
	PricePerNight float64  `json:"price_per_night"`
	AmadeusID     string   `json:"amadeus_id"`
	Availability  bool     `json:"availability,omitempty"` // Campo dinámico
}

type SearchResult struct {
	Hotels []SolrHotel `json:"hotels"`
	Total  int         `json:"total"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
}

type SolrResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int         `json:"numFound"`
		Start    int         `json:"start"`
		Docs     []SolrHotel `json:"docs"`
	} `json:"response"`
}

func NewSearchService(solrURL, userBookingURL, hotelInfoURL string) *SearchService {
	return &SearchService{
		solrURL:        solrURL,
		userBookingURL: userBookingURL,
		hotelInfoURL:   hotelInfoURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ss *SearchService) SearchHotels(c *gin.Context) {
	// Parámetros de búsqueda
	city := c.Query("city")
	checkIn := c.Query("checkIn")
	checkOut := c.Query("checkOut")
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	// Convertir parámetros
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	start := (pageInt - 1) * sizeInt

	// Construir consulta Solr
	query := fmt.Sprintf("city:\"%s\"", city)

	// Construir URL de Solr
	params := url.Values{}
	params.Set("q", query)
	params.Set("wt", "json")
	params.Set("start", strconv.Itoa(start))
	params.Set("rows", strconv.Itoa(sizeInt))
	params.Set("sort", "rating desc")

	solrURL := fmt.Sprintf("%s/select?%s", ss.solrURL, params.Encode())

	// Realizar búsqueda en Solr
	resp, err := ss.client.Get(solrURL)
	if err != nil {
		log.Printf("Error querying Solr: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search service unavailable"})
		return
	}
	defer resp.Body.Close()

	var solrResp SolrResponse
	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		log.Printf("Error decoding Solr response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid search response"})
		return
	}

	// Si hay fechas, verificar disponibilidad concurrentemente
	if checkIn != "" && checkOut != "" {
		ss.checkAvailabilityConcurrent(solrResp.Response.Docs, checkIn, checkOut)
	}

	// Construir resultado
	result := SearchResult{
		Hotels: solrResp.Response.Docs,
		Total:  solrResp.Response.NumFound,
		Page:   pageInt,
		Size:   sizeInt,
	}

	c.JSON(http.StatusOK, result)
}

func (ss *SearchService) checkAvailabilityConcurrent(hotels []SolrHotel, checkIn, checkOut string) {
	var wg sync.WaitGroup

	for i := range hotels {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Consultar disponibilidad
			available := ss.checkHotelAvailability(hotels[index].ID, checkIn, checkOut)
			hotels[index].Availability = available
		}(i)
	}

	wg.Wait()
}

func (ss *SearchService) checkHotelAvailability(hotelID, checkIn, checkOut string) bool {
	url := fmt.Sprintf("%s/availability/%s?checkIn=%s&checkOut=%s",
		ss.userBookingURL, hotelID, checkIn, checkOut)

	resp, err := ss.client.Get(url)
	if err != nil {
		log.Printf("Error checking availability for hotel %s: %v", hotelID, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var availabilityResp struct {
		Available bool `json:"available"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&availabilityResp); err != nil {
		log.Printf("Error decoding availability response: %v", err)
		return false
	}

	return availabilityResp.Available
}

func (ss *SearchService) HandleHotelUpdate(messageBody []byte) error {
	var message struct {
		Action    string      `json:"action"`
		HotelID   string      `json:"hotel_id"`
		HotelData interface{} `json:"hotel_data"`
	}

	if err := json.Unmarshal(messageBody, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	switch message.Action {
	case "created", "updated":
		return ss.indexHotel(message.HotelID)
	case "deleted":
		return ss.deleteHotelFromIndex(message.HotelID)
	default:
		log.Printf("Unknown action: %s", message.Action)
		return nil
	}
}

func (ss *SearchService) indexHotel(hotelID string) error {
	// Obtener datos del hotel desde el servicio de ficha
	url := fmt.Sprintf("%s/api/hotels/%s", ss.hotelInfoURL, hotelID)
	resp, err := ss.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get hotel data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("hotel service returned status %d", resp.StatusCode)
	}

	var hotel map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&hotel); err != nil {
		return fmt.Errorf("failed to decode hotel data: %w", err)
	}

	// Convertir a formato Solr
	solrDoc := map[string]interface{}{
		"id":              hotel["id"],
		"name":            hotel["name"],
		"description":     hotel["description"],
		"city":            hotel["city"],
		"address":         hotel["address"],
		"photos":          hotel["photos"],
		"thumbnail":       hotel["thumbnail"],
		"amenities":       hotel["amenities"],
		"rating":          hotel["rating"],
		"price_per_night": hotel["price_per_night"],
		"amadeus_id":      hotel["amadeus_id"],
	}

	// Indexar en Solr
	return ss.indexDocumentInSolr(solrDoc)
}

func (ss *SearchService) indexDocumentInSolr(doc map[string]interface{}) error {
	solrUpdate := map[string]interface{}{
		"add": map[string]interface{}{
			"doc": doc,
		},
	}

	jsonData, err := json.Marshal(solrUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal Solr document: %w", err)
	}

	// POST a Solr
	url := fmt.Sprintf("%s/update?commit=true", ss.solrURL)
	resp, err := ss.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to post to Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Solr returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully indexed hotel %v in Solr", doc["id"])
	return nil
}

func (ss *SearchService) deleteHotelFromIndex(hotelID string) error {
	solrDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": hotelID,
		},
	}

	jsonData, err := json.Marshal(solrDelete)
	if err != nil {
		return fmt.Errorf("failed to marshal Solr delete: %w", err)
	}

	// POST a Solr
	url := fmt.Sprintf("%s/update?commit=true", ss.solrURL)
	resp, err := ss.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to delete from Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Solr delete returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully deleted hotel %s from Solr", hotelID)
	return nil
}
