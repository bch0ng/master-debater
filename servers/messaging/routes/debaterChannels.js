const router = require("express").Router();
const MongoClient = require("../models/mongoClient");
const RabbitMQ = require("../controllers/amqp");

// Returns a list of all channels a member can see
router.get("/", function (req, res, next) {
  console.log("HERE!");
  let mongoDb = MongoClient.getDB();
  let memberName = req.body.memberID;
  mongoDb.collection('channel')
    .find()
    .toArray(function(err, doc) {
      if (err) {
        return res.status(400).json(err);
      }
      let visibleChannels = doc.map((channel) => {
        return {
            id: channel.id,
            channel: channel.name
        };
      });
      return res.json(visibleChannels);
    });
});

module.exports = router;
