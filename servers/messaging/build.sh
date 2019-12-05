use info441-a5
db.createCollection("channel")
db.createCollection("message")
db.channel.createIndex( { "email" : 1 }, { unique : true } )