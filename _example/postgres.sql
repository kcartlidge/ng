/*
WARNINGS

This is a fallback script only, NOT a structural database backup.
The SQL below only represents what was needed to generate code
and is intended solely for emergency use when other backups fail.
Schemas created from it may omit things (eg sequences, functions).

The script should NOT be simply executed in one go!
Several deliberate restrictions *force* you to take care:
  - DROP statements assume things already exist
  - Tables are in RANDOM order, NOT in order of dependencies

DATABASE SETUP

You'll need a Postgres server and a database.
Here's some example SQL to create a suitable login and schema in Postgres.
Remember to CHANGE THE PASSWORD.
The details should be the same as the ones in your connection string
environment variable (the last environment variable name used was
'DB_CONNSTR').

CREATE USER example WITH PASSWORD 'example';
CREATE DATABASE example
WITH
    ENCODING = 'UTF8'
    OWNER = example
    CONNECTION LIMIT = 100;
CREATE SCHEMA example AUTHORIZATION example;
*/




-------- Account --------

DROP TABLE example.account CASCADE;

CREATE TABLE IF NOT EXISTS example.account (
    id  BIGSERIAL NOT NULL,
    email_address  character varying(250) NOT NULL,
    display_name  character varying(50) NOT NULL,
    created_at  timestamp with time zone NOT NULL,
    updated_at  timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at  timestamp with time zone,
    CONSTRAINT account_pkey PRIMARY KEY (id),
    CONSTRAINT uniq_account_email_address UNIQUE (email_address)
);

ALTER TABLE example.account OWNER TO example;
COMMENT ON TABLE example.account IS 'A table of user accounts.';
COMMENT ON COLUMN example.account.id IS 'The unique account ID.';
COMMENT ON COLUMN example.account.email_address IS 'The account-holder''s contact email address.';
COMMENT ON COLUMN example.account.display_name IS 'The account-holder''s display name.';
COMMENT ON COLUMN example.account.created_at IS 'When the account was created.';
COMMENT ON COLUMN example.account.updated_at IS 'When the account details were last updated.';
COMMENT ON COLUMN example.account.deleted_at IS 'When (if) the account was deleted.';


-------- Account Setting --------

DROP TABLE example.account_setting CASCADE;

CREATE TABLE IF NOT EXISTS example.account_setting (
    id  BIGSERIAL NOT NULL,
    account_id  bigint NOT NULL,
    setting_id  bigint NOT NULL,
    value  character varying(250) NOT NULL DEFAULT ''::character varying,
    updated_at  timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT account_setting_pkey PRIMARY KEY (id),
    CONSTRAINT fk_account_setting_account FOREIGN KEY (account_id)
        REFERENCES example.account (id) MATCH SIMPLE 
        ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT fk_account_setting_setting FOREIGN KEY (setting_id)
        REFERENCES example.setting (id) MATCH SIMPLE 
        ON UPDATE NO ACTION ON DELETE NO ACTION
);

ALTER TABLE example.account_setting OWNER TO example;
COMMENT ON TABLE example.account_setting IS 'Settings for a particular account.';
COMMENT ON COLUMN example.account_setting.id IS 'The unique ID for this setting for this account.';
COMMENT ON COLUMN example.account_setting.account_id IS 'The account this setting''s value applies to.';
COMMENT ON COLUMN example.account_setting.setting_id IS 'The ID of the setting which has this value.';
COMMENT ON COLUMN example.account_setting.value IS 'The current value for this account setting.';
COMMENT ON COLUMN example.account_setting.updated_at IS 'When the value was last updated.';


-------- Setting --------

DROP TABLE example.setting CASCADE;

CREATE TABLE IF NOT EXISTS example.setting (
    id  BIGSERIAL NOT NULL,
    display_name  character varying(50) NOT NULL,
    details  character varying(500) NOT NULL,
    max_value_length  bigint NOT NULL DEFAULT 30,
    is_enabled  boolean NOT NULL,
    CONSTRAINT setting_pkey PRIMARY KEY (id)
);

ALTER TABLE example.setting OWNER TO example;
COMMENT ON TABLE example.setting IS 'The settings available for an account.';
COMMENT ON COLUMN example.setting.id IS 'The unique setting ID.';
COMMENT ON COLUMN example.setting.display_name IS 'The displayable brief name for this setting.';
COMMENT ON COLUMN example.setting.details IS 'Descriptive details for this setting.';
COMMENT ON COLUMN example.setting.max_value_length IS 'The longest a value for this setting is allowed to be.';
