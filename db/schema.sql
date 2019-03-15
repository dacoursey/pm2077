DROP TABLE IF EXISTS system_notifications;
DROP TABLE IF EXISTS internal_resources;
DROP TABLE IF EXISTS task;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;


CREATE TABLE IF NOT EXISTS roles
(
    id serial primary key,
    name varchar(40)
);

INSERT INTO roles (name) VALUES ('Admin');
INSERT INTO roles (name) VALUES ('Project Manager');
INSERT INTO roles (name) VALUES ('User');



CREATE TABLE IF NOT EXISTS users
(
    id serial primary key,
    username varchar(100) not null,
    password varchar(100) not null,
    role int not null references roles(id)
);

INSERT INTO users (username, password, role) VALUES ('admin','password','1');
INSERT INTO users (username, password, role) VALUES ('dave','password','3');
INSERT INTO users (username, password, role) VALUES ('pm','password','2');


CREATE TABLE IF NOT EXISTS customers
(
    id serial primary key,
    name varchar(100) not null,
    poc int 
);

INSERT INTO customers (name, poc) VALUES ('Tesla Motor Company','1');
INSERT INTO customers (name, poc) VALUES ('Lindner Salt Factory, Inc.','4');


CREATE TABLE IF NOT EXISTS contacts
(
    id serial primary key,
    cust_id int not null,
    fname varchar(25) not null,
    lname varchar(25) not null,
    email varchar(50) not null
);

INSERT INTO contacts (cust_id, fname, lname, email) VALUES ('1','Scooby','Doo','scooby@gmail.com');
INSERT INTO contacts (cust_id, fname, lname, email) VALUES ('1','George','Washington','originalgw@gmail.com');
INSERT INTO contacts (cust_id, fname, lname, email) VALUES ('1','John','Doe','doej@aol.com');
INSERT INTO contacts (cust_id, fname, lname, email) VALUES ('2','Steve','Jobs','ceo@apple.com');
INSERT INTO contacts (cust_id, fname, lname, email) VALUES ('2','Other','Dullah','blah@gmail.com');


CREATE TABLE IF NOT EXISTS projects
(
    id serial primary key,
    cust_id int not null references customers(id),
    name varchar(100) not null,
    start_date date,
    hours int
);

INSERT INTO projects (name, cust_id, start_date, hours) VALUES ('Public Website','1','2017-06-01','40');
INSERT INTO projects (name, cust_id, start_date, hours) VALUES ('Auto Software','1','2017-06-01','40');
INSERT INTO projects (name, cust_id, start_date, hours) VALUES ('Saltiness Magnifyer','1','2017-06-01','40');


CREATE TABLE IF NOT EXISTS system_notifications
(
    id serial primary key,
    title varchar(50) not null,
    message varchar(150) not null
);

INSERT INTO system_notifications (title, message) VALUES ('The 3rd Floor ice cream machine is down!', 'Third floor ice cream machine is out of order.  Again!!');


CREATE TABLE IF NOT EXISTS internal_resources
(
    id serial primary key,
    title varchar(50) not null,
    url varchar(150) not null
);

INSERT INTO internal_resources (title, url) VALUES ('Standard Form 42', 'http://sharepoint/forms/42');
INSERT INTO internal_resources (title, url) VALUES ('How to submit a bug report.', 'http://sharepoint/howto/submitbugreport');


CREATE TABLE IF NOT EXISTS task
(
    id serial primary key,
    user_id int not null,
    project_id int not null,
    title varchar(150) not null,
    is_completed int not null,
    attachment_path varchar(1000),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);
