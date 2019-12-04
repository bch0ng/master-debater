// Load env variables
require('dotenv').config();

const MongoClient = require( 'mongodb' ).MongoClient;

let _db = null;

const MONGO_NOT_CONNECTED_ERROR_MESSAGE = "Mongo DB not connected. Please run connectDB() before getDB()";

// Connects to MongoDB and runs the given callback.
function connectDB(callback) {
    MongoClient.connect(process.env.MONGO_URI, { useNewUrlParser: true, useUnifiedTopology: true }, function(err, client) {
        if (err) {
            return err;
        }
        _db = client.db(process.env.MONGO_DB);
        return callback();
    })
}

// Returns the initialized database object.
function getDB() {
    if (_db === null) {
        throw MONGO_NOT_CONNECTED_ERROR_MESSAGE;
    }
    return _db;
}

// Disconnect the database.
function disconnectDB() {
    if (_db === null) {
        throw MONGO_NOT_CONNECTED_ERROR_MESSAGE; 
    }
    _db.close();
}

module.exports = {
    connectDB,
    getDB,
    disconnectDB
};