do
$$
    begin
        execute 'ALTER DATABASE ' || current_database() || ' SET timezone = ''+06''';
    end;
$$;

create table cfg
(
    v jsonb not null default '{}'
);

create table usr
(
    id         bigserial   not null,
    created_at timestamptz not null default now(),
    type_id    smallint    not null default 0,
    phone      text        not null,
    name       text        not null default '',
    token      text        not null default '',

    constraint usr_pk primary key (id),
    constraint usr_unique_phone unique (phone)
);
create index usr_created_at_idx
    on usr (created_at);
create index usr_type_id_idx
    on usr (type_id);
create index usr_phone_idx
    on usr (phone);
create index usr_token_idx
    on usr (token);

do
$$
    declare
    begin
        -- cfg row
        insert into cfg (v) values ('{}');

        -- Admin user
        insert into usr(type_id, phone, name)
        values (1, '70000000000', 'Admin');
    end ;
$$;
