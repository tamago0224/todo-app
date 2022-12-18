CREATE TABLE todos (
    id  integer not null primary key auto_increment,
    title VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE users (
    id integer not null primary key auto_increment,
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
)
