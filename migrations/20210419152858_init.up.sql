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

-- create index tablename_colname_idx
--     on tablename (colname);

do
$$
    declare
    begin
    end ;
$$;
