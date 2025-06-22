import React from 'react';
import {
  Container,
  Typography,
  Box,
  Paper,
  Button,
  Grid,
  Card,
  CardContent,
  Divider,
  Chip
} from '@mui/material';
import {
  CheckCircle,
  Error,
  Home,
  CalendarToday,
  LocationOn,
  People,
  Receipt
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

const Confirmation = () => {
  const navigate = useNavigate();
  const location = useLocation();
  
  // Obtener datos de la reserva del state de navegación
  const booking = location.state?.booking;
  const success = location.state?.success ?? true;

  // Datos mock si no hay información
  const mockBooking = {
    id: Date.now(),
    hotel_name: 'Hotel de Ejemplo',
    check_in_date: '2024-01-15',
    check_out_date: '2024-01-18',
    guests: 2,
    total_price: 45000,
    status: 'confirmed'
  };

  const reservationData = booking || mockBooking;

  const calculateNights = () => {
    const start = new Date(reservationData.check_in_date);
    const end = new Date(reservationData.check_out_date);
    const diffTime = Math.abs(end - start);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays || 1;
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Box textAlign="center" mb={4}>
        {success ? (
          <>
            <CheckCircle 
              sx={{ 
                fontSize: 80, 
                color: 'success.main', 
                mb: 2 
              }} 
            />
            <Typography variant="h3" component="h1" gutterBottom color="success.main">
              ¡Reserva Confirmada!
            </Typography>
            <Typography variant="h6" color="text.secondary">
              Tu reserva ha sido procesada exitosamente
            </Typography>
          </>
        ) : (
          <>
            <Error 
              sx={{ 
                fontSize: 80, 
                color: 'error.main', 
                mb: 2 
              }} 
            />
            <Typography variant="h3" component="h1" gutterBottom color="error.main">
              Reserva Rechazada
            </Typography>
            <Typography variant="h6" color="text.secondary">
              Lo sentimos, no pudimos procesar tu reserva
            </Typography>
          </>
        )}
      </Box>

      {success && (
        <>
          <Paper elevation={3} sx={{ p: 4, mb: 4 }}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
              <Typography variant="h5">
                Detalles de la Reserva
              </Typography>
              <Chip 
                icon={<Receipt />}
                label={`#${reservationData.id}`}
                color="primary"
                variant="outlined"
              />
            </Box>

            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Box display="flex" alignItems="center" mb={2}>
                      <LocationOn sx={{ mr: 1, color: 'primary.main' }} />
                      <Typography variant="h6">
                        {reservationData.hotel_name}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary">
                      Hotel confirmado para tu estadía
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Box display="flex" alignItems="center" mb={2}>
                      <CalendarToday sx={{ mr: 1, color: 'primary.main' }} />
                      <Typography variant="h6">
                        {calculateNights()} noche(s)
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary">
                      Duración de la estadía
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12}>
                <Divider />
              </Grid>

              <Grid item xs={12} md={6}>
                <Box>
                  <Typography variant="subtitle1" gutterBottom>
                    <strong>Check-in</strong>
                  </Typography>
                  <Typography variant="body1">
                    {formatDate(reservationData.check_in_date)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    A partir de las 15:00
                  </Typography>
                </Box>
              </Grid>

              <Grid item xs={12} md={6}>
                <Box>
                  <Typography variant="subtitle1" gutterBottom>
                    <strong>Check-out</strong>
                  </Typography>
                  <Typography variant="body1">
                    {formatDate(reservationData.check_out_date)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Antes de las 12:00
                  </Typography>
                </Box>
              </Grid>

              <Grid item xs={12} md={6}>
                <Box display="flex" alignItems="center">
                  <People sx={{ mr: 1, color: 'text.secondary' }} />
                  <Typography variant="body1">
                    {reservationData.guests} huésped(es)
                  </Typography>
                </Box>
              </Grid>

              <Grid item xs={12} md={6}>
                <Box>
                  <Typography variant="h6" color="primary.main">
                    Total: ${reservationData.total_price?.toLocaleString()} ARS
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          </Paper>

          <Paper elevation={1} sx={{ p: 3, bgcolor: 'info.light', mb: 4 }}>
            <Typography variant="h6" gutterBottom>
              📧 Información importante
            </Typography>
            <Typography variant="body2" paragraph>
              • Te enviaremos un email de confirmación con todos los detalles
            </Typography>
            <Typography variant="body2" paragraph>
              • Presentate en recepción con tu documento de identidad
            </Typography>
            <Typography variant="body2" paragraph>
              • El pago se puede realizar al momento del check-in
            </Typography>
            <Typography variant="body2">
              • Para cancelaciones, contacta al hotel con al menos 24hs de anticipación
            </Typography>
          </Paper>
        </>
      )}

      {!success && (
        <Paper elevation={3} sx={{ p: 4, mb: 4 }}>
          <Typography variant="h6" gutterBottom>
            Motivos posibles del rechazo:
          </Typography>
          <Typography variant="body1" paragraph>
            • El hotel no tiene disponibilidad para las fechas seleccionadas
          </Typography>
          <Typography variant="body1" paragraph>
            • Error en el procesamiento del pago
          </Typography>
          <Typography variant="body1" paragraph>
            • Información incompleta en la reserva
          </Typography>
          <Typography variant="body1">
            Por favor, intenta nuevamente o contacta a nuestro servicio al cliente.
          </Typography>
        </Paper>
      )}

      <Box display="flex" justifyContent="center" gap={2}>
        <Button
          variant="outlined"
          startIcon={<Home />}
          onClick={() => navigate('/')}
          size="large"
        >
          Volver al Inicio
        </Button>
        
        {success && (
          <Button
            variant="contained"
            onClick={() => navigate('/search')}
            size="large"
          >
            Hacer otra reserva
          </Button>
        )}
        
        {!success && (
          <Button
            variant="contained"
            onClick={() => navigate(-2)}
            size="large"
          >
            Intentar nuevamente
          </Button>
        )}
      </Box>
    </Container>
  );
};

export default Confirmation;