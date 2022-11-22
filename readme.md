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

CREATE TABLE whitelisted_addresses (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    address TEXT UNIQUE NOT NULL,
    chain_id INTEGER NOT NULL,
    note TEXT NOT NULL,

    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO whitelisted_addresses (address, chain_id, note) VALUES ('0x2b64E6F89dB700272cBDB740C099a460754a8DA5', '5', 'goerli abs purchase address') ON CONFLICT DO NOTHING;
INSERT INTO whitelisted_addresses (address, chain_id, note) VALUES ('0x48e6f3e175C58181086AEC640f21815C5EbF4fC0', '1', 'mainnet abs purchase address') ON CONFLICT DO NOTHING;
INSERT INTO whitelisted_addresses (address, chain_id, note) VALUES ('0x70B6a7b6768a9dCE4569316C573aa806c59c8160', '5', 'goerli presale purchase address') ON CONFLICT DO NOTHING;
INSERT INTO whitelisted_addresses (address, chain_id, note) VALUES ('0x0f3FB0E3800927Fa4757eAF9BeBD9982c5534CC3', '1', 'mainnet presale purchase address') ON CONFLICT DO NOTHING;

CREATE TABLE kv (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT UNIQUE NOT NULL DEFAULT '',
    value TEXT NOT NULL DEFAULT '',

    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO kv (key, value) VALUES ('last_block_goerli_sups', '7859764') ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value;
INSERT INTO kv (key, value) VALUES ('last_block_goerli_eth', '7859764') ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value;
INSERT INTO kv (key, value) VALUES ('last_block_mainnet_sups', '15879854') ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value;
INSERT INTO kv (key, value) VALUES ('last_block_mainnet_eth', '15879854') ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value;

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
    timestamp INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tx_id, log_index, block)
);
```

## Random commands

`docker run -d -p 5432:5432 --name xsyn-pricefeed-db -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=dev -e POSTGRES_DB=xsyn-pricefeed postgres:13-alpine`
