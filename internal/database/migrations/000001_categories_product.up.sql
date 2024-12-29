CREATE TABLE categories (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  enabled BOOL NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  sku VARCHAR(100) NOT NULL UNIQUE,
  description VARCHAR(255) NOT NULL,
  category_id INTEGER,
  price_in_cents INTEGER,
  status INTEGER NOT NULL DEFAULT 0, -- in stock | low stock | out of stock
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);

INSERT INTO categories VALUES (NULL, "Prints", TRUE, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "Comics", TRUE, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "Accessories", TRUE, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "Stickers", TRUE, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "Decors", TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

