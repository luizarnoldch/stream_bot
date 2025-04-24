-- reset_master.sql
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE datname = 'chatbot' AND pid <> pg_backend_pid();

DROP DATABASE IF EXISTS chatbot;
CREATE DATABASE chatbot;