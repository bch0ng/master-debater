var mongoose = require("mongoose");

var channelSchema = new mongoose.Schema({
    id: Number,
    name: String,
    description: String,
    private: Boolean,
    members: [],
    createdAt: Date,
    creator: {
        id: Number
    },
    editedAt:  Date
})


module.exports = mongoose.model("Channel", channelSchema);