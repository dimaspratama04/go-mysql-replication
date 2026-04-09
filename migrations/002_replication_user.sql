-- Create dedicated replication user on Primary
-- This file runs automatically via docker-entrypoint-initdb.d

USE mysql;

-- Create replication user
CREATE USER IF NOT EXISTS 'replicator'@'%' IDENTIFIED WITH mysql_native_password BY 'replicatorpass';
GRANT REPLICATION SLAVE ON *.* TO 'replicator'@'%';
FLUSH PRIVILEGES;
