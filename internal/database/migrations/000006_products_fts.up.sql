CREATE VIRTUAL TABLE products_fts USING fts5 (product_id, name, description);

INSERT INTO products_fts (product_id, name, description)
  SELECT id, name, description FROM products;

CREATE TRIGGER insert_products_fts after INSERT ON products
begin
  INSERT INTO products_fts (product_id, name, description)
  VALUES (NEW.id, NEW.name, NEW.description);
end;

CREATE TRIGGER update_products_fts after UPDATE ON products
begin
  UPDATE products_fts
  SET
    name = NEW.name,
    description = NEW.description
  WHERE product_id = NEW.id;
end;

CREATE TRIGGER delete_products_fts after DELETE ON products
begin
  DELETE FROM products_fts
  WHERE product_id = OLD.id;
end;
