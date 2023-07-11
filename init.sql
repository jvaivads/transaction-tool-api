create table user
(
    id           int                                  not null auto_increment,
    user_name    varchar(100)                         not null,
    email        varchar(100)                         not null,
    date_created datetime default current_timestamp() not null,
    constraint user_pk primary key (id)
);

create table transaction
(
    id           int        not null auto_increment,
    user_id      int        not null,
    amount       float      not null,
    date_created datetime   not null,
    constraint transaction_pk primary key (id),
    constraint user_id_fk foreign key (user_id) references user (id)
);

insert into user (id,user_name,email) values (100, "juan perez", "xxxxxxx@gmail.com");
