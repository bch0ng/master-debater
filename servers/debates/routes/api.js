const router = require("express").Router();
const MongoClient = require("../models/mongoClient");

const chatroom = require("./chatrooms");
const messages = require("./messages");

router.use("/chatroom", chatroom);
router.use("/messages", messages);

module.exports = router;