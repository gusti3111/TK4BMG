# --- Stage 1: Build Image (Untuk kompilasi kode Go) ---
FROM golang:1.25.3-alpine AS builder

# Set direktori kerja di dalam container builder
WORKDIR /app/backend

# Menyalin file go.mod dan go.sum untuk mengunduh dependensi
COPY go.mod .
COPY go.sum .

# Mengunduh dependensi (hanya perlu dijalankan ulang jika go.mod berubah)
RUN go mod download

# Menyalin seluruh kode sumber
COPY . .

# Kompilasi aplikasi Go
# -o /app/bin/bmg-server: Menghasilkan binary di lokasi ini
# -tags netgo: Menghasilkan binary statis
# ./cmd/server/main.go: Entry point aplikasi
RUN go build -o /app/bin/bmg-server ./cmd/server/main.go

#
# Tambahkan ini 👇
ENV APP_MODE dev



# --- Stage 2: Final Image (Image runtime yang sangat kecil) ---
FROM alpine:latest

# Menginstal sertifikat CA untuk koneksi HTTPS/SSL
# Penting untuk koneksi keluar atau jika DB menggunakan SSL
RUN apk --no-cache add ca-certificates

# Menyetel timezone ke Asia/Jakarta (opsional, tetapi baik untuk konsistensi log)
ENV TZ Asia/Jakarta

# Set direktori kerja
WORKDIR /app

# Menyalin executable dari build stage
COPY --from=builder /app/bin/bmg-server /app/bmg-server

# Menetapkan port yang akan diekspos (sama dengan yang digunakan Gin :8080)
EXPOSE 8080

# Menentukan perintah saat container dijalankan
CMD ["/app/bmg-server"]