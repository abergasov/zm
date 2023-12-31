CREATE TABLE merkle_trees
(
    mt_id varchar PRIMARY KEY,
    tree  jsonb
);

CREATE TABLE files
(
    f_id       uuid PRIMARY KEY,
    tree_id    varchar CONSTRAINT files_merkle_trees_mt_id_fk REFERENCES merkle_trees,
    file_index INTEGER NOT NULL,
    file_name  VARCHAR NOT NULL,
    file_hash  VARCHAR NOT NULL
);

ALTER TABLE files ADD CONSTRAINT files_pk UNIQUE (tree_id, file_index);