FROM oven/bun:1.1.13-alpine AS builder

WORKDIR /app

COPY package.json ./
COPY bun.lockb ./
RUN bun install

COPY . ./
RUN bun run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]