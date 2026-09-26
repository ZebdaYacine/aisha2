DROP TRIGGER IF EXISTS cart_items_touch_updated_at ON cart_items;
DROP TRIGGER IF EXISTS carts_touch_updated_at ON carts;
DROP FUNCTION IF EXISTS touch_cart_updated_at();
DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
