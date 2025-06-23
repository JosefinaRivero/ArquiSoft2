import React, { useState } from 'react';
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Box,
  Grid,
  Card,
  CardContent,
  CardMedia,
  Fade,
  Grow,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { Search, LocationOn, Event, Star } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { format } from 'date-fns';

const Home = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useState({
    city: '',
    checkIn: null,
    checkOut: null,
  });

  const handleSearch = () => {
    if (!searchParams.city || !searchParams.checkIn || !searchParams.checkOut) {
      alert('Por favor completa todos los campos');
      return;
    }

    const params = new URLSearchParams({
      city: searchParams.city,
      checkIn: format(searchParams.checkIn, 'yyyy-MM-dd'),
      checkOut: format(searchParams.checkOut, 'yyyy-MM-dd'),
    });

    navigate(`/search?${params.toString()}`);
  };

  const featuredHotels = [
    {
      id: 1,
      name: 'Hotel Boutique Central',
      image: 'https://images.unsplash.com/photo-1566073771259-6a8506099945?w=400',
      rating: 4.8,
      location: 'Centro, Córdoba',
    },
    {
      id: 2,
      name: 'Resort Vista Sierras',
      image: 'https://images.unsplash.com/photo-1520250497591-112f2f40a3f4?w=400',
      rating: 4.6,
      location: 'Villa Carlos Paz',
    },
    {
      id: 3,
      name: 'Hotel Ejecutivo Plaza',
      image: 'https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?w=400',
      rating: 4.7,
      location: 'Nueva Córdoba',
    },
  ];

  return (
    <Box>
      {/* Hero Section */}
      <Box
        sx={{
          backgroundImage: 'linear-gradient(rgba(0,0,0,0.4), rgba(0,0,0,0.4)), url(https://images.unsplash.com/photo-1571003123894-1f0594d2b5d9?w=1200)',
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          minHeight: '70vh',
          display: 'flex',
          alignItems: 'center',
          color: 'white',
        }}
      >
        <Container maxWidth="lg">
          <Fade in timeout={1000}>
            <Box textAlign="center" mb={6}>
              <Typography
                variant="h2"
                component="h1"
                gutterBottom
                sx={{ fontWeight: 'bold', mb: 2 }}
              >
                Encuentra tu Hotel Perfecto
              </Typography>
              <Typography variant="h5" sx={{ mb: 4, opacity: 0.9 }}>
                Descubre experiencias únicas en los mejores hoteles de Argentina
              </Typography>
            </Box>
          </Fade>

          <Grow in timeout={1500}>
            <Paper
              elevation={8}
              sx={{
                p: 4,
                borderRadius: 3,
                backgroundColor: 'rgba(255,255,255,0.95)',
                backdropFilter: 'blur(10px)',
              }}
            >
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <Grid container spacing={3} alignItems="center">
                  <Grid item xs={12} md={3}>
                    <TextField
                      fullWidth
                      label="¿A dónde vamos?"
                      variant="outlined"
                      value={searchParams.city}
                      onChange={(e) =>
                        setSearchParams({ ...searchParams, city: e.target.value })
                      }
                      InputProps={{
                        startAdornment: <LocationOn sx={{ mr: 1, color: 'primary.main' }} />,
                      }}
                    />
                  </Grid>
                  <Grid item xs={12} md={3}>
                    <DatePicker
                      label="Fecha de entrada"
                      value={searchParams.checkIn}
                      onChange={(date) =>
                        setSearchParams({ ...searchParams, checkIn: date })
                      }
                      renderInput={(params) => <TextField {...params} fullWidth />}
                      minDate={new Date()}
                    />
                  </Grid>
                  <Grid item xs={12} md={3}>
                    <DatePicker
                      label="Fecha de salida"
                      value={searchParams.checkOut}
                      onChange={(date) =>
                        setSearchParams({ ...searchParams, checkOut: date })
                      }
                      renderInput={(params) => <TextField {...params} fullWidth />}
                      minDate={searchParams.checkIn || new Date()}
                    />
                  </Grid>
                  <Grid item xs={12} md={3}>
                    <Button
                      fullWidth
                      variant="contained"
                      size="large"
                      onClick={handleSearch}
                      sx={{
                        height: '56px',
                        borderRadius: 2,
                        fontSize: '1.1rem',
                        textTransform: 'none',
                      }}
                      startIcon={<Search />}
                    >
                      Buscar
                    </Button>
                  </Grid>
                </Grid>
              </LocalizationProvider>
            </Paper>
          </Grow>
        </Container>
      </Box>

      {/* Featured Hotels Section */}
      <Container maxWidth="lg" sx={{ py: 8 }}>
        <Typography
          variant="h3"
          component="h2"
          textAlign="center"
          gutterBottom
          sx={{ mb: 6, fontWeight: 'bold' }}
        >
          Hoteles Destacados
        </Typography>
        
        <Grid container spacing={4}>
          {featuredHotels.map((hotel, index) => (
            <Grid item xs={12} md={4} key={hotel.id}>
              <Grow in timeout={1000 + index * 200}>
                <Card
                  sx={{
                    height: '100%',
                    cursor: 'pointer',
                    transition: 'transform 0.3s, box-shadow 0.3s',
                    '&:hover': {
                      transform: 'translateY(-8px)',
                      boxShadow: '0 12px 24px rgba(0,0,0,0.15)',
                    },
                  }}
                  onClick={() => navigate(`/hotel/${hotel.id}`)}
                >
                  <CardMedia
                    component="img"
                    height="240"
                    image={hotel.image}
                    alt={hotel.name}
                  />
                  <CardContent>
                    <Typography variant="h6" component="h3" gutterBottom>
                      {hotel.name}
                    </Typography>
                    <Box display="flex" alignItems="center" mb={1}>
                      <LocationOn sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                      <Typography variant="body2" color="text.secondary">
                        {hotel.location}
                      </Typography>
                    </Box>
                    <Box display="flex" alignItems="center">
                      <Star sx={{ fontSize: 16, mr: 0.5, color: 'orange' }} />
                      <Typography variant="body2">
                        {hotel.rating} (Excelente)
                      </Typography>
                    </Box>
                  </CardContent>
                </Card>
              </Grow>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

export default Home;