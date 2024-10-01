CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100)   NOT NULL,
    username VARCHAR(100)   NOT NULL,
    email    VARCHAR(100)   NOT NULL,
    balance  DECIMAL(10, 2) NOT NULL DEFAULT 0
);

-- Create roles table
CREATE TABLE roles
(
    id        SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
);

-- Insert default roles
INSERT INTO roles (role_name)
VALUES ('admin'),
       ('user');

-- Create user_roles table (junction table)
CREATE TABLE user_roles
(
    auth_id INT NOT NULL,
    role_id INT NOT NULL,
    PRIMARY KEY (auth_id, role_id),
    FOREIGN KEY (auth_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

-- Set all existing users to have the 'user' role
INSERT INTO user_roles (auth_id, role_id)
SELECT id, 2
FROM auth;



CREATE TABLE products
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100)   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    stock       INT            NOT NULL,
    image_path  VARCHAR(255)
);

-- Create an index on the name column to improve search performance
CREATE INDEX idx_products_name ON products (name);

VALUES ('Product 1', 100, 10);
INSERT INTO products (name, price, stock)
VALUES ('Product 2', 200, 20);
INSERT INTO products (name, price, stock)
VALUES ('Product 3', 300, 30);
INSERT INTO products (name, price, stock)
VALUES ('Product 4', 400, 40);
INSERT INTO products (name, price, stock)
VALUES ('Product 5', 500, 50);

CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    user_id     INT            NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_lines
(
    id         SERIAL PRIMARY KEY,
    order_id   INT            NOT NULL,
    product_id INT            NOT NULL,
    qty        INT            NOT NULL DEFAULT 1,
    total      DECIMAL(10, 2) NOT NULl DEFAULT 0
);

CREATE TABLE auth
( -- table for authentication
    id       SERIAL PRIMARY KEY,
    user_id  INT          NOT NULL,
    username VARCHAR(100) NOT NULL,
    email    VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL
);



