# Stage 1: Build static files
FROM oven/bun:1.1.13-alpine AS builder

WORKDIR /app

ARG VITE_MODE
ARG VITE_API_URL
ENV VITE_MODE=$VITE_MODE
ENV VITE_API_URL=$VITE_API_URL

COPY package.json bun.lock ./
RUN bun install

COPY . ./
RUN bun run build

# Pre-compress assets menggunakan Brotli
RUN apk add --no-cache brotli && \
    find /app/dist -type f \( -name "*.js" -o -name "*.css" -o -name "*.html" -o -name "*.json" -o -name "*.svg" -o -name "*.xml" \) \
    -exec brotli -q 11 -f {} \;

# Stage 2: Serve dengan nginx + brotli
FROM fholzer/nginx-brotli:latest

# Buat direktori cache
RUN mkdir -p /var/cache/nginx/fastcgi /var/cache/nginx/proxy && \
    chown -R nginx:nginx /var/cache/nginx

COPY --from=builder /app/dist /usr/share/nginx/html

# Pastikan nginx.conf tersedia di root proyekmu
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["-g", "daemon off;"]
