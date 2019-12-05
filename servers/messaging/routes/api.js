const router = require("express").Router();
const MongoClient = require("../models/mongoClient");

const channels = require("./channels");
const messages = require("./messages");

router.use("/channels", channels);
router.use("/messages", messages);

module.exports = router;