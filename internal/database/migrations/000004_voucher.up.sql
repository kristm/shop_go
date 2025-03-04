CREATE TABLE vouchers (
  id INTEGER NOT NULL PRIMARY KEY,
  voucher_type_id INTEGER NOT NULL,
  code VARCHAR(100) NOT NULL UNIQUE,
  valid BOOL NOT NULL DEFAULT TRUE,
  minimum_spend INTEGER DEFAULT 0,
  expires_at TIMESTAMP(3) NOT NULL,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (voucher_type_id) REFERENCES voucher_types(id)
);

CREATE TABLE voucher_types (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description VARCHAR(100) NOT NULL,
  amount INTEGER DEFAULT 0,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO voucher_types VALUES 
  (NULL, "less30", "30% off", 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (NULL, "less50", "50% off", 50, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (NULL, "freeship", "Free Shipping", 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

ALTER TABLE orders ADD voucher VARCHAR(50);

insert into vouchers (voucher_type_id, code, valid, minimum_spend, expires_at) values (2, 'snake', true, 1000, date('now', '+1 month'));
