DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        INSERT INTO users (username, password_hash, role, full_name)
        VALUES ('admin', '$2a$14$TEVIy39sh8CC/RmDHdi7E.GeDsa3XAXa.FrT6waWR4ZVeWpLEYU.q', 'admin', 'Super Admin');
    ELSE
        RAISE NOTICE 'Username already exists';
    END IF;
END $$;

-- Seeder for categories (only if the specific category is not already present)
INSERT INTO categories (name, created_by, updated_by, created_at)
SELECT 'Elektronik', 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Elektronik');

INSERT INTO categories (name, created_by, updated_by, created_at)
SELECT 'Grosir', 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Grosir');

INSERT INTO categories (name, created_by, updated_by, created_at)
SELECT 'Pakaian', 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Pakaian');

INSERT INTO categories (name, created_by, updated_by, created_at)
SELECT 'Buku', 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Buku');

-- Seeder for products (only if the specific product is not already present)
INSERT INTO products (name, price, stock, category_id, created_by, updated_by, created_at)
SELECT 'Smartphone', 3000000, 20, 1, 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Smartphone');

INSERT INTO products (name, price, stock, category_id, created_by, updated_by, created_at)
SELECT 'Laptop', 5000000, 10, 1, 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Laptop');

INSERT INTO products (name, price, stock, category_id, created_by, updated_by, created_at)
SELECT 'Katong Kresek', 100, 1000, 2, 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Katong Kresek');

INSERT INTO products (name, price, stock, category_id, created_by, updated_by, created_at)
SELECT 'Kaos Anime', 85000, 200, 3, 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Kaos Anime');

INSERT INTO products (name, price, stock, category_id, created_by, updated_by, created_at)
SELECT 'Novel Fantasy', 45000, 150, 4, 1, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Novel Cinta');

-- Seeder for product_history (only if the specific transaction is not already present)
INSERT INTO product_history (trx_ref, product_id, quantity_change, type, reason, created_by, created_at)
SELECT 'TRX001', 1, 20, 'in', 'New Stock', 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM product_history WHERE trx_ref = 'TRX001');

INSERT INTO product_history (trx_ref, product_id, quantity_change, type, reason, created_by, created_at)
SELECT 'TRX002', 2, 10, 'in', 'New Stock', 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM product_history WHERE trx_ref = 'TRX002');

INSERT INTO product_history (trx_ref, product_id, quantity_change, type, reason, created_by, created_at)
SELECT 'TRX003', 3, 1000, 'in', 'New Stock', 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM product_history WHERE trx_ref = 'TRX003');

INSERT INTO product_history (trx_ref, product_id, quantity_change, type, reason, created_by, created_at)
SELECT 'TRX004', 4, 200, 'in', 'New Stock', 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM product_history WHERE trx_ref = 'TRX004');

INSERT INTO product_history (trx_ref, product_id, quantity_change, type, reason, created_by, created_at)
SELECT 'TRX005', 5, 150, 'in', 'New Stock', 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM product_history WHERE trx_ref = 'TRX005');
