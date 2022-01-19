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
INSERT INTO services (id, service_id, name, description, version) VALUES 
    ('53ED660E-24BE-4662-9564-F7464991D651', '1', 'Locate Us', 'Service for retrieving location info', '1.0'),
    ('AD187F14-B774-4D25-A939-864AA57903A0', '2', 'Collect Monday', '', '1.0'),
    ('1D75DF8D-1639-428D-AAD3-C78CD71A250F', '3', 'Contact Us', 'Service for retrieving contact us info', '1.0');

-- sample existing service with new version
INSERT INTO services (id, service_id, name, description, version) VALUES 
    ('6ACB5E21-DBC5-4ED7-83E3-EDE75A3255D1', '1', 'Locate Us', 'Service for retrieving location info v1.1', '1.1');



CREATE TABLE services_latest (
    service_id INTEGER,
    latest_record_id UUID,
    name VARCHAR(128),
    description VARCHAR(1024),
    version DECIMAL,
    versions INTEGER,
    UNIQUE(service_id, latest_record_id)
);

INSERT INTO services_latest (service_id, latest_record_id, name, description, version, versions) VALUES
    ('1', '6ACB5E21-DBC5-4ED7-83E3-EDE75A3255D1', 'Locate Us', 'Service for retrieving location info v1.1', '1.1', '2'),
    ('2', 'AD187F14-B774-4D25-A939-864AA57903A0', 'Collect Monday', '', '1.0', '1'),
    ('3', '1D75DF8D-1639-428D-AAD3-C78CD71A250F', 'Contact Us', 'Service for retrieving contact us info', '1.0', '1');
