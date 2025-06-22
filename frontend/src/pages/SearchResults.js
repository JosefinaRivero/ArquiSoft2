import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Grid,
  Card,
  CardMedia,
  CardContent,
  CardActions,
  Button,
  Box,
  Chip,
  Rating,
  Skeleton,
  Alert,
  Pagination,
  Paper,
  Divider,
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
} from '@mui/icons-material';
import { useLocation, useNavigate } from 'react-router-dom';
import { hotelService } from '../services/api';

const SearchResults = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [hotels, setHotels] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [searchParams, setSearchParams] = useState({});

  const amenityIcons = {
    'WiFi': <Wifi />,
    'Piscina': <Pool />,
    'Restaurante': <Restaurant />,
    'Gimnasio': <FitnessCenter />,
    'Spa': <Spa />,
    'Estacionamiento': <LocalParking />,
  };

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const searchData = {
      city: params.get('city'),
      checkIn: params.get('checkIn'),
      checkOut: params.get('checkOut'),
    };
    setSearchParams(searchData);
    searchHotels(searchData, 1);
  }, [location.search]);

  const searchHotels = async (params, pageNum = 1) => {
    setLoading(true);
    setError('');

    try {
      const searchQuery = {
        city: params.city,
        checkIn: params.checkIn,
        checkOut: params.checkOut,
        page: pageNum,
        size: 6,
      };

      const response = await hotelService.search(searchQuery);
      setHotels(response.data.hotels || []);
      setTotalPages(Math.ceil(response.data.total / 6));
    } catch (err) {
      setError('Error al buscar hoteles. Intenta nuevamente.');
      console.error('Search error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handlePageChange = (event, value) => {
    setPage(value);
    searchHotels(searchParams, value);
    window.scrollTo(0, 0);
  };

  const handleHotelClick = (hotelId) => {
    const params = new URLSearchParams({
      checkIn: searchParams.checkIn || '',
      checkOut: searchParams.checkOut || '',
    });
    navigate(`/hotel/${hotelId}?${params.toString()}`);
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
    }).format(price);
  };

  const renderHotelCard = (hotel) => (
    <Grid item xs={12} md={6} lg={4} key={hotel.id}>
      <Card
        sx={{
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          cursor: 'pointer',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
          },
        }}
        onClick={() => handleHotelClick(hotel.id)}
      >
        <CardMedia
          component="img"
          height="200"
          image={hotel.thumbnail || 'https://images.unsplash.com/photo-1566073771259-6a8506099945?w=400'}
          alt={hotel.name}
          sx={{
            objectFit: 'cover',
          }}
        />
        <CardContent sx={{ flexGrow: 1, p: 2 }}>
          <Typography variant="h6" component="h3" gutterBottom noWrap>
            {hotel.name}
          </Typography>
          
          <Box display="flex" alignItems="center" mb={1}>
            <LocationOn sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
            <Typography variant="body2" color="text.secondary" noWrap>
              {hotel.address || hotel.city}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" mb={2}>
            <Rating
              value={hotel.rating || 4.5}
              precision={0.1}
              size="small"
              readOnly
            />
            <Typography variant="body2" sx={{ ml: 1 }}>
              ({hotel.rating || 4.5})
            </Typography>
            {hotel.availability !== undefined && (
              <Chip
                label={hotel.availability ? 'Disponible' : 'No disponible'}
                color={hotel.availability ? 'success' : 'error'}
                size="small"
                sx={{ ml: 'auto' }}
              />
            )}
          </Box>

          <Typography
            variant="body2"
            color="text.secondary"
            sx={{
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              mb: 2,
            }}
          >
            {hotel.description}
          </Typography>

          {hotel.amenities && hotel.amenities.length > 0 && (
            <Box sx={{ mb: 2 }}>
              <Box display="flex" flexWrap="wrap" gap={0.5}>
                {hotel.amenities.slice(0, 3).map((amenity, index) => (
                  <Box key={index} display="flex" alignItems="center">
                    {amenityIcons[amenity] || <Star sx={{ fontSize: 14 }} />}
                    <Typography variant="caption" sx={{ ml: 0.5, fontSize: 11 }}>
                      {amenity}
                    </Typography>
                  </Box>
                ))}
                {hotel.amenities.length > 3 && (
                  <Typography variant="caption" color="text.secondary">
                    +{hotel.amenities.length - 3} más
                  </Typography>
                )}
              </Box>
            </Box>
          )}
        </CardContent>

        <Divider />
        
        <CardActions sx={{ p: 2, justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="h6" color="primary.main">
              {formatPrice(hotel.price_per_night)}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              por noche
            </Typography>
          </Box>
          <Button
            variant="contained"
            size="small"
            sx={{ textTransform: 'none' }}
            disabled={hotel.availability === false}
          >
            {hotel.availability === false ? 'No disponible' : 'Ver detalles'}
          </Button>
        </CardActions>
      </Card>
    </Grid>
  );

  const renderSkeleton = () => (
    <Grid item xs={12} md={6} lg={4}>
      <Card sx={{ height: '100%' }}>
        <Skeleton variant="rectangular" height={200} />
        <CardContent>
          <Skeleton variant="text" sx={{ fontSize: '1.2rem' }} />
          <Skeleton variant="text" width="60%" />
          <Skeleton variant="text" width="40%" />
          <Box display="flex" justifyContent="space-between" mt={2}>
            <Skeleton variant="text" width="30%" />
            <Skeleton variant="rectangular" width={80} height={32} />
          </Box>
        </CardContent>
      </Card>
    </Grid>
  );

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* Search Header */}
      <Paper elevation={2} sx={{ p: 3, mb: 4, borderRadius: 2 }}>
        <Typography variant="h4" gutterBottom>
          Hoteles en {searchParams.city}
        </Typography>
        
        {searchParams.checkIn && searchParams.checkOut && (
          <Typography variant="body1" color="text.secondary">
            {searchParams.checkIn} - {searchParams.checkOut}
          </Typography>
        )}
        
        {!loading && hotels.length > 0 && (
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            {hotels.length} hoteles encontrados
          </Typography>
        )}
      </Paper>

      {/* Error Message */}
      {error && (
        <Alert severity="error" sx={{ mb: 4 }}>
          {error}
          <Button
            color="inherit"
            size="small"
            onClick={() => searchHotels(searchParams, page)}
            sx={{ ml: 2 }}
          >
            Reintentar
          </Button>
        </Alert>
      )}

      {/* Results Grid */}
      <Grid container spacing={3}>
        {loading
          ? Array.from(new Array(6)).map((_, index) => (
              <React.Fragment key={index}>{renderSkeleton()}</React.Fragment>
            ))
          : hotels.map((hotel) => renderHotelCard(hotel))}
      </Grid>

      {/* No Results */}
      {!loading && hotels.length === 0 && !error && (
        <Box textAlign="center" py={8}>
          <Typography variant="h5" gutterBottom>
            No se encontraron hoteles
          </Typography>
          <Typography variant="body1" color="text.secondary" mb={3}>
            Intenta buscar en otra ciudad o modifica las fechas
          </Typography>
          <Button
            variant="contained"
            onClick={() => navigate('/')}
            sx={{ textTransform: 'none' }}
          >
            Nueva búsqueda
          </Button>
        </Box>
      )}

      {/* Pagination */}
      {!loading && hotels.length > 0 && totalPages > 1 && (
        <Box display="flex" justifyContent="center" mt={6}>
          <Pagination
            count={totalPages}
            page={page}
            onChange={handlePageChange}
            color="primary"
            size="large"
          />
        </Box>
      )}
    </Container>
  );
};

export default SearchResults;