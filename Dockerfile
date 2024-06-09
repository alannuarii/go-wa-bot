# Gunakan image dasar Go yang resmi
FROM golang:1.22.3-alpine3.18

# Install dependencies yang dibutuhkan termasuk gcc untuk kompilasi sqlite
RUN apk add --no-cache git sqlite gcc g++ libc-dev

# Set environment variables
ENV GO111MODULE=on

WORKDIR /app

# Copy the entire project to the working directory
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# # Install the package
RUN go install -v ./...

# Build aplikasi
RUN go build -o wa-bot

# Ekspose port yang digunakan oleh aplikasi (sesuaikan jika perlu)
EXPOSE 8080

# Tentukan perintah yang akan dijalankan saat container berjalan
CMD ./wa-bot
