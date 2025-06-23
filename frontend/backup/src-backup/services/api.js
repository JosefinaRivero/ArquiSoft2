import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
});

// Interceptor para manejar errores globalmente
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Servicios de hoteles
export const hotelService = {
  search: (params) => api.get('/hotels/search', { params }),
  getById: (id) => api.get(`/hotels/${id}`),
  create: (data) => api.post('/hotels', data),
  update: (id, data) => api.put(`/hotels/${id}`, data),
  delete: (id) => api.delete(`/hotels/${id}`),
  checkAvailability: (hotelId, checkIn, checkOut) => 
    api.get(`/hotels/${hotelId}/availability`, { 
      params: { checkIn, checkOut } 
    }),
};

// Servicios de reservas
export const bookingService = {
  create: (data) => api.post('/bookings', data),
  getUserBookings: () => api.get('/bookings/user'),
  getAll: () => api.get('/bookings'),
  updateStatus: (id, status) => api.patch(`/bookings/${id}/status`, { status }),
};

// Servicios de autenticación
export const authService = {
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
  me: () => api.get('/auth/me'),
};

// Servicios de usuario
export const userService = {
  getProfile: () => api.get('/users/profile'),
  updateProfile: (data) => api.put('/users/profile', data),
  getAll: () => api.get('/users'),
};

export default api;