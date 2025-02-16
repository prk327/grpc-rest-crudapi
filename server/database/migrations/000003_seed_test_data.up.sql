-- Insert sample data
INSERT INTO omniq.items (name, description) VALUES
('Test Item 1', 'First test item'),
('Test Item 2', 'Second test item'),
('Test Item 3', 'Third test item');

-- Commit the transaction
COMMIT;

-- Refresh projections
SELECT REFRESH('items');