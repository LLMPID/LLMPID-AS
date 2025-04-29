# Stage 1: Build the React Vite App
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

# Copy the rest of the application and build it
COPY . .
RUN yarn build

# Stage 2: Serve the App with NGINX
FROM nginx:alpine

# Set working directory
WORKDIR /usr/share/nginx/html

# Remove default NGINX static files
RUN rm -rf ./*

# Copy built files from the builder stage
COPY --from=builder /app/dist .

# Copy custom NGINX configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080

# Start NGINX
CMD ["nginx", "-g", "daemon off;"]
