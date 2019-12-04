let error = true;

let res = [
    db.createCollection("channel"),
    db.channel.createIndex( { "username" : 1 }, { unique : true } ),
    db.createCollection("channelCounters"),
    db.channelCounters.createIndex( {"seq": 1 }),
    db.channelCounters.insert({id:"counterID", seq:0}),
    db.createCollection("message"),
    db.createCollection("messageCounters"),
    db.messageCounters.createIndex( {"seq": 1 }),
    db.messageCounters.insert({id:"counterID", seq:0})
];

printjson(res);

if (error) {
  print('Error, exiting');
  quit(1);
}