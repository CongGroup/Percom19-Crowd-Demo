FROM node:10.15
# install simple http server for serving static content
RUN npm install -g http-server

WORKDIR /app

COPY package*.json ./

RUN npm i

COPY . .

ARG SERVER_PATH

RUN npm run build

EXPOSE 8080

CMD ["http-server","dist"]