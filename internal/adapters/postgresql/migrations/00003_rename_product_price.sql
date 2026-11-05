-- +goose Up
ALTER TABLE products
  RENAME COLUMN price_in_centers TO price_in_cents;

-- +goose Down
ALTER TABLE products
  RENAME COLUMN price_in_cents TO price_in_centers;