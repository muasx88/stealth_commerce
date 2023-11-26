CREATE TABLE IF NOT EXISTS admin (
    id SERIAL PRIMARY KEY,
	username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    sku VARCHAR(15) NOT NULL,
    qty INTEGER NOT NULL,
    description VARCHAR(255) NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS buyers (
	id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    secret_key VARCHAR(255) NOT NULL UNIQUE,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS carts (
	id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    buyer_id INT NOT NULL,
    qty INT NOT NULL,
    note VARCHAR(255) NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_buyer_id FOREIGN KEY (buyer_id) REFERENCES buyers (id),
 	CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE INDEX IF NOT EXISTS idx_buyer_id ON carts (buyer_id);
CREATE INDEX idx_product_id ON carts (product_id);

DROP TYPE IF EXISTS order_status_type;
CREATE TYPE order_status_type AS ENUM (
	'PENDING',
	'COMPLETED',
	'CANCELED'
);

CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	order_number VARCHAR(10) NOT NULL UNIQUE,
    buyer_id INT NOT NULL,
    grand_total INT NOT NULL DEFAULT 0,
    status "order_status_type" NOT NULL DEFAULT 'PENDING'::order_status_type,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_buyer_id FOREIGN KEY (buyer_id) REFERENCES buyers (id)
);
CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders (buyer_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders (order_number);

CREATE TABLE IF NOT EXISTS order_details (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    qty INT NOT NULL DEFAULT 1,
    price INT NOT NULL DEFAULT 0,
    total INT NOT NULL DEFAULT 0,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_order_product_id FOREIGN KEY (product_id) REFERENCES products (id)
);
CREATE INDEX IF NOT EXISTS idx_order_detail_product_id ON order_details (product_id);
CREATE INDEX IF NOT EXISTS idx_order_detail_order_id ON order_details (order_id);

-- insert admin data
INSERT INTO admin (username, password, fullname) VALUES ('admin', '$2a$10$.3A/kqoU7kspEHdxOHxpBOfP50vkrC1BbPgWUruuqI6VuBEYWXveu', 'Si Admin');


-- insert buyer data
INSERT INTO buyers (fullname, address, secret_key) VALUES ('Peter Parker', 'Queens', 'jaringlaba123');
INSERT INTO buyers (fullname, address, secret_key) VALUES ('Bruce Wayne', 'Gotham City', 'kelelawar123');


-- payment table
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    status VARCHAR(20),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payment_order_id FOREIGN KEY (order_id) REFERENCES orders (id)
);
CREATE INDEX IF NOT EXISTS idx_payment_order_id ON payments (order_id);