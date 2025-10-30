# TAHAP 1: Build Aplikasi React
# Menggunakan image Node.js 18-alpine sebagai dasar
FROM node:18-alpine AS builder

# Tentukan direktori kerja di dalam container
WORKDIR /app

# Salin package.json dan package-lock.json
# Ini memanfaatkan cache Docker. Jika file ini tidak berubah, 'npm install' tidak akan dijalankan ulang.
COPY package*.json ./

# Instal dependensi
RUN npm install

# Salin sisa kode sumber frontend
COPY . .

# Bangun aplikasi untuk produksi
# Ini akan membuat folder /app/build (atau /app/dist jika Anda pakai Vite)
RUN npm run build

# ---

# TAHAP 2: Sajikan Aplikasi dengan Nginx
# Menggunakan image Nginx yang ringan
FROM nginx:stable-alpine

# Salin file build statis dari tahap 'builder' ke direktori default Nginx
# Ganti '/app/build' dengan '/app/dist' jika Anda menggunakan Vite
COPY --from=builder /app/dist /usr/share/nginx/html

# Salin file konfigurasi Nginx kustom kita
# Ini PENTING untuk routing di sisi klien (React Router) dan proxy API
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Ekspos port 80 (port default Nginx)
EXPOSE 80

# Perintah untuk menjalankan Nginx
CMD ["nginx", "-g", "daemon off;"]
