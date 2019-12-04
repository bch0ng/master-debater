require('dotenv').config();
const http = require("http");
const express = require("express");
const morgan = require("morgan");
const bodyParser = require('body-parser');

const app = express();
const port = process.env.PORT;

const MongoClient = require('./models/mongoClient');
const RabbitMQ = require('./controllers/amqp');
const apiRoute = require('./routes/api');

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.use(morgan('combined'))

const server = http.createServer(app);

//RabbitMQ.connect(process.env.RABBITMQADDR);

app.use("/api/debate/", apiRoute);

MongoClient.connectDB(function () {
    server.listen(port, _ => {
        console.log(`Example app listening on port ${port}!`);
    });
});