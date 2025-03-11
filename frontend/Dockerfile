# Use Node.js as base image
FROM node:18

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the entire project
COPY . .

# Expose frontend port
EXPOSE 3000

# Start React app
CMD ["npm", "start"]
