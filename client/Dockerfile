FROM oven/bun:1.1.13-alpine AS builder

WORKDIR /app

COPY package.json ./
COPY bun.lock ./
RUN bun install

COPY . ./
RUN bun run build

FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# This will be processed at container startup
COPY public/env.template.js /usr/share/nginx/html/

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh

# Make entrypoint executable
RUN chmod +x /entrypoint.sh

EXPOSE 80

# This will generate env.js before starting nginx
ENTRYPOINT ["/entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]