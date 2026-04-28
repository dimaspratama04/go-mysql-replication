-- Create dedicated replication user on Primary
-- This file runs automatically via docker-entrypoint-initdb.d
USE mysql;

-- Create replication user
CREATE USER IF NOT EXISTS 'replicator'@'%' IDENTIFIED WITH caching_sha2_password BY 'replicatorpass';
GRANT REPLICATION SLAVE ON *.* TO 'replicator'@'%';

-- Create app user
CREATE USER IF NOT EXISTS 'appuser'@'%' IDENTIFIED WITH caching_sha2_password BY 'apppassword';
GRANT ALL PRIVILEGES ON *.* TO 'appuser'@'%';

-- Create Orchestrator user
CREATE USER 'orchestrator'@'%' IDENTIFIED WITH caching_sha2_password BY 'orchestratorpass';
GRANT SUPER, PROCESS, REPLICATION SLAVE, RELOAD ON *.* TO 'orchestrator'@'%';
GRANT SELECT ON mysql.* TO 'orchestrator'@'%';
GRANT ALL PRIVILEGES ON orchestrator.* TO 'orchestrator'@'%';
FLUSH PRIVILEGES;
