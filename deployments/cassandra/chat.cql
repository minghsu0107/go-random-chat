CREATE KEYSPACE randomchat WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
USE randomchat;
CREATE TABLE channels (
    id varint,
    user_id varint,
    PRIMARY KEY((id), user_id)
);
CREATE TABLE messages (
    id varint,
    event int,
    channel_id varint,
    user_id varint,
    payload text,
    seen boolean,
    timestamp timestamp,
    PRIMARY KEY((channel_id), id)
) WITH CLUSTERING ORDER BY (id DESC);