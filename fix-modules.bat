@echo off
echo 🔧 Arreglando módulos Go y creando archivos faltantes...
echo.

REM API Gateway
echo 📁 Configurando API Gateway...
cd services\api-gateway

echo module api-gateway > go.mod
echo. >> go.mod
echo go 1.21 >> go.mod
echo. >> go.mod
echo require ^( >> go.mod
echo     github.com/gin-contrib/cors v1.4.0 >> go.mod
echo     github.com/gin-gonic/gin v1.9.1 >> go.mod
echo     github.com/golang-jwt/jwt/v4 v4.5.0 >> go.mod
echo     github.com/joho/godotenv v1.4.0 >> go.mod
echo ^) >> go.mod

go mod tidy

cd ..\..

REM Hotel Info Service  
echo 📁 Configurando Hotel Info Service...
cd services\hotel-info

echo module hotel-info > go.mod
echo. >> go.mod
echo go 1.21 >> go.mod
echo. >> go.mod
echo require ^( >> go.mod
echo     github.com/gin-gonic/gin v1.9.1 >> go.mod
echo     github.com/joho/godotenv v1.4.0 >> go.mod
echo     github.com/streadway/amqp v1.1.0 >> go.mod
echo     go.mongodb.org/mongo-driver v1.12.1 >> go.mod
echo ^) >> go.mod

go mod tidy

cd ..\..

REM Hotel Search Service
echo 📁 Configurando Hotel Search Service...
cd services\hotel-search

echo module hotel-search > go.mod
echo. >> go.mod  
echo go 1.21 >> go.mod
echo. >> go.mod
echo require ^( >> go.mod
echo     github.com/gin-gonic/gin v1.9.1 >> go.mod
echo     github.com/joho/godotenv v1.4.0 >> go.mod
echo     github.com/streadway/amqp v1.1.0 >> go.mod
echo ^) >> go.mod

go mod tidy

cd ..\..

REM User Booking Service
echo 📁 Configurando User Booking Service...
cd services\user-booking

echo module user-booking > go.mod
echo. >> go.mod
echo go 1.21 >> go.mod  
echo. >> go.mod
echo require ^( >> go.mod
echo     github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874 >> go.mod
echo     github.com/gin-gonic/gin v1.9.1 >> go.mod
echo     github.com/go-sql-driver/mysql v1.7.1 >> go.mod
echo     github.com/golang-jwt/jwt/v4 v4.5.0 >> go.mod
echo     github.com/joho/godotenv v1.4.0 >> go.mod
echo     golang.org/x/crypto v0.12.0 >> go.mod
echo ^) >> go.mod

go mod tidy

cd ..\..

echo ✅ Módulos Go configurados correctamente
echo.
echo 📄 Ahora necesitas crear los archivos .go con el código
echo    Usa los códigos que te proporcioné anteriormente
pause