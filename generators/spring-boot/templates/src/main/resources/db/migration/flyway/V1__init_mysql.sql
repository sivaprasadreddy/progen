create table users
(
    id         bigint       not null auto_increment,
    email      varchar(200) not null,
    password   varchar(200) not null,
    name       varchar(200) not null,
    role       varchar(20)  not null,
    created_at timestamp    not null,
    updated_at timestamp,
    primary key (id),
    constraint user_email_unique unique (email)
);

insert into users(email, password, name, role, created_at) values
('admin@gmail.com', '$2a$10$hKDVYxLefVHV/vtuPhWD3OigtRyOykRLDdUAp80Z1crSoS1lFqaFS', 'Administrator', 'ROLE_ADMIN', CURRENT_TIMESTAMP),
('demouser@gmail.com', '$2a$10$CDAk0r/V8cvHUj.gzstta.xfvbkJY6kHrjL8pzyQWLsWMzUmByWVC', 'Demo User','ROLE_USER', CURRENT_TIMESTAMP);
