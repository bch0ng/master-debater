const router = require("express").Router();
const MongoClient = require("../models/mongoClient");
const RabbitMQ = require("../controllers/amqp");

//Returns a list of visible channels
function getListOfVisibleChannels(channelsArray, memberName) {
  return channelsArray.filter((channel) => {
    return (!channel.private || isMemberOfChannel(channel.members, memberName))
  });
}

//Returns the a timestamp as an iso string
function getCurrentTime() {
  let newDateStr=new Date().toISOString();
  return newDateStr;
}

//Returns a list of all members of a channel
function isMemberOfChannel(members, memberName) {
  let membersLength = members.length;
  for (let i = 0; i < membersLength; i++) {
    let targetMember = members[i];
    if (memberName == targetMember.name) {
      return true;
    }
  }
  return false;
}

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
  let memberName = req.body.memberID;
  mongoDb.collection('channel')
    .find()
    .toArray(function(err, doc) {
      if (err) {
        return res.status(400).json(err);
      }
      let visibleChannels = getListOfVisibleChannels(doc, memberName).map((channel) => {
        return channel.name;
      });
      return res.json(visibleChannels);
    });
});

// Makes a new channel
router.post("/", async function (req, res, next) {
  console.log(req.body);
  if (req.body.name) {
    let description = req.body.description;
    const insertID = await getNextSequenceValue('channelCounters', "counterID");
    let newEntry = {
      id: insertID,
      name: req.body.name,
      description: "",
      createdAt: new Date(),
      creator: JSON.parse(req.headers['x-user']),
      editedAt: new Date()
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
        if (doc.ops[0].private) {
          newChannelNotify.userIDs = doc.ops[0].members;
        }
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

// Returns last 100 messages for specified channel.
// If the channel is private, the requesting user must be
// a member of the channel.
router.get("/:channelID", function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID) {
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
        if (newChannel && (!newChannel.private || isMemberOfChannel(newChannel, req.body.memberID))) {
          let messages = mongoDb.collection("message");
          const mongoQuery = {
            "channelID": parseInt(attemptedChannelID)
          };
          if (req.query.before) {
            mongoQuery.id = {
              $lte: req.query.before
            };
          }
          messages.find(mongoQuery)
            .sort({ createdAt: -1 })
            .limit(100)
            .toArray(function(err, docs) {
              if (err) {
                return res.status(400).json(err);
              }
              return res.json(docs);
            });
        } else {
          return res.status(403).send({
            message: 'forbidden!'
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
// Post a message to a specific channel.
// If the channel is private, the requesting user must be
// a member of the channel.
router.post("/:channelID", function (req, res, next) {

  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID) {
    let mongoDb = MongoClient.getDB();
    try {
      let channels = mongoDb.collection("channel");
      channels.findOne({
        id: parseInt(attemptedChannelID)
      }, async function(err, newChannel) {
        if (err) {
          return res.status(400).json(err);
        }
        if (newChannel && (!newChannel.private || isMemberOfChannel(newChannel, req.body.memberID))) {
          const messageID = await getNextSequenceValue('messageCounters', "counterID");
          let messages = mongoDb.collection("message");
          let updatedEntry = {
            id: messageID,
            channelID: newChannel.id,
            body: req.body.body,
            createdAt: getCurrentTime(),
            creator: req.body.member,
            editedAt: getCurrentTime()
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
// Update specific channel's name and description.
// If the channel is private, the requesting user must be
// a member of the channel.
router.patch("/:channelID", function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID) {
    let mongoDb = MongoClient.getDB();
    try {
      let channels = mongoDb.collection("channel");
      channels.findOne({
        id: parseInt(attemptedChannelID)
      }, function(err, newChannel) {
        if (newChannel && (!newChannel.private || isMemberOfChannel(newChannel, req.body.memberID))) {
          let updatedEntry = {
            name: req.body.name,
            description: req.body.description
          };
          channels.findOneAndUpdate(
            { id: parseInt(attemptedChannelID) },
            { $set: updatedEntry },
            { returnOriginal: false },
            function(err, doc) {
              if (err) {
                return res.status(400).json(err);
              }
              const updateChannelNotify = {
                type: "channel-update",
                channel: doc
              };
              if (doc.private) {
                updateChannelNotify.userIDs = doc.members;
              }
              RabbitMQ.publish(updateChannelNotify);
              return res.json(doc);
            });
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
// Deletes a specific channel if the requesting user
// is the creator of the channel.
router.delete("/:channelID", function (req, res, next) {
  let attemptedChannelID = req.params.channelID;
  if (attemptedChannelID) {
    let mongoDb = MongoClient.getDB();
    try {
      let channels = mongoDb.collection("channel");
      channels.findOne({
        id: parseInt(attemptedChannelID)
      }, function(err, newChannel) {
        if (newChannel && newChannel.creator.id == req.body.memberID) {
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

// Adds given member to the specified channel if the requesting user
// is the creator of the channel.
router.post("/:channelID/members", function (req, res, next) {
  let channels = mongoDb.collection("channel");
  let channel = channels.find({
    id: parseInt(req.params.channelID)
  });
  if (channel.creator != req.body.memberID) {
    return res.status(403).json("Forbidden");
  }
  const members = channel.members;
  members.push(req.body.newMemberID);
  channels.findOneAndUpdate({id: parseInt(req.params.channelID)},
    { $set:{"members": members} },
    function(err, doc){
      if (err) {
        return res.status(400).json(err);
      }
      res.setHeader('content-type', 'text/plain');
      return res.status(201).send("User was added to the channel.")
    });
});
// Removes a member from the specific channel if the requesting user
// is the creator of the channel.
router.delete("/:channelID/members", function (req, res, next) {
  let channels = mongoDb.collection("channel");
  let channel = channels.findOne({
    id: parseInt(req.params.channelID)
  }, function(err, channel) {
    if (channel.creator != req.body.memberID) {
      return res.status(403).json("Forbidden");
    }z
  });
  channels.findOneAndUpdate({id: parseInt(req.params.channelID)},
    { $pull:{"members": req.body.deleteMemberID} },
    function(err, doc){
      if (err) {
        return res.status(400).json(err);
      }
      res.setHeader('content-type', 'text/plain');
      return res.status(201).send("User was deleted from the channel.")
    });
});

module.exports = router;
