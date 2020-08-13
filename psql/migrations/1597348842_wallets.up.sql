CREATE TABLE IF NOT EXISTS wallets
(
    PRIMARY KEY (id),
    id          SERIAL,
    currency_id int,
    score       float(32)
);

INSERT INTO wallets (currency_id, score) VALUES (1, 100500);
INSERT INTO wallets (currency_id, score) VALUES (1, 20000);
INSERT INTO wallets (currency_id, score) VALUES (2, 100500);
INSERT INTO wallets (currency_id, score) VALUES (2, 20000);
