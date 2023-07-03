-- Script to create the test role, database, schema, and tables.
-- The first few entries may be easier done in a GUI (eg pgAdmin).

-- USER

CREATE ROLE example WITH
  LOGIN
  NOSUPERUSER
  NOCREATEDB
  NOCREATEROLE
  INHERIT
  NOREPLICATION
  CONNECTION LIMIT -1;

ALTER ROLE example
	PASSWORD 'example';

-- DATABASE

CREATE DATABASE example
    WITH
    OWNER = example
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

COMMENT ON DATABASE example
    IS 'An example database for the NearGothic Go tool.';


-- SCHEMA

CREATE SCHEMA IF NOT EXISTS example
  AUTHORIZATION example;



-- TABLE (account)

CREATE TABLE example.account
(
  id bigserial NOT NULL,
  email_address character varying(250) NOT NULL,
  display_name character varying(50) NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  deleted_at timestamp with time zone,
  PRIMARY KEY (id),
  CONSTRAINT uniq_account_email_address UNIQUE (email_address)
);

ALTER TABLE IF EXISTS example.account
  OWNER to example;

COMMENT ON TABLE example.account
    IS 'A table of user accounts.';

COMMENT ON COLUMN example.account.id
    IS 'The unique account ID.';

COMMENT ON COLUMN example.account.email_address
    IS 'The account-holder''s contact email address.';

COMMENT ON COLUMN example.account.display_name
    IS 'The account-holder''s display name.';

COMMENT ON COLUMN example.account.created_at
    IS 'When the account was created.';

COMMENT ON COLUMN example.account.updated_at
    IS 'When the account details were last updated.';

COMMENT ON COLUMN example.account.deleted_at
    IS 'When (if) the account was deleted.';


-- TABLE (setting)

CREATE TABLE example.setting
(
  id bigserial NOT NULL,
  display_name character varying(50) NOT NULL,
  details character varying(500) NOT NULL,
  max_value_length bigint NOT NULL DEFAULT 30,
  is_enabled boolean NOT NULL,
  PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS example.setting
  OWNER to example;

COMMENT ON TABLE example.setting
    IS 'The settings available for an account.';

COMMENT ON COLUMN example.setting.id
    IS 'The unique setting ID.';

COMMENT ON COLUMN example.setting.display_name
    IS 'The displayable brief name for this setting.';

COMMENT ON COLUMN example.setting.details
    IS 'Descriptive details for this setting.';

COMMENT ON COLUMN example.setting.max_value_length
    IS 'The longest a value for this setting is allowed to be.';


-- TABLE (account_setting)

CREATE TABLE example.account_setting
(
  id bigserial NOT NULL,
  account_id bigint NOT NULL,
  setting_id bigint NOT NULL,
  value character varying(250) NOT NULL DEFAULT '',
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id),
  CONSTRAINT fk_account_setting_account FOREIGN KEY (account_id)
    REFERENCES example.account (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID,
  CONSTRAINT fk_account_setting_setting FOREIGN KEY (setting_id)
    REFERENCES example.setting (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
);

ALTER TABLE IF EXISTS example.account_setting
  OWNER to example;

COMMENT ON TABLE example.account_setting
    IS 'Settings for a particular account.';

COMMENT ON COLUMN example.account_setting.id
    IS 'The unique ID for this setting for this account.';

COMMENT ON COLUMN example.account_setting.account_id
    IS 'The account this setting''s value applies to.';

COMMENT ON COLUMN example.account_setting.setting_id
    IS 'The ID of the setting which has this value.';

COMMENT ON COLUMN example.account_setting.value
    IS 'The current value for this account setting.';

COMMENT ON COLUMN example.account_setting.updated_at
    IS 'When the value was last updated.';

COMMENT ON CONSTRAINT fk_account_setting_account ON example.account_setting
    IS 'Link to the account this value is for.';
COMMENT ON CONSTRAINT fk_account_setting_setting ON example.account_setting
    IS 'Link to the setting this value is for.';


-- VIEW (active_account)

CREATE OR REPLACE VIEW example.active_account
 AS
 SELECT a.id AS account_id,
    a.email_address,
    a.display_name
   FROM example.account a
  WHERE a.deleted_at IS NULL;

ALTER TABLE example.active_account
    OWNER TO example;
COMMENT ON VIEW example.active_account
    IS 'Sample view listing basic details for accounts NOT soft-deleted.';
