# Usa una imagen base oficial de Go
FROM golang:alpine AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos de go.mod y go.sum
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Compila la aplicación
RUN go build -o main ./main.go

# Usa una imagen más ligera para el contenedor final
FROM alpine:latest

# Copia el binario desde la etapa de construcción
COPY --from=builder /app/main /app/

# Exponer el puerto en el que tu aplicación escucha
EXPOSE 3000

# Comando para ejecutar la aplicación
CMD ["/app/main"]
