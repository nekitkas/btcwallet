CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    amount NUMERIC(20, 6) NOT NULL,
    spent INTEGER NOT NULL CHECK (spent IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO transactions (id, amount, spent, created_at) VALUES ('ec8c04c2-8e03-4310-b4c0-363455b87c4f', 2, 0, CURRENT_TIMESTAMP);
INSERT INTO transactions (id, amount, spent, created_at) VALUES ('d7c07d4e-ca5f-4b68-9ec2-efad4c7d7c54', 3, 0, CURRENT_TIMESTAMP);
INSERT INTO transactions (id, amount, spent, created_at) VALUES ('d444477f-9bd2-4cf2-a6ad-5af35b8f12ba', 5, 0, CURRENT_TIMESTAMP);
INSERT INTO transactions (id, amount, spent, created_at) VALUES ('335e4ab2-6d35-42ef-9b59-50713f05870a', 1.5, 0, CURRENT_TIMESTAMP);