create table blobs (
    id int not null auto_increment primary key,
    created_at datetime not null,
    deleted_at datetime default null,
    hash text not null,
    content longblob not null
);

create table items (
    id int not null auto_increment primary key,
    created_at datetime not null,
    deleted_at datetime default null,
    label text not null,
    hits int not null,
    last_hit_at datetime null,
    uuid text not null,
    kind text not null,
    blob_hash text not null,
    content_type text not null
);

create table logs (
    id int not null auto_increment primary key,
    created_at datetime not null,
    item_id int not null,
    method varchar(10) not null,
    request varchar(200) not null,
    ip_address varchar(15) not null,
    user_agent varchar(200) not null
);

create table tokens (
    id int not null auto_increment primary key,
    created_at datetime not null,
    deleted_at datetime default null,
    item_id int not null,
    token varchar(48) not null,
    expires_at datetime null,
    unique key (token)
);

create table deployments (
    id int not null auto_increment primary key,
    created_at datetime not null,
    token_id int not null,
    note text not null,
    ip_address varchar(45) not null,
    user_agent varchar(200) not null
);
