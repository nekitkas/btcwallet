CREATE TABLE transactions (
    id varchar not null primary key,
    amount NUMERIC(20, 6) not null,
    spent boolean not null,
    created_at timestamp
);

INSERT INTO transactions (id, amount, spent, created_at) VALUES (gen_random_uuid(), 2, false, NOW());
INSERT INTO transactions (id, amount, spent, created_at) VALUES (gen_random_uuid(), 3, false, NOW());
INSERT INTO transactions (id, amount, spent, created_at) VALUES (gen_random_uuid(), 5, false, NOW());
INSERT INTO transactions (id, amount, spent, created_at) VALUES (gen_random_uuid(), 1.5, false, NOW());
