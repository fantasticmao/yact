CREATE TABLE clash_traffic (
  up        BIGINT,
  down      BIGINT,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clash_rule_match (
  id        VARCHAR,
  type      VARCHAR,
  duration  BIGINT,
  error     VARCHAR,
  proxy     VARCHAR,
  rule      VARCHAR,
  payload   VARCHAR,
  metadata  JSON,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clash_proxy_dial (
  id        VARCHAR,
  type      VARCHAR,
  duration  BIGINT,
  error     VARCHAR,
  proxy     VARCHAR,
  chain     JSON,
  address   VARCHAR,
  host      VARCHAR,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clash_dns_request (
  id        VARCHAR,
  type      VARCHAR,
  duration  BIGINT,
  error     VARCHAR,
  dnsType   VARCHAR,
  name      VARCHAR,
  qType     VARCHAR,
  answer    JSON,
  source    VARCHAR,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
