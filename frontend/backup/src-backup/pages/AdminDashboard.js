import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Grid,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Chip,
  IconButton,
  Tabs,
  Tab,
  Alert,
  Fab,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import {
  Add,
  Edit,
  Delete,
  Hotel,
  People,
  BookOnline,
  Visibility,
  Close,
} from '@mui/icons-material';
import { useAuth } from '../context/AuthContext';
import { hotelService, bookingService, userService } from '../services/api';

const AdminDashboard = () => {
  const { isAdmin } = useAuth();
  const [activeTab, setActiveTab] = useState(0);
  const [hotels, setHotels] = useState([]);
  const [bookings, setBookings] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  // Hotel form state
  const [hotelDialog, setHotelDialog] = useState(false);
  const [editingHotel, setEditingHotel] = useState(null);
  const [hotelForm, setHotelForm] = useState({
    name: '',
    description: '',
    city: '',
    address: '',
    price_per_night: '',
    rating: '',
    amenities: [],
    photos: [],
    thumbnail: '',
  });

  const [newAmenity, setNewAmenity] = useState('');
  const [newPhoto, setNewPhoto] = useState('');

  const commonAmenities = [
    'WiFi', 'Piscina', 'Restaurante', 'Gimnasio', 'Spa', 
    'Estacionamiento', 'Aire acondicionado', 'TV', 'Minibar'
  ];

  useEffect(() => {
    if (isAdmin) {
      loadData();
    }
  }, [isAdmin, activeTab]);

  const loadData = async () => {
    setLoading(true);
    setError('');

    try {
      switch (activeTab) {
        case 0: // Hotels
          const hotelsResponse = await hotelService.search({ city: '', size: 100 });
          setHotels(hotelsResponse.data.hotels || []);
          break;
        case 1: // Bookings
          const bookingsResponse = await bookingService.getAll();
          setBookings(bookingsResponse.data || []);
          break;
        case 2: // Users
          const usersResponse = await userService.getAll();
          setUsers(usersResponse.data || []);
          break;
      }
    } catch (err) {
      setError('Error al cargar datos');
      console.error('Load data error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateHotel = () => {
    setEditingHotel(null);
    setHotelForm({
      name: '',
      description: '',
      city: '',
      address: '',
      price_per_night: '',
      rating: '',
      amenities: [],
      photos: [],
      thumbnail: '',
    });
    setHotelDialog(true);
  };

  const handleEditHotel = (hotel) => {
    setEditingHotel(hotel);
    setHotelForm({
      name: hotel.name || '',
      description: hotel.description || '',
      city: hotel.city || '',
      address: hotel.address || '',
      price_per_night: hotel.price_per_night || '',
      rating: hotel.rating || '',
      amenities: hotel.amenities || [],
      photos: hotel.photos || [],
      thumbnail: hotel.thumbnail || '',
    });
    setHotelDialog(true);
  };

  const handleSaveHotel = async () => {
    try {
      const hotelData = {
        ...hotelForm,
        price_per_night: parseFloat(hotelForm.price_per_night),
        rating: parseFloat(hotelForm.rating),
      };

      if (editingHotel) {
        await hotelService.update(editingHotel.id, hotelData);
      } else {
        await hotelService.create(hotelData);
      }

      setHotelDialog(false);
      loadData();
    } catch (err) {
      setError('Error al guardar hotel');
      console.error('Save hotel error:', err);
    }
  };

  const handleDeleteHotel = async (hotelId) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este hotel?')) {
      try {
        await hotelService.delete(hotelId);
        loadData();
      } catch (err) {
        setError('Error al eliminar hotel');
        console.error('Delete hotel error:', err);
      }
    }
  };

  const handleUpdateBookingStatus = async (bookingId, status) => {
    try {
      await bookingService.updateStatus(bookingId, status);
      loadData();
    } catch (err) {
      setError('Error al actualizar reserva');
      console.error('Update booking error:', err);
    }
  };

  const addAmenity = (amenity) => {
    if (amenity && !hotelForm.amenities.includes(amenity)) {
      setHotelForm({
        ...hotelForm,
        amenities: [...hotelForm.amenities, amenity],
      });
    }
    setNewAmenity('');
  };

  const removeAmenity = (amenity) => {
    setHotelForm({
      ...hotelForm,
      amenities: hotelForm.amenities.filter(a => a !== amenity),
    });
  };

  const addPhoto = () => {
    if (newPhoto && !hotelForm.photos.includes(newPhoto)) {
      setHotelForm({
        ...hotelForm,
        photos: [...hotelForm.photos, newPhoto],
      });
    }
    setNewPhoto('');
  };

  const removePhoto = (photo) => {
    setHotelForm({
      ...hotelForm,
      photos: hotelForm.photos.filter(p => p !== photo),
    });
  };

  if (!isAdmin) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">
          No tienes permisos para acceder a esta página.
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h3" gutterBottom>
        Panel de Administración
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <Hotel color="primary" sx={{ fontSize: 40, mr: 2 }} />
                <Box>
                  <Typography variant="h4">{hotels.length}</Typography>
                  <Typography color="text.secondary">Hoteles</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <BookOnline color="success" sx={{ fontSize: 40, mr: 2 }} />
                <Box>
                  <Typography variant="h4">{bookings.length}</Typography>
                  <Typography color="text.secondary">Reservas</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <People color="info" sx={{ fontSize: 40, mr: 2 }} />
                <Box>
                  <Typography variant="h4">{users.length}</Typography>
                  <Typography color="text.secondary">Usuarios</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Tabs */}
      <Paper sx={{ mb: 3 }}>
        <Tabs
          value={activeTab}
          onChange={(e, newValue) => setActiveTab(newValue)}
          indicatorColor="primary"
          textColor="primary"
        >
          <Tab label="Hoteles" />
          <Tab label="Reservas" />
          <Tab label="Usuarios" />
        </Tabs>
      </Paper>

      {/* Hotels Tab */}
      {activeTab === 0 && (
        <Box>
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Nombre</TableCell>
                  <TableCell>Ciudad</TableCell>
                  <TableCell>Precio/Noche</TableCell>
                  <TableCell>Rating</TableCell>
                  <TableCell>Acciones</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {hotels.map((hotel) => (
                  <TableRow key={hotel.id}>
                    <TableCell>{hotel.name}</TableCell>
                    <TableCell>{hotel.city}</TableCell>
                    <TableCell>
                      {new Intl.NumberFormat('es-AR', {
                        style: 'currency',
                        currency: 'ARS',
                      }).format(hotel.price_per_night)}
                    </TableCell>
                    <TableCell>{hotel.rating}</TableCell>
                    <TableCell>
                      <IconButton onClick={() => handleEditHotel(hotel)} size="small">
                        <Edit />
                      </IconButton>
                      <IconButton 
                        onClick={() => handleDeleteHotel(hotel.id)} 
                        size="small"
                        color="error"
                      >
                        <Delete />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>

          <Fab
            color="primary"
            aria-label="add"
            sx={{ position: 'fixed', bottom: 16, right: 16 }}
            onClick={handleCreateHotel}
          >
            <Add />
          </Fab>
        </Box>
      )}

      {/* Bookings Tab */}
      {activeTab === 1 && (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Usuario</TableCell>
                <TableCell>Hotel</TableCell>
                <TableCell>Fechas</TableCell>
                <TableCell>Total</TableCell>
                <TableCell>Estado</TableCell>
                <TableCell>Acciones</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {bookings.map((booking) => (
                <TableRow key={booking.id}>
                  <TableCell>{booking.id}</TableCell>
                  <TableCell>{booking.user_email}</TableCell>
                  <TableCell>{booking.hotel_name}</TableCell>
                  <TableCell>
                    {booking.check_in_date} - {booking.check_out_date}
                  </TableCell>
                  <TableCell>
                    {new Intl.NumberFormat('es-AR', {
                      style: 'currency',
                      currency: 'ARS',
                    }).format(booking.total_price)}
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={booking.status}
                      color={
                        booking.status === 'confirmed' ? 'success' :
                        booking.status === 'cancelled' ? 'error' :
                        booking.status === 'rejected' ? 'error' : 'warning'
                      }
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <FormControl size="small" sx={{ minWidth: 120 }}>
                      <Select
                        value={booking.status}
                        onChange={(e) => handleUpdateBookingStatus(booking.id, e.target.value)}
                      >
                        <MenuItem value="pending">Pendiente</MenuItem>
                        <MenuItem value="confirmed">Confirmada</MenuItem>
                        <MenuItem value="cancelled">Cancelada</MenuItem>
                        <MenuItem value="rejected">Rechazada</MenuItem>
                      </Select>
                    </FormControl>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      {/* Users Tab */}
      {activeTab === 2 && (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Nombre</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Teléfono</TableCell>
                <TableCell>Rol</TableCell>
                <TableCell>Fecha de registro</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.id}</TableCell>
                  <TableCell>{user.name}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{user.phone || '-'}</TableCell>
                  <TableCell>
                    <Chip
                      label={user.role}
                      color={user.role === 'admin' ? 'primary' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    {new Date(user.created_at).toLocaleDateString()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      {/* Hotel Dialog */}
      <Dialog open={hotelDialog} onClose={() => setHotelDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingHotel ? 'Editar Hotel' : 'Crear Hotel'}
          <IconButton
            aria-label="close"
            onClick={() => setHotelDialog(false)}
            sx={{ position: 'absolute', right: 8, top: 8 }}
          >
            <Close />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Nombre"
                value={hotelForm.name}
                onChange={(e) => setHotelForm({ ...hotelForm, name: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Ciudad"
                value={hotelForm.city}
                onChange={(e) => setHotelForm({ ...hotelForm, city: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Descripción"
                multiline
                rows={3}
                value={hotelForm.description}
                onChange={(e) => setHotelForm({ ...hotelForm, description: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Dirección"
                value={hotelForm.address}
                onChange={(e) => setHotelForm({ ...hotelForm, address: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Precio por noche"
                type="number"
                value={hotelForm.price_per_night}
                onChange={(e) => setHotelForm({ ...hotelForm, price_per_night: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Rating"
                type="number"
                inputProps={{ min: 0, max: 5, step: 0.1 }}
                value={hotelForm.rating}
                onChange={(e) => setHotelForm({ ...hotelForm, rating: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Imagen principal (URL)"
                value={hotelForm.thumbnail}
                onChange={(e) => setHotelForm({ ...hotelForm, thumbnail: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom>
                Amenidades
              </Typography>
              <Box display="flex" flexWrap="wrap" gap={1} mb={2}>
                {hotelForm.amenities.map((amenity) => (
                  <Chip
                    key={amenity}
                    label={amenity}
                    onDelete={() => removeAmenity(amenity)}
                    size="small"
                  />
                ))}
              </Box>
              <Box display="flex" gap={1}>
                <FormControl size="small" sx={{ minWidth: 150 }}>
                  <InputLabel>Amenidad</InputLabel>
                  <Select
                    value={newAmenity}
                    onChange={(e) => setNewAmenity(e.target.value)}
                    label="Amenidad"
                  >
                    {commonAmenities
                      .filter(a => !hotelForm.amenities.includes(a))
                      .map((amenity) => (
                        <MenuItem key={amenity} value={amenity}>
                          {amenity}
                        </MenuItem>
                      ))}
                  </Select>
                </FormControl>
                <Button
                  variant="outlined"
                  onClick={() => addAmenity(newAmenity)}
                  disabled={!newAmenity}
                >
                  Agregar
                </Button>
              </Box>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setHotelDialog(false)}>Cancelar</Button>
          <Button onClick={handleSaveHotel} variant="contained">
            {editingHotel ? 'Actualizar' : 'Crear'}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default AdminDashboard;