FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Expose port
EXPOSE 5173

# Start the development server with WebSocket support
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"] 