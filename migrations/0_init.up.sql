CREATE TABLE prices
(
    ts        bigserial PRIMARY KEY,
    ask_price REAL NOT NULL,
    bid_price REAL NOT NULL
);
