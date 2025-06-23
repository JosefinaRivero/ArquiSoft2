FROM node:18-alpine

# Instalar dependencias para node-gyp
RUN apk add --no-cache python3 make g++

WORKDIR /app

# Copiar package.json
COPY package.json package-lock.json* ./

# Instalar dependencias
RUN npm install

# Copiar solo archivos esenciales (sin tests)
COPY src/components/ ./src/components/
COPY src/pages/ ./src/pages/
COPY src/context/ ./src/context/
COPY src/services/ ./src/services/
COPY src/utils/ ./src/utils/
COPY src/App.js ./src/
COPY src/index.js ./src/
COPY public/index.html ./public/

# Variables de entorno
ENV GENERATE_SOURCEMAP=false
ENV DISABLE_ESLINT_PLUGIN=true

EXPOSE 3000
CMD ["npm", "start"]