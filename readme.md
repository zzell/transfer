### Run app
```bash
docker-compose up -d
```

### Requests example
- POST transfer
```bash
curl --location --request POST 'localhost:5580/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": 1,
    "to": 2,
    "amount": 100
}'
```

- GET wallet
```bash
curl --location --request GET 'localhost:5580/wallet/5'
```

### Technical gaps:
- Code is not covered by unit-tests intentionally, `web/transfer/transfer_test.go` has a simple example how it could be done in real life.
- Commission subtracted from each transaction is not stored anywhere for simplicity.
- `SERIAL` is used as ID for simplicity, UUID is more preferable
- Money stored as floats, but it's a poor decision due to problems with rounding and math operations
- No logging
- No ORM
- No API for wallets management