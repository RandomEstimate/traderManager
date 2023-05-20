CREATE TABLE real_strategy_table
(
    strategy VARCHAR(100) PRIMARY KEY
);

CREATE TABLE test_strategy_table
(
    strategy VARCHAR(100) PRIMARY KEY
);

CREATE TABLE real_order_table
(
    strategy VARCHAR(100),
    symbol   VARCHAR(100),
    price    float,
    side     int,
    qty      FLOAT,
    time     bigint
);

CREATE TABLE test_order_table
(
    strategy VARCHAR(100),
    symbol   VARCHAR(100),
    price    float,
    side     int,
    qty      FLOAT,
    time     bigint
);

CREATE TABLE real_value_table
(
    strategy VARCHAR(100),
    symbol   VARCHAR(100),
    profit   float
);

CREATE TABLE test_value_table
(
    strategy VARCHAR(100),
    symbol   VARCHAR(100),
    profit   float
);