
create database if not EXISTS psql_db;
  use psql_db;
create table if not exists contacts (
    id SERIAL primary key,
    email varchar(320) not null UNIQUE,
    passhash varchar(255) not null,
    user_name varchar(255) not null UNIQUE,
    first_name varchar(128) not null,
    last_name varchar(128) not null,
    photo_url varchar(255) not null
);
create table if not exists sessions (
    id SERIAL primary key,
    sign_in_time timestamp not null,
    ip varchar(15) not null UNIQUE
);

/*
  2 differences from the Go model.
    channel members are not a field but fetched from the channel_members join table.
    creator is a integer value corresponding to the contact id of the creator.
*/
create table if not exists channels (
    id SERIAL primary key,
    name varchar(320) not null UNIQUE,
    description varchar(255) not null,
    private BOOLEAN not null,
    createdAt timestamp not null,
    creator Int not null,
    editedAt timestamp not null
);
create table if not exists messages(
  id SERIAL primary key,
  channelId Int not null,
  body varchar(5000),
  createdAt timestamp not null,
  creator Int not null,
  editedAt timestamp not null
);
create table if not exists channelMembers (
    channelId int,
    contactId int,
    Primary Key(channelId, contactId)
);
DELIMITER //
DROP PROCEDURE IF EXISTS getChannelMembers;
CREATE PROCEDURE getChannelMembers(IN targetChannelId int)
BEGIN
  SELECT contactId FROM channel_members
  WHERE channelId = targetChannelId;
END //
DELIMITER ;

INSERT INTO contacts (email, passhash, user_name, first_name,last_name,photo_url)
VALUES ("mail", "passhash", "leoTran", "Leo", "Tran", "sss");
