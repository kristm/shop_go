CREATE TABLE customers (
  id INTEGER NOT NULL PRIMARY KEY,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  phone VARCHAR(50) NOT NULL,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
  id INTEGER NOT NULL PRIMARY KEY,
  shipping_id INTEGER,
  customer_id INTEGER,
  amount_in_cents INTEGER,
  status INTEGER NOT NULL DEFAULT 0, -- pending | canceled | paid
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
  FOREIGN KEY (shipping_id) REFERENCES shipping(id)
);

CREATE TABLE order_products (
  id INTEGER NOT NULL PRIMARY KEY,
  order_id INTEGER,
  product_id INTEGER,
  qty INTEGER DEFAULT 1,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE shipping (
  id INTEGER NOT NULL PRIMARY KEY,
  customer_id INTEGER,
  status INTEGER NOT NULL DEFAULT 0, -- pending | in transit | failed | completed
  address VARCHAR(100) NOT NULL,
  city VARCHAR(100) NOT NULL,
  country VARCHAR(100) NOT NULL,
  zip VARCHAR(10) NOT NULL,
  phone VARCHAR(50) NOT NULL,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);
