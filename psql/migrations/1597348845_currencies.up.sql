CREATE TABLE IF NOT EXISTS currencies
(
    id     SERIAL,
    name   text,
    symbol text
);

INSERT INTO currencies (name, symbol) VALUES ('bitcoin', 'btc');
INSERT INTO currencies (name, symbol) VALUES ('ethereum', 'eth');
