CREATE TABLE categories (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(255) NOT NULL,
  category_id INTEGER,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);

INSERT INTO categories VALUES (NULL, "Prints", CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "Stickers", CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO products VALUES (NULL, "Time Spent with Cats", "Artwork", 1, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "Breton", "Reprint of Gouache artwork", 1, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO products VALUES (NULL, "Gamer Cat", "", 2, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP), (NULL, "WFH Cat", "", 2, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);