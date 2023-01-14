CREATE TABLE users (
    id integer not null primary key auto_increment,
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE todos (
    id  integer not null primary key auto_increment,
    user_id integer not null,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    foreign key (user_id) references users (id) on delete cascade
);
