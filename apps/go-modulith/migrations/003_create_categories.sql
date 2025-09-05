-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert some sample categories
INSERT INTO categories (name) VALUES 
('Electronics'),
('Books'),
('Clothing'),
('Home & Garden'),
('Sports & Outdoors');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd