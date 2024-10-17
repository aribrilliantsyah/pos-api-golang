DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        INSERT INTO users (username, password_hash, role, full_name)
        VALUES ('admin', '$2a$14$TEVIy39sh8CC/RmDHdi7E.GeDsa3XAXa.FrT6waWR4ZVeWpLEYU.q', 'admin', 'Super Admin');
    ELSE
        RAISE NOTICE 'Username already exists';
    END IF;
END $$;