create table todos (
    id varchar(36) not null,
    title varchar(150) not null,
    tags varchar(1500) not null default '',
    completed_at TIMESTAMP DEFAULT NULL,
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    primary key todos_pk(id)
);

create view v_todos as 
    select id, title, tags
    from todos
;

CREATE TRIGGER todos_row_update BEFORE UPDATE
ON todos
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;
