CREATE TABLE transactions (
    id bigserial not null primary key,
    amount varchar not null,
    spent boolean not null,
    created_at timestamp
);

CREATE TABLE wallet (
    id bigserial not null primary key,
    btc_balance varchar not null,
    eur_balance varchar not null
)