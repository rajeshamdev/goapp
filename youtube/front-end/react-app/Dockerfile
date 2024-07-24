
# Use an official Node.js runtime as the base image
FROM node:14-alpine

WORKDIR /app

COPY package.json ./

# install npm dependencies
RUN npm install

# copy the rest of the application code to the working dir
COPY . .

EXPOSE 3000

CMD ["npm", "start"]
