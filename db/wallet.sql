CREATE DATABASE wallet
    WITH 
    OWNER = pgx
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    CONNECTION LIMIT = -1;
\c wallet;
CREATE SCHEMA  wallet;
CREATE TABLE wallet (idx SERIAL NOT NULL, funds REAL NOT NULL CHECK (funds >= 0), uuid CHAR(36) NOT NULL, owner_id INT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, last_modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (idx), UNIQUE (owner_id));
ALTER TABLE wallet
    OWNER to pgx;

CREATE TABLE transaction (idx SERIAL NOT NULL, sum REAL NOT NULL, ref varchar(255) ,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, last_modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (idx));
ALTER TABLE transaction
    OWNER to pgx;


CREATE USER pgx_test WITH PASSWORD 'test_user_2021';
CREATE DATABASE wallet_test
    WITH 
    OWNER = pgx_test
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    CONNECTION LIMIT = -1;
\c wallet_test;
CREATE TABLE wallet (idx SERIAL NOT NULL, funds REAL NOT NULL CHECK (funds >= 0), uuid CHAR(36) NOT NULL, owner_id INT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, last_modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (idx), UNIQUE (owner_id));
ALTER TABLE wallet
    OWNER to pgx_test;
CREATE TABLE transaction (idx SERIAL NOT NULL, sum REAL NOT NULL, ref varchar(255) ,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, last_modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (idx));
ALTER TABLE transaction
    OWNER to pgx_test;
