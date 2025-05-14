-- Drop triggers
DROP TRIGGER IF EXISTS update_ingredients_updated_at ON ingredients;

DROP TRIGGER IF EXISTS update_recipes_updated_at ON recipes;

DROP TRIGGER IF EXISTS update_product_translations_updated_at ON product_translations;

DROP TRIGGER IF EXISTS update_products_updated_at ON products;

DROP TRIGGER IF EXISTS update_countries_updated_at ON countries;

DROP TRIGGER IF EXISTS update_languages_updated_at ON languages;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop tables
DROP TABLE IF EXISTS ingredients;

DROP TABLE IF EXISTS recipes;

DROP TABLE IF EXISTS product_translations;

DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS countries;

DROP TABLE IF EXISTS languages;

DROP TABLE IF EXISTS users;

-- Drop types
DROP TYPE IF EXISTS unit;

DROP TYPE IF EXISTS role;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column;