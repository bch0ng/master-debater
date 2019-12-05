var mongoose = require("mongoose");

var messageSchema = new mongoose.Schema({
    id: Number,
    channelID: Number,
    body: String,
    createdAt: Date,
    creator: {
        id: Number
    },
    editedAt:  Date
})

module.exports = mongoose.model("Message", messageSchema);