CREATE DATABASE IF NOT EXISTS go_market;
USE go_market;

CREATE TABLE IF NOT EXISTS products
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255),
    description TEXT,
    price       DECIMAL(10, 2),
    quantity    INT
);

CREATE TABLE IF NOT EXISTS clients
(
    id      INT AUTO_INCREMENT PRIMARY KEY,
    name    VARCHAR(255),
    surname VARCHAR(255),
    phone   VARCHAR(20),
    address TEXT,
    email   VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS orders
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    client_id     INT,
    product_id    INT,
    quantity      INT,
    price         DECIMAL(10, 2),
    purchase_date DATE,
    FOREIGN KEY (client_id) REFERENCES clients (id),
    FOREIGN KEY (product_id) REFERENCES products (id)
);
