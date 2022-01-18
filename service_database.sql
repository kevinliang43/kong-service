CREATE DATABASE kong_database;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE services (
    id UUID DEFAULT uuid_generate_v4(),
    service_id SERIAL,
    name VARCHAR(128),
    description VARCHAR(1024),
    version DECIMAL,
    PRIMARY KEY (id),
    UNIQUE(service_id, version)
);

CREATE INDEX service_version_index
    ON services (service_id, version);

-- initial mock data to populate the database with new services
INSERT INTO services (name, description, version) VALUES 
    ('Locate Us', 'Service for retrieving location info', '1.0'),
    ('Collect Monday', '', '1.0'),
    ('Contact Us', 'Service for retrieving contact us info', '1.0');

-- sample existing service with new version
INSERT INTO services (service_id, name, description, version) VALUES 
    ('1', 'Locate Us', 'Service for retrieving location info v1.1', '1.1');
