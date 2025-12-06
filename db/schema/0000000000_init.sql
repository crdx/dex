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
    last_hit datetime null,
    uuid text not null,
    kind text not null,
    blob_hash text not null
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
