CREATE TABLE IF NOT EXISTS events
(
    id SERIAL PRIMARY KEY,
    mq_event_id VARCHAR(255) UNIQUE,
    type VARCHAR(30),
    confirmations_amount BIGINT,
    confirmers_amount BIGINT
);

CREATE TABLE IF NOT EXISTS compensatory_events
(
    id SERIAL PRIMARY KEY,
    mq_event_id VARCHAR(255) UNIQUE,
    event_id INT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products
(
    id SERIAL PRIMARY KEY, 
    user_id int not null, 
    name varchar(255) not null, 
    price BIGINT not null, 
    description varchar(255) not null, 
    amount BIGINT not null, 
    category varchar(255) not null,
    photo_url varchar(255),
    status VARCHAR(20) CHECK (status IN ('available', 'deleted')) DEFAULT 'available',
    status_event_id VARCHAR(255)
);

ALTER TABLE compensatory_events ADD FOREIGN KEY (event_id) REFERENCES events(id);
