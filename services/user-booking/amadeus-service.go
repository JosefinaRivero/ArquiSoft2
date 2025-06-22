package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AmadeusService struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
	token     string
	tokenExp  time.Time
}

type AmadeusTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type AmadeusHotelOffer struct {
	ID        string `json:"id"`
	Available bool   `json:"available"`
	Price     struct {
		Total    string `json:"total"`
		Currency string `json:"currency"`
	} `json:"price"`
}

type AmadeusOffersResponse struct {
	Data []AmadeusHotelOffer `json:"data"`
}

func NewAmadeusService(apiKey, apiSecret, baseURL string) *AmadeusService {
	if baseURL == "" {
		baseURL = "https://test.api.amadeus.com"
	}

	return &AmadeusService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (as *AmadeusService) getAccessToken() error {
	// Si tenemos un token válido, no necesitamos uno nuevo
	if as.token != "" && time.Now().Before(as.tokenExp) {
		return nil
	}

	// Preparar datos para OAuth2
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", as.apiKey)
	data.Set("client_secret", as.apiSecret)

	// Realizar petición
	resp, err := as.client.PostForm(as.baseURL+"/v1/security/oauth2/token", data)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp AmadeusTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	// Guardar token
	as.token = tokenResp.AccessToken
	as.tokenExp = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // 60 segundos de margen

	return nil
}

func (as *AmadeusService) ValidateBooking(hotelID, checkIn, checkOut string, guests int) (string, error) {
	// Si no tenemos credenciales de Amadeus, simular validación
	if as.apiKey == "" || as.apiKey == "your_amadeus_api_key_here" {
		return as.simulateBookingValidation(hotelID, checkIn, checkOut, guests)
	}

	// Obtener token de acceso
	if err := as.getAccessToken(); err != nil {
		return "", fmt.Errorf("failed to authenticate with Amadeus: %w", err)
	}

	// Mapear hotel ID interno a Amadeus ID (en producción esto vendría de la base de datos)
	amadeusHotelID := as.mapToAmadeusHotelID(hotelID)

	// Verificar disponibilidad en Amadeus
	available, err := as.checkAmadeusAvailability(amadeusHotelID, checkIn, checkOut, guests)
	if err != nil {
		return "", fmt.Errorf("failed to check availability: %w", err)
	}

	if !available {
		return "", fmt.Errorf("hotel not available in Amadeus system")
	}

	// Simular creación de reserva en Amadeus (en producción sería una llamada real)
	bookingID := fmt.Sprintf("AMD_%s_%d", hotelID, time.Now().Unix())

	return bookingID, nil
}

func (as *AmadeusService) checkAmadeusAvailability(amadeusHotelID, checkIn, checkOut string, guests int) (bool, error) {
	// Construir URL para verificar ofertas
	params := url.Values{}
	params.Set("hotelIds", amadeusHotelID)
	params.Set("checkInDate", checkIn)
	params.Set("checkOutDate", checkOut)
	params.Set("adults", strconv.Itoa(guests))

	reqURL := fmt.Sprintf("%s/v3/shopping/hotel-offers?%s", as.baseURL, params.Encode())

	// Crear petición
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+as.token)
	req.Header.Set("Content-Type", "application/json")

	// Realizar petición
	resp, err := as.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Si no encontramos ofertas, consideramos que no está disponible
		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("amadeus API error %d: %s", resp.StatusCode, string(body))
	}

	var offersResp AmadeusOffersResponse
	if err := json.NewDecoder(resp.Body).Decode(&offersResp); err != nil {
		return false, fmt.Errorf("failed to decode offers response: %w", err)
	}

	// Si hay ofertas disponibles, el hotel está disponible
	return len(offersResp.Data) > 0, nil
}

func (as *AmadeusService) simulateBookingValidation(hotelID, checkIn, checkOut string, guests int) (string, error) {
	// Simulación para desarrollo/testing cuando no hay credenciales reales de Amadeus

	// Validaciones básicas
	checkInTime, err := time.Parse("2006-01-02", checkIn)
	if err != nil {
		return "", fmt.Errorf("invalid check-in date format")
	}

	checkOutTime, err := time.Parse("2006-01-02", checkOut)
	if err != nil {
		return "", fmt.Errorf("invalid check-out date format")
	}

	if checkOutTime.Before(checkInTime) {
		return "", fmt.Errorf("check-out date must be after check-in date")
	}

	if guests < 1 || guests > 10 {
		return "", fmt.Errorf("invalid number of guests")
	}

	// Simular algunas condiciones de fallo (10% de probabilidad)
	if time.Now().Unix()%10 == 0 {
		return "", fmt.Errorf("hotel temporarily unavailable")
	}

	// Generar ID de reserva simulado
	bookingID := fmt.Sprintf("SIM_%s_%d", hotelID, time.Now().Unix())

	return bookingID, nil
}

func (as *AmadeusService) mapToAmadeusHotelID(internalHotelID string) string {
	// En producción, esto consultaría la tabla hotel_mappings en la base de datos
	// Por ahora, devolvemos un ID simulado basado en el interno

	// Mapeo simulado para algunos hoteles comunes
	hotelMappings := map[string]string{
		"1": "YXPARKPR", // Hotel ejemplo Amadeus
		"2": "TXPARKNY", // Otro hotel ejemplo
		"3": "LXPARKLD", // Otro hotel ejemplo
	}

	if amadeusID, exists := hotelMappings[internalHotelID]; exists {
		return amadeusID
	}

	// Para hoteles no mapeados, generar un ID basado en el interno
	return fmt.Sprintf("HTL_%s", strings.ToUpper(internalHotelID))
}

func (as *AmadeusService) GetHotelsByCity(cityCode string) ([]string, error) {
	// Si no tenemos credenciales reales, simular
	if as.apiKey == "" || as.apiKey == "your_amadeus_api_key_here" {
		return as.simulateHotelsByCity(cityCode), nil
	}

	// Obtener token
	if err := as.getAccessToken(); err != nil {
		return nil, err
	}

	// Construir URL para buscar hoteles por ciudad
	reqURL := fmt.Sprintf("%s/v1/reference-data/locations/hotels/by-city?cityCode=%s", as.baseURL, cityCode)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+as.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get hotels by city: status %d", resp.StatusCode)
	}

	// Procesar respuesta (implementación específica según la API de Amadeus)
	var hotels []string
	// ... procesamiento de la respuesta real de Amadeus

	return hotels, nil
}

func (as *AmadeusService) simulateHotelsByCity(cityCode string) []string {
	// Simulación para desarrollo
	cityHotels := map[string][]string{
		"COR": {"HTL_COR_001", "HTL_COR_002", "HTL_COR_003"},
		"BUE": {"HTL_BUE_001", "HTL_BUE_002", "HTL_BUE_003"},
		"MDZ": {"HTL_MDZ_001", "HTL_MDZ_002"},
		"BAR": {"HTL_BAR_001", "HTL_BAR_002", "HTL_BAR_003"},
		"PAR": {"YXPARKPR", "TXPARKNY", "LXPARKLD"},
	}

	if hotels, exists := cityHotels[strings.ToUpper(cityCode)]; exists {
		return hotels
	}

	// Devolver hoteles genéricos si no se encuentra la ciudad
	return []string{fmt.Sprintf("HTL_%s_001", strings.ToUpper(cityCode))}
}
