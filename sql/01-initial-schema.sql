CREATE TABLE IF NOT EXISTS users (
    id integer not null primary key auto_increment,
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS todos (
    id  integer not null primary key auto_increment,
    user_id integer not null,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    done BOOL DEFAULT FALSE,
    deadline TIMESTAMP NULL DEFAULT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    foreign key (user_id) references users (id) on delete cascade
);
