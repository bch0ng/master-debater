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
const debateChannels = require("./routes/debaterChannels");

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.use(morgan('combined'))

const server = http.createServer(app);

RabbitMQ.connect(process.env.RABBITMQADDR);

app.use("/v1/", function(req, res, next) {
    if (!req.headers['x-user']) {
        return res.status(401).json("Unauthorized");
    }
    next();
});
app.use("/v1/", apiRoute);

app.use("/v2/openchannels", debateChannels);

MongoClient.connectDB(function () {
    console.log("Mongo ConnectDb Start:SanityCheck 123llfawef");
    server.listen(port, _ => {
        console.log(`Example app listening on port ${port}!`);
    });
});
