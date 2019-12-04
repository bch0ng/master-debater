const router = require("express").Router();
const MongoClient = require("../models/mongoClient");
const RabbitMQ = require("../controllers/amqp");


async function getNextSequenceValue(dbName, name) {
  let db = MongoClient.getDB();
  try {
    const doc = await db.collection(dbName).findOneAndUpdate(
      { id: name },
      { $inc: { seq: 1 } },
      { returnOriginal: false });
    return doc.value.seq;
  } catch (err) {
    throw err;
  }
}


// Returns a list of all channels a member can see
router.get("/", function (req, res, next) {
  let mongoDb = MongoClient.getDB();
  mongoDb.collection('channel')
    .find()
    .toArray(function(err, doc) {
      if (err) {
        return res.status(400).json(err);
      }
      return res.json(doc);
    });
});

// Makes a new channel
router.post("/create", async function (req, res, next) {
  if (req.body.name && req.headers["x-user"]) {
    let description = req.body.description;
    const insertID = await getNextSequenceValue('channelCounters', "counterID");
    let newEntry = {
      id: insertID,
      name: req.body.name,
      description: "",
      creator: JSON.parse(req.headers["x-user"]),
      debaters: [
        JSON.parse(req.headers["x-user"])
      ],
      voted: [],
      votes: 0,
      createdAt: new Date()
    };
    if (description) {
      newEntry.description = description;
    }
    let mongoDb = MongoClient.getDB();
    try {
      let channels = mongoDb.collection("channel");
      channels.insertOne(newEntry, function(err, doc) {
        if (err) {
          return res.status(400).json(err);
        }
        const newChannelNotify = {
          type: "channel-new",
          channel: doc.ops[0],
        };
        RabbitMQ.publish(newChannelNotify);
        return res.json(doc.ops[0]);
      });
    } catch (error) {
      return res.status(400).send({
        message: 'Error inserting channel!'
      });
    }
  } else {
    return res.status(400).send({
      message: 'Invalid channel name!'
    });
  }
});

router.post("/:channelID", function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID && req.headers["x-user"]) {
    let mongoDb = MongoClient.getDB();
    try {
      mongoDb.collection('channel').findOne({
        "id": parseInt(attemptedChannelID)
      }, function(err, newChannel) {
        if (err) {
          return res.status(400).json(err);
        } else if (newChannel === null) {
          return res.status(403).send({
            message: 'Error fetching channel!'
          });
        }
        const debaters = newChannel.debaters;
        if (debaters.length >= 10) {
          return res.status(200).json("Joined as audience.")
        } else {
          const currUserID = JSON.parse(req.headers["x-user"]).id;
          if (debaters.filter(function(debater) { return debater.id === currUserID; }).length > 0) {
            return res.status(400).json("Already joined as a debater.")
          } else {
            debaters.push(JSON.parse(req.headers["x-user"]));
          }
          mongoDb.collection('channel').findOneAndUpdate(
            { "id": parseInt(attemptedChannelID) },
            { $set: { debaters: debaters }},
            { returnOriginal: false },
            function(err, doc) {
              return res.status(200).json("Joined as a debater.");
            });
        }
      });
    } catch (error) {
      return res.status(403).send({
        message: 'Error fetching channel!'
      });
    }
  }
});

// req.body.vote +1 = for, -1 = against
router.post("/:channelID/vote", function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID && req.headers["x-user"]) {
    let mongoDb = MongoClient.getDB();
    try {
      mongoDb.collection('channel').findOne({
        "id": parseInt(attemptedChannelID)
      }, function(err, newChannel) {
        if (err) {
          return res.status(400).json(err);
        } else if (newChannel === null) {
          return res.status(403).send({
            message: 'Error fetching channel!'
          });
        }
        const voted = newChannel.voted;
        if (voted.filter(function(voter) { return voter.id === currUserID; }).length > 0) {
          return res.status(400).json("You have already voted for this debate.");
        } else {
          voted.push(JSON.parse(req.headers["x-user"]));
        }
        mongoDb.collection('channel').findOneAndUpdate(
          { "id": parseInt(attemptedChannelID) },
          {
            $inc: { votes: req.body.vote },
            $set: { voted: voted }
          },
          { returnOriginal: false },
          function(err, doc) {
            return res.status(200).json("Vote successfully entered!");
          });
      });
    } catch (error) {
      return res.status(403).send({
        message: 'Error fetching channel!'
      });
    }
  }
});

// Deletes a specific channel if the requesting user
// is the creator of the channel.
router.delete("/:channelID", function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID && req.headers["x-user"]) {
    let mongoDb = MongoClient.getDB();
    try {
      let channels = mongoDb.collection("channel");
      channels.findOne({
        id: parseInt(attemptedChannelID)
      }, function(err, newChannel) {
        if (newChannel && newChannel.creator == req.headers["x-user"]) {
          channels.remove(
            { id: parseInt(attemptedChannelID) },
          function(err, doc) {
            if (err) {
              return res.status(400).json(err);
            }
            const deleteChannelNotify = {
              type: "channel-delete",
              channelID: doc.id
            };
            if (doc.private) {
              deleteChannelNotify.userIDs = doc.members;
            }
            RabbitMQ.publish(deleteChannelNotify);
            return res.status(200).type("text").send('Delete successful!');
          })
        } else {
          return res.status(403).send({
            message: 'forbidden!'
          });
        }
      });
    } catch (error) {
      return res.status(400).send({
        message: 'Error posting message!'
      });
    }
  }
});

// Returns a list of all channels a member can see
router.get("/:channelID/posts", function (req, res, next) {
  let mongoDb = MongoClient.getDB();
  mongoDb.collection('message')
    .find({ channelID: parseInt(req.params.channelID) })
    .toArray(function(err, doc) {
      if (err) {
        return res.status(400).json(err);
      }
      return res.json(doc);
    });
});

/**
 * Updates specific message's body and returns
 * updated document.
 */
router.post("/:channelID/post", async function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID && req.headers["x-user"]) {
    let mongoDb = MongoClient.getDB();
    try {
      mongoDb.collection('channel').findOne({
        "id": parseInt(attemptedChannelID)
      }, async function(err, newChannel) {
        if (err) {
          return res.status(400).json(err);
        } else if (newChannel === null) {
          return res.status(403).send({
            message: 'Error fetching channel!'
          });
        }
        const debaters = newChannel.debaters;
        const currUserID = JSON.parse(req.headers["x-user"]).id;
        if (debaters.filter(function(debater) { return debater.id === currUserID; }).length > 0) {
          const messageID = await getNextSequenceValue('messageCounters', "counterID");
          let messages = mongoDb.collection("message");
          let updatedEntry = {
            id: messageID,
            channelID: parseInt(req.params.channelID),
            body: req.body.message,
            author: JSON.parse(req.headers["x-user"]),
            createdAt: new Date()
          };
          messages.insertOne(updatedEntry, function(err, doc) {
            if (err) {
              return res.status(400).json(err);
            }
            const newMessageNotify = {
              type: "message-new",
              message: doc.ops[0]
            };
            if (newChannel.private) {
              newMessageNotify.userIDs = newChannel.members;
            }
            RabbitMQ.publish(newMessageNotify);
            return res.json(doc.ops[0]);
          });
        } else {
          return res.status(400).json("You are not a debater.")
        }
      });
    } catch (error) {
      return res.status(403).send({
        message: 'Error fetching channel!'
      });
    }
  }
});

module.exports = router;
