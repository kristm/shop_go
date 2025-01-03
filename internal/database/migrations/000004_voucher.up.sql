CREATE TABLE vouchers (
  id INTEGER NOT NULL PRIMARY KEY,
  voucher_type_id INTEGER NOT NULL,
  code VARCHAR(100) NOT NULL UNIQUE,
  valid BOOL NOT NULL DEFAULT TRUE,
  expires_at TIMESTAMP(3) NOT NULL,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (voucher_type_id) REFERENCES voucher_types(id)
);

CREATE TABLE voucher_types (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description VARCHAR(100) NOT NULL,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO voucher_types VALUES 
  (NULL, "less30", "30% off", CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (NULL, "less50", "50% off", CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (NULL, "freeship", "Free Shipping", CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

ALTER TABLE orders ADD voucher VARCHAR(50);

