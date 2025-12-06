create table tokens (
    id int not null auto_increment primary key,
    created_at datetime not null,
    deleted_at datetime default null,
    item_id int not null,
    token varchar(32) not null,
    expires_at datetime null,
    unique key (token)
);

create table publishes (
    id int not null auto_increment primary key,
    created_at datetime not null,
    token_id int not null,
    note text not null,
    ip_address varchar(45) not null,
    user_agent varchar(200) not null
);
