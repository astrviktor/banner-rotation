--CREATE DATABASE banner_rotation;
--CREATE USER "user" WITH ENCRYPTED PASSWORD 'password';
--GRANT ALL PRIVILEGES ON DATABASE banner_rotation TO "user";

CREATE SCHEMA banner_rotation;

CREATE TABLE banner_rotation.banner (
  id uuid NOT NULL unique,
  description text NOT NULL
);

CREATE INDEX banner_id_index ON banner_rotation.banner (id);
CREATE INDEX banner_description_index ON banner_rotation.banner (description);

CREATE TABLE banner_rotation.slot (
  id uuid NOT NULL unique,
  description text NOT NULL
);

CREATE INDEX slot_id_index ON banner_rotation.slot (id);
CREATE INDEX slot_description_index ON banner_rotation.slot (description);

CREATE TABLE banner_rotation.segment (
  id uuid NOT NULL unique,
  description text NOT NULL
);

CREATE INDEX segment_id_index ON banner_rotation.segment (id);
CREATE INDEX segment_description_index ON banner_rotation.segment (description);

CREATE TABLE banner_rotation.rotation (
  slot_id uuid NOT NULL,
  banner_id uuid NOT NULL,
  PRIMARY KEY (slot_id, banner_id)
);

CREATE INDEX rotation_slot_id_index ON banner_rotation.rotation (slot_id);
CREATE INDEX rotation_banner_id_index ON banner_rotation.rotation (banner_id);

CREATE TABLE banner_rotation.stat (
  banner_id uuid NOT NULL,
  segment_id uuid NOT NULL,
  show_count  integer,
  click_count integer,
  PRIMARY KEY (banner_id, segment_id)
);

CREATE INDEX stat_banner_id_index ON banner_rotation.stat (banner_id);
CREATE INDEX stat_segment_id_index ON banner_rotation.stat (segment_id);

CREATE TYPE action_type AS ENUM ('show', 'click');

CREATE TABLE banner_rotation.event (
  slot_id uuid NOT NULL,
  banner_id uuid NOT NULL,
  segment_id uuid NOT NULL,
  action action_type,
  date timestamp with time zone NOT NULL
);
