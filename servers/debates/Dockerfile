FROM node:10
WORKDIR /debates

COPY package*.json ./

RUN npm install
COPY . /debates
CMD ["npm","start"]
