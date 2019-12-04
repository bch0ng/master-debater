FROM node:10
WORKDIR /messaging

COPY package*.json ./

RUN npm install
COPY . /messaging
CMD ["npm","start"]
