CREATE TABLE IF NOT EXISTS products
(
    id SERIAL PRIMARY KEY, 
    user_id int not null, 
    name varchar(255) not null, 
    price BIGINT not null, 
    description varchar(255) not null, 
    amount BIGINT not null, 
    category varchar(255) not null 
);