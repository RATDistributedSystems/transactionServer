CREATE KEYSPACE userdb WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  AND durable_writes = true;

CREATE TABLE userdb.triggers (
    tid uuid PRIMARY KEY,
    operation boolean,
    pendingcash int,
    pendingstocks int,
    stock text,
    stockamount int,
    triggervalue int,
    userid text
);

CREATE TABLE userdb.users (
    userid text PRIMARY KEY,
    pendingcash int,
    stock int,
    usablecash int
);

CREATE TABLE userdb.userstocks (
    usid uuid,
    stockamount int,
    stock text,
    userid text,
    PRIMARY KEY((userid), stock)
);

CREATE TABLE userdb.buypendingtransactions (
    pid timeuuid,
    pendingcash int,
    stock text,
    stockvalue int,
    userid text,
    PRIMARY KEY((userid), pid)
);

CREATE TABLE userdb.sellpendingtransactions (
    pid timeuuid,
    pendingcash int,
    stock text,
    stockvalue int,
    userid text,
    PRIMARY KEY((userid), pid)
);
    
create table userdb.buyTriggers (
    tid uuid, 
    userId varchar, 
    stock varchar, 
    pendingCash int, 
    triggerValue int, 
    stockAmount int, 
    pendingStocks int, 
    primary key (userId, stock));

create table userdb.sellTriggers (
    tid uuid, 
    userId varchar, 
    stock varchar, 
    pendingCash int, 
    triggerValue int, 
    stockAmount int, 
    pendingStocks int, 
    primary key (userId, stock));