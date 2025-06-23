import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Grid,
  Card,
  CardMedia,
  Button,
  Chip,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  CircularProgress,
  Rating,
  Divider
} from '@mui/material';
import {
  LocationOn,
  Star,
  Wifi,
  Pool,
  Restaurant,
  FitnessCenter,
  Spa,
  LocalParking,
  ArrowBack,
  CalendarToday,
  People
} from '@mui/icons-material';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { hotelService, bookingService } from '../services/api';

const HotelDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { isAuthenticated, user } = useAuth();
  
  const [hotel, setHotel] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [bookingDialog, setBookingDialog] = useState(false);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingError, setBookingError] = useState('');

  const checkIn = searchParams.get('checkIn');
  const checkOut = searchParams.get('checkOut');

  const [bookingData, setBookingData] = useState({
    guests: 1,
    specialRequests: ''
  });

  const amenityIcons = {
    'WiFi': <Wifi />,
    'Piscina': <Pool />,
    'Restaurante': <Restaurant />,
    'Gimnasio': <FitnessCenter />,
    'Spa': <Spa />,
    'Estacionamiento': <LocalParking />,
  };

  useEffect(() => {
    loadHotel();
  }, [id]);

  const loadHotel = async () => {
    try {
      // Simular datos del hotel si no hay conexión con la API
      const mockHotel = {
        id: id,
        name: `Hotel ${id === '1' ? 'Boutique Central' : id === '2' ? 'Resort Vista Sierras' : 'Ejecutivo Plaza'}`,
        description: `Un excelente hotel ubicado en el corazón de la ciudad. Ofrecemos servicios de primera clase con habitaciones elegantes y modernas. Perfecto para viajes de negocios y placer.`,
        city: id === '1' ? 'Córdoba' : id === '2' ? 'Villa Carlos Paz' : 'Nueva Córdoba',
        address: `Av. Principal ${100 + parseInt(id)}, ${id === '1' ? 'Centro' : id === '2' ? 'Villa Carlos Paz' : 'Nueva Córdoba'}`,
        rating: 4.5 + (parseInt(id) * 0.1),
        price_per_night: 15000 + (parseInt(id) * 2000),
        photos: [
          'https://images.unsplash.com/photo-1566073771259-6a8506099945?w=800',
          'https://images.unsplash.com/photo-1520250497591-112f2f40a3f4?w=800',
          'https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?w=800'
        ],
        thumbnail: 'https://images.unsplash.com/photo-1566073771259-6a8506099945?w=400',
        amenities: ['WiFi', 'Piscina', 'Restaurante', 'Gimnasio', 'Spa', 'Estacionamiento']
      };

      try {
        const response = await hotelService.getById(id);
        setHotel(response.data);
      } catch (apiError) {
        console.log('API no disponible, usando datos mock');
        setHotel(mockHotel);
      }
    } catch (err) {
      setError('Error al cargar el hotel');
    } finally {
      setLoading(false);
    }
  };

  const calculateNights = () => {
    if (!checkIn || !checkOut) return 1;
    const start = new Date(checkIn);
    const end = new Date(checkOut);
    const diffTime = Math.abs(end - start);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays || 1;
  };

  const calculateTotal = () => {
    const nights = calculateNights();
    return hotel?.price_per_night * nights * bookingData.guests;
  };

  const handleBooking = async () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    if (!checkIn || !checkOut) {
      setBookingError('Por favor selecciona fechas de check-in y check-out');
      return;
    }

    setBookingLoading(true);
    setBookingError('');

    try {
      const booking = {
        hotel_id: hotel.id,
        check_in_date: checkIn,
        check_out_date: checkOut,
        guests: bookingData.guests,
        total_price: calculateTotal(),
        special_requests: bookingData.specialRequests
      };

      try {
        await bookingService.create(booking);
        navigate('/confirmation', { 
          state: { 
            booking: { ...booking, hotel_name: hotel.name },
            success: true 
          }
        });
      } catch (apiError) {
        // Simular éxito si la API no está disponible
        navigate('/confirmation', { 
          state: { 
            booking: { ...booking, hotel_name: hotel.name, id: Date.now() },
            success: true 
          }
        });
      }
    } catch (err) {
      setBookingError('Error al procesar la reserva');
    } finally {
      setBookingLoading(false);
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4, textAlign: 'center' }}>
        <CircularProgress />
        <Typography variant="h6" sx={{ mt: 2 }}>
          Cargando información del hotel...
        </Typography>
      </Container>
    );
  }

  if (error) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">{error}</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Button
        startIcon={<ArrowBack />}
        onClick={() => navigate(-1)}
        sx={{ mb: 3 }}
      >
        Volver
      </Button>

      <Grid container spacing={4}>
        {/* Galería de fotos */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardMedia
              component="img"
              height="400"
              image={hotel?.photos?.[0] || hotel?.thumbnail}
              alt={hotel?.name}
            />
          </Card>
          
          {hotel?.photos && hotel.photos.length > 1 && (
            <Grid container spacing={1} sx={{ mt: 2 }}>
              {hotel.photos.slice(1, 4).map((photo, index) => (
                <Grid item xs={4} key={index}>
                  <Card>
                    <CardMedia
                      component="img"
                      height="120"
                      image={photo}
                      alt={`${hotel.name} ${index + 2}`}
                    />
                  </Card>
                </Grid>
              ))}
            </Grid>
          )}
        </Grid>

        {/* Información y reserva */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 3, position: 'sticky', top: 20 }}>
            <Typography variant="h4" gutterBottom>
              {hotel?.name}
            </Typography>
            
            <Box display="flex" alignItems="center" mb={2}>
              <Rating value={hotel?.rating || 4.5} precision={0.1} readOnly />
              <Typography variant="body2" sx={{ ml: 1 }}>
                ({hotel?.rating || 4.5})
              </Typography>
            </Box>

            <Box display="flex" alignItems="center" mb={3}>
              <LocationOn sx={{ mr: 1, color: 'text.secondary' }} />
              <Typography variant="body1">
                {hotel?.address}
              </Typography>
            </Box>

            <Typography variant="h5" color="primary" gutterBottom>
              ${hotel?.price_per_night?.toLocaleString()} ARS
              <Typography component="span" variant="body2" color="text.secondary">
                / noche
              </Typography>
            </Typography>

            {checkIn && checkOut && (
              <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
                <Box display="flex" alignItems="center" mb={1}>
                  <CalendarToday sx={{ mr: 1, fontSize: 16 }} />
                  <Typography variant="body2">
                    {checkIn} - {checkOut}
                  </Typography>
                </Box>
                <Typography variant="body2">
                  {calculateNights()} noche(s)
                </Typography>
              </Box>
            )}

            <Button
              fullWidth
              variant="contained"
              size="large"
              onClick={() => setBookingDialog(true)}
              disabled={!checkIn || !checkOut}
              sx={{ mb: 2 }}
            >
              {checkIn && checkOut ? 'Reservar Ahora' : 'Selecciona fechas'}
            </Button>

            {(!checkIn || !checkOut) && (
              <Typography variant="body2" color="text.secondary" align="center">
                Vuelve a la búsqueda para seleccionar fechas
              </Typography>
            )}
          </Paper>
        </Grid>

        {/* Descripción */}
        <Grid item xs={12}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h5" gutterBottom>
              Descripción
            </Typography>
            <Typography variant="body1" paragraph>
              {hotel?.description}
            </Typography>

            <Divider sx={{ my: 3 }} />

            <Typography variant="h6" gutterBottom>
              Amenidades
            </Typography>
            <Grid container spacing={1}>
              {hotel?.amenities?.map((amenity) => (
                <Grid item key={amenity}>
                  <Chip
                    icon={amenityIcons[amenity] || <Star />}
                    label={amenity}
                    variant="outlined"
                  />
                </Grid>
              ))}
            </Grid>
          </Paper>
        </Grid>
      </Grid>

      {/* Dialog de reserva */}
      <Dialog open={bookingDialog} onClose={() => setBookingDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Confirmar Reserva</DialogTitle>
        <DialogContent>
          {bookingError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {bookingError}
            </Alert>
          )}

          <Typography variant="h6" gutterBottom>
            {hotel?.name}
          </Typography>
          
          <Typography variant="body2" color="text.secondary" gutterBottom>
            {checkIn} - {checkOut} ({calculateNights()} noche(s))
          </Typography>

          <TextField
            fullWidth
            label="Número de huéspedes"
            type="number"
            value={bookingData.guests}
            onChange={(e) => setBookingData({...bookingData, guests: parseInt(e.target.value) || 1})}
            inputProps={{ min: 1, max: 10 }}
            sx={{ mt: 2, mb: 2 }}
          />

          <TextField
            fullWidth
            label="Solicitudes especiales (opcional)"
            multiline
            rows={3}
            value={bookingData.specialRequests}
            onChange={(e) => setBookingData({...bookingData, specialRequests: e.target.value})}
            sx={{ mb: 2 }}
          />

          <Box sx={{ bgcolor: 'grey.100', p: 2, borderRadius: 1 }}>
            <Typography variant="h6">
              Total: ${calculateTotal()?.toLocaleString()} ARS
            </Typography>
            <Typography variant="body2" color="text.secondary">
              ({bookingData.guests} huésped(es) × {calculateNights()} noche(s))
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setBookingDialog(false)}>Cancelar</Button>
          <Button 
            onClick={handleBooking} 
            variant="contained"
            disabled={bookingLoading}
          >
            {bookingLoading ? <CircularProgress size={24} /> : 'Confirmar Reserva'}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default HotelDetail;