CREATE TABLE postRec( 
    id   integer primary key not null,
    title text,
    text text,
    tag_ids text,
    owner_id text,
    is_private integer,
    date text,
    create_date text,
    path text
);
INSERT INTO postRec(id, title, text, tag_ids, owner_id, is_private, date, create_date, path)
   SELECT id, title, text, tag_ids, owner_id, is_private, date, create_date, path FROM post;
DROP TABLE post;
ALTER TABLE postRec RENAME TO post;