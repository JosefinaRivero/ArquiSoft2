@echo off
setlocal enabledelayedexpansion

echo.
echo  ██╗  ██╗ ██████╗ ████████╗███████╗██╗         ██████╗  ██████╗  ██████╗ ██╗  ██╗██╗███╗   ██╗ ██████╗ 
echo  ██║  ██║██╔═══██╗╚══██╔══╝██╔════╝██║         ██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝██║████╗  ██║██╔════╝ 
echo  ███████║██║   ██║   ██║   █████╗  ██║         ██████╔╝██║   ██║██║   ██║█████╔╝ ██║██╔██╗ ██║██║  ███╗
echo  ██╔══██║██║   ██║   ██║   ██╔══╝  ██║         ██╔══██╗██║   ██║██║   ██║██╔═██╗ ██║██║╚██╗██║██║   ██║
echo  ██║  ██║╚██████╔╝   ██║   ███████╗███████╗    ██████╔╝╚██████╔╝╚██████╔╝██║  ██╗██║██║ ╚████║╚██████╔╝
echo  ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚══════╝╚══════╝    ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝ ╚═════╝ 
echo.
echo                                  🏨 Sistema de Reservas de Hoteles
echo                                      Arquitectura de Microservicios
echo.

REM Verificar dependencias
echo 🔍 Verificando dependencias...

docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker no está instalado. Descarga desde: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js no está instalado. Descarga desde: https://nodejs.org/
    pause
    exit /b 1
)

go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go no está instalado. Descarga desde: https://golang.org/dl/
    pause
    exit /b 1
)

echo ✅ Todas las dependencias están instaladas

REM Crear estructura de proyecto
echo.
echo 📁 Creando estructura del proyecto...

REM Crear directorios principales
for %%d in (
    "frontend" 
    "services\api-gateway" 
    "services\hotel-info" 
    "services\hotel-search" 
    "services\user-booking"
    "solr\configsets\hotels\conf"
    "haproxy"
    "nginx"
    "uploads"
) do (
    if not exist "%%d" mkdir "%%d"
)

echo ✅ Estructura de directorios creada

REM Crear archivos de configuración
echo.
echo 📄 Creando archivos de configuración...

REM .env file
echo # Configuración del Sistema Hotel Booking > .env
echo JWT_SECRET=your-super-secret-jwt-key-change-in-production >> .env
echo GIN_MODE=debug >> .env
echo. >> .env
echo # Frontend >> .env
echo REACT_APP_API_URL=http://localhost:8080/api >> .env
echo. >> .env
echo # Amadeus API ^(Obtener de https://developers.amadeus.com/^) >> .env
echo AMADEUS_API_KEY=your_amadeus_api_key_here >> .env
echo AMADEUS_API_SECRET=your_amadeus_api_secret_here >> .env
echo AMADEUS_API_URL=https://test.api.amadeus.com >> .env

REM Solr Schema
echo ^<?xml version="1.0" encoding="UTF-8"?^> > solr\configsets\hotels\conf\schema.xml
echo ^<schema name="hotels" version="1.6"^> >> solr\configsets\hotels\conf\schema.xml
echo   ^<field name="id" type="string" indexed="true" stored="true" required="true" multiValued="false"/^> >> solr\configsets\hotels\conf\schema.xml
echo   ^<field name="name" type="text_general" indexed="true" stored="true"/^> >> solr\configsets\hotels\conf\schema.xml
echo   ^<field name="description" type="text_general" indexed="true" stored="true"/^> >> solr\configsets\hotels\conf\schema.xml
echo   ^<field name="city" type="string" indexed="true" stored="true"/^> >> solr\configsets\hotels\conf\schema.xml
echo