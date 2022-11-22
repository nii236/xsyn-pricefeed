# XSYN Pricefeed

A REST API that returns the price of Ether and SUPS.

To regenerate the contract bindings, install `solc` and `abigen`:

```
cd ethusd
solc --abi ethusd.sol -o .
abigen --abi=AggregatorV3Interface.abi --pkg=ethusd --out=ethusd.go

cd ..
cd supseth
solc --abi supseth.sol -o .
abigen --abi=IUniswapV3PoolState.abi.abi --pkg=supseth --out=supseth.go
```

To run:

```
go run main.go --rpc_url {{RPC_URL}}
```

## Migration

```sql
CREATE TABLE kv (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT UNIQUE NOT NULL DEFAULT '',
    value TEXT NOT NULL DEFAULT '',

    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO kv (key, value) VALUES ('last_block', '0') ON CONFLICT DO NOTHING;
INSERT INTO kv (key, value) VALUES ('block_height', '0') ON CONFLICT DO NOTHING;

CREATE TABLE prices (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sups_price_cents TEXT NOT NULL,
    eth_price_cents TEXT NOT NULL,
    bnb_price_cents TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transfers (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    log_index INTEGER NOT NULL,
    block INTEGER NOT NULL,
    chain_id INTEGER NOT NULL,
    contract TEXT NOT NULL,
    symbol TEXT NOT NULL,
    decimals INTEGER NOT NULL,
    tx_id TEXT NOT NULL,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    amount NUMERIC(28),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tx_id, log_index, block)
);
```

## Random commands

`docker run -d -p 5432:5432 --name xsyn-pricefeed-db -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=dev -e POSTGRES_DB=xsyn-pricefeed postgres:13-alpine`
