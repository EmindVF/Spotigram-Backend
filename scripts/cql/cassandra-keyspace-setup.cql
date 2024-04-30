CREATE KEYSPACE IF NOT EXISTS spotigram 
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };