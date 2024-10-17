CREATE OR REPLACE FUNCTION generate_member_code(OUT member_code TEXT)
AS $$
DECLARE
    new_code BIGINT;
BEGIN
    -- Mengambil nilai urutan terbaru dari sequence atau table
    SELECT COALESCE(MAX(SUBSTRING(c.member_code, 2)::BIGINT), 0) + 1 INTO new_code
    FROM customers c;

    -- Menggabungkan nomor urut dengan prefix '0' menjadi 12 digit
    member_code := LPAD(new_code::TEXT, 12, '0');
END;
$$ LANGUAGE plpgsql;
