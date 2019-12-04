const router = require("express").Router();
const MongoClient = require("../models/mongoClient");

/**
 * Checks if user is the creator of the specific message.
 * If not, throws a 403 Forbidden error.
 */
router.use("/:messageID", function(req, res, next) {
  try {
    let db = MongoClient.getDB();
    db.collection("message").findOne({"id": parseInt(req.params.messageID)}
      , function(err, doc) {
        if (err) {
          console.log(err);
          return res.status(400).json(err);
        }
        if (doc.creator !== req.headers['x-user']) {
          return res.status(403).json("Forbidden.");
        }
        next();
      });
  } catch (err) {
    console.log(err);
    return res.status(400).json(err);
  }
});

/**
 * Updates specific message's body and returns
 * updated document.
 */
router.patch("/:messageID", async function (req, res, next) {
  try {
    let db = MongoClient.getDB();
    await db.collection("message").findOneAndUpdate(
      {id: parseInt(req.params.messageID)},
      {$set: {body:req.body.body}},
      { returnOriginal: false },
      function (err, doc) {
        if (err) {
          console.log(err);
          return res.status(400).json(err);
        }
        db.collection("channel").findOne(
          {"id": doc.channelID},
          function(err, channel) {
            const updateMessageNotify = {
              type: "message-update",
              message: doc
            };
            if (channel.private) {
              updateMessageNotify.userIDs = channel.members;
            }
            RabbitMQ.publish(updateMessageNotify);
          });
        return res.status(200).json(doc);
      });
  } catch (err) {
    console.log(err);
    return res.status(400).json(err);
  }
});

/**
 * Deletes specific message.
 */
router.delete("/:messageID", function (req, res, next) {
  try {
    let db = MongoClient.getDB();
    db.collection("message").remove(
      { "id": parseInt(req.params.messageID) },
      function (err, doc) {
        if (err) {
          return res.status(400).json(err);
        }
        db.collection("channel").findOne(
          {"id": doc.channelID},
          function(err, channel) {
            const deleteMessageNotify = {
              type: "message-delete",
              messageID: doc.id
            };
            if (channel.private) {
              deleteMessageNotify.userIDs = channel.members;
            }
            RabbitMQ.publish(deleteMessageNotify);
          });
        return res.status(200)
          .send(`Successfully deleted message ${req.params.messageID}.`);
      });
  } catch (err) {
    return res.status(400).json(err);
  }
});

module.exports = router;
