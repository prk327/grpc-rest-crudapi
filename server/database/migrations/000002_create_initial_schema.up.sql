-- Create sequence for auto-incrementing ID
CREATE SEQUENCE IF NOT EXISTS items_id_seq;

-- Create main items table
CREATE TABLE IF NOT EXISTS items (
    id INT NOT NULL DEFAULT NEXTVAL('items_id_seq') PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT SYSDATE,
    updated_at TIMESTAMP DEFAULT SYSDATE
);

-- Create projection for better performance
CREATE PROJECTION IF NOT EXISTS items_projection 
/*+basename(items),createtype(L)*/
(
    id ENCODING AUTO,
    name ENCODING AUTO,
    description ENCODING AUTO,
    created_at ENCODING AUTO,
    updated_at ENCODING AUTO
)
AS
SELECT * FROM items
ORDER BY id
SEGMENTED BY HASH(id) ALL NODES KSAFE;

-- Create refresh command
SELECT MARK_DESIGN_KSAFE(1);