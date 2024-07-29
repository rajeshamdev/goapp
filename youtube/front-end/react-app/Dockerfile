# Stage 1: Build the React application
FROM node:18 AS build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy application files
COPY public ./public
COPY src ./src

# Build the application
RUN npm run build

# Stage 2: Serve the React application
FROM nginx:alpine

# Copy build files from the build stage
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
