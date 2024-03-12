CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       passhash BYTEA NOT NULL
);

CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL
);

CREATE TABLE categoryies (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE product_category (
                                  id SERIAL PRIMARY KEY,
                                  product_id INTEGER NOT NULL,
                                  category_id INTEGER NOT NULL,
                                  FOREIGN KEY (product_id) REFERENCES products(id),
                                  FOREIGN KEY (category_id) REFERENCES categoryies(id),
                                  UNIQUE (product_id, category_id)
);