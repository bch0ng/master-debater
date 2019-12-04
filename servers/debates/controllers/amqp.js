const amqp = require('amqplib/callback_api');

const QUEUE_NAME = 'debate';

let _conn = null;
let _channel = null;

function connect(host) {
    if (_conn === null) {
        amqp.connect(host, function(error0, connection) {
            if (error0) {
                throw error0;
            }
            console.log("RABBITMQ CONNECTED");
            _conn = connection;
            _conn.createChannel(function(error1, channel) {
                if (error1) {
                    throw error1;
                }
                channel.assertQueue(QUEUE_NAME, {
                    durable: true
                });
                _channel = channel;
            });
        });
    }
}

function publish(json) {
    if (_conn === null || _channel === null) {
        throw "RabbitMQ connection NOT established."
    } else {
        _channel.sendToQueue(QUEUE_NAME, Buffer.from(JSON.stringify(json)));
        console.log(" [x] Sent %s",json);
    }
}

function disconnect() {
    if (_conn === null) {
        throw "RabbitMQ connection NOT established."
    } else {
        _conn.close(); 
        process.exit(0);
    }
}

module.exports = {
    connect,
    publish,
    disconnect
};