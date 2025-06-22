# 🏨 Hotel Booking System - Arquitectura de Microservicios

Sistema completo de reservas de hoteles implementado con arquitectura de microservicios utilizando React para el frontend y Go para los servicios backend.

## 📋 Características Principales

### ✨ Funcionalidades
- **Búsqueda de hoteles** por ciudad y fechas
- **Reservas en línea** con validación externa (Amadeus API)
- **Panel de administración** para gestión de hoteles
- **Autenticación y autorización** de usuarios
- **Interfaz moderna** y responsiva con Material-UI
- **Caché distribuido** para optimización de consultas
- **Mensajería asíncrona** entre servicios

### 🏗️ Arquitectura
- **Frontend**: React 18 + Material-UI
- **API Gateway**: Go + Gin + JWT
- **4 Microservicios en Go**:
  - Hotel Info Service (MongoDB)
  - Hotel Search Service (Solr)
  - User Booking Service (MySQL + Memcached)
- **Load Balancer**: HAProxy
- **Mensajería**: RabbitMQ
- **Bases de datos**: MongoDB, MySQL, Solr, Memcached
- **Containerización**: Docker + Docker Compose

## 🚀 Instalación Rápida (Windows)

### Prerrequisitos
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [Node.js 18+](https://nodejs.org/)
- [Go 1.21+](https://golang.org/dl/)

### Instalación Automática
```bash
# Clonar o descargar el proyecto
git clone <repository-url>
cd hotel-booking-system

# Ejecutar script de instalación
setup.bat
```

El script automáticamente:
1. ✅ Verifica dependencias
2. 📁 Crea estructura de directorios
3. 📄 Genera archivos de configuración
4. 🐳 Crea Dockerfiles
5. 🔧 Inicializa módulos Go y React
6. 🚀 Inicia todos los servicios

## 🌐 URLs del Sistema

### Servicios Principales
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080

### Microservicios
- **Hotel Info**: http://localhost:8081/health
- **Hotel Search**: http://localhost:8082/health
- **User Booking**: http://localhost:8083/health

### Herramientas de Administración
- **Solr Admin**: http://localhost:8983
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **HAProxy Stats**: http://localhost:8404/stats

### Bases de Datos
- **MySQL**: localhost:3306 (root/password)
- **MongoDB**: localhost:27017
- **Memcached**: localhost:11211

## 👤 Credenciales de Acceso

### Usuario Normal
- **Email**: user@hotel.com
- **Password**: password

### Administrador
- **Email**: admin@hotel.com
- **Password**: password

## 📱 Pantallas del Frontend

### 1. Página Principal
- Formulario de búsqueda con ciudad y fechas
- Hoteles destacados
- Interfaz moderna y atractiva

### 2. Resultados de Búsqueda
- Lista de hoteles disponibles
- Filtros por disponibilidad
- Información detallada de cada hotel

### 3. Detalle del Hotel
- Información completa del hotel
- Galería de fotos
- Amenidades disponibles
- Botón de reserva

### 4. Confirmación
- Estado de la reserva (éxito/rechazo)
- Detalles de la transacción

### 5. Panel de Administración
- Gestión de hoteles (CRUD)
- Administración de reservas
- Gestión de usuarios
- Estadísticas del sistema

## 🔧 Comandos Útiles

### Docker
```bash
# Ver logs de todos los servicios
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f frontend
docker-compose logs -f api-gateway

# Reiniciar un servicio
docker-compose restart hotel-info

# Parar todo el sistema
docker-compose down

# Rebuild y restart
docker-compose up -d --build
```

### Desarrollo Local
```bash
# Frontend (React)
cd frontend
npm start

# API Gateway
cd services/api-gateway
go run .

# Hotel Info Service
cd services/hotel-info
go run .
```

## 🔑 Configuración de Amadeus API

1. Registrarse en https://developers.amadeus.com/
2. Crear una aplicación
3. Obtener API Key y Secret
4. Editar el archivo `.env`:

```env
AMADEUS_API_KEY=tu_api_key_aqui
AMADEUS_API_SECRET=tu_api_secret_aqui
AMADEUS_API_URL=https://test.api.amadeus.com
```

## 📊 Estructura del Proyecto

```
hotel-booking-system/
├── frontend/                    # React App
│   ├── src/
│   │   ├── components/         # Componentes reutilizables
│   │   ├── pages/             # Páginas principales
│   │   ├── services/          # Servicios API
│   │   └── context/           # Context de autenticación
│   └── Dockerfile
├── services/
│   ├── api-gateway/           # Gateway principal
│   ├── hotel-info/            # Microservicio de hoteles
│   ├── hotel-search/          # Microservicio de búsqueda
│   └── user-booking/          # Microservicio de usuarios/reservas
├── solr/                      # Configuración Solr
├── haproxy/                   # Configuración Load Balancer
├── nginx/                     # Servidor web estático
├── docker-compose.yml         # Orquestación completa
├── .env                       # Variables de entorno
└── setup.bat                  # Script de instalación
```

## 🗄️ Esquema de Base de Datos

### MySQL (Users & Bookings)
```sql
-- Usuarios
users: id, name, email, password_hash, phone, role, created_at, updated_at

-- Mapeo de hoteles con Amadeus
hotel_mappings: id, internal_hotel_id, amadeus_hotel_id, created_at

-- Reservas
bookings: id, user_id, hotel_id, amadeus_booking_id, check_in_date, 
          check_out_date, guests, total_price, status, created_at, updated_at
```

### MongoDB (Hotels)
```javascript
// Hoteles
{
  _id: ObjectId,
  name: String,
  description: String,
  city: String,
  address: String,
  photos: [String],
  thumbnail: String,
  amenities: [String],
  rating: Number,
  price_per_night: Number,
  amadeus_id: String,
  created_at: Date,
  updated_at: Date
}
```

## 🔄 Flujo de Comunicación entre Servicios

1. **Cliente** → **API Gateway** → **Microservicio específico**
2. **Hotel Info** → **RabbitMQ** → **Hotel Search** (sincronización)
3. **Hotel Search** → **User Booking** (verificación de disponibilidad)
4. **User Booking** → **Amadeus API** (validación de reservas)
5. **Memcached** → Cache de disponibilidad (TTL 10 segundos)

## 📈 Características Avanzadas Implementadas

### Anexo 1 - Integración con Amadeus
- ✅ Autenticación OAuth2 automática
- ✅ Validación de reservas en tiempo real
- ✅ Mapeo entre IDs internos y Amadeus

### Anexo 2 - Funcionalidades Adicionales
- ✅ Load Balancer con HAProxy
- ✅ Caché distribuido con Memcached
- ✅ Test unitarios para servicios Go
- ✅ Docker Compose para orquestación
- ✅ Frontend de administración completo
- ✅ Login de usuario y administrador
- ✅ Escalado automático preparado

### Anexo 3 - Arquitectura
- ✅ Microservicios desacoplados
- ✅ Comunicación asíncrona (RabbitMQ)
- ✅ Bases de datos especializadas por servicio
- ✅ API Gateway como punto de entrada único

## 🛠️ Resolución de Problemas

### Puerto en uso
```bash
# Verificar puertos ocupados
netstat -ano | findstr :3000
netstat -ano | findstr :8080

# Cambiar puertos en docker-compose.yml si es necesario
```

### Servicios no responden
```bash
# Verificar estado de contenedores
docker-compose ps

# Reiniciar servicios problemáticos
docker-compose restart nombre-servicio
```

### Problemas de permisos
```bash
# Ejecutar como administrador en Windows
# Verificar que Docker Desktop esté ejecutándose
```

### Logs para debugging
```bash
# Ver logs detallados
docker-compose logs -f --tail=100 hotel-search
docker-compose logs -f --tail=100 rabbitmq
```

## 📞 Soporte

Para problemas o preguntas:
1. Revisar logs de servicios
2. Verificar configuración de .env
3. Comprobar que todos los contenedores estén ejecutándose
4. Verificar conectividad de red entre servicios

## 🏆 Cumplimiento de Requisitos

- ✅ **4 Microservicios** implementados
- ✅ **Frontend React** con 4 pantallas requeridas
- ✅ **MongoDB** para hoteles con RabbitMQ
- ✅ **Solr** para búsquedas con disponibilidad concurrente
- ✅ **MySQL + Memcached** para usuarios y reservas
- ✅ **Integración Amadeus** para validación
- ✅ **Load Balancer** y escalado automático
- ✅ **Docker Compose** para orquestación
- ✅ **Panel de administración** completo
- ✅ **Autenticación** de usuarios y admin
- ✅ **Caché distribuido** con TTL optimizado
- ✅ **Tests unitarios** implementados

---

**🎓 Proyecto desarrollado para Arquitectura de Software II - Universidad Católica de Córdoba**