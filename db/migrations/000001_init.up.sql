-- Create ENUM types
CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TYPE unit AS ENUM ('g', 'ml', 'tbsp', 'tsp', 'cup', 'oz', 'lb');

CREATE TYPE language AS ENUM ('en', 'pl');

-- Create function to update timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();

RETURN NEW;

END;

$$ LANGUAGE plpgsql;

INSERT INTO products (name, brand,..costam)
VALUES ('Pepsi', 'Pepsi',..costam)
INSERT INTO product_translations (name, product_id, language_id)
VALUES ('Pepsi', 'Pepsi', 'en') CREATE TABLE product_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    ingredients TEXT NOT NULL DEFAULT '',
    product_id UUID REFERENCES products (id) ON DELETE CASCADE,
    language_id UUID REFERENCES languages (id) ON DELETE CASCADE
  );

-- Create tables
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email VARCHAR(255) UNIQUE NOT NULL,
  PASSWORD VARCHAR(255) NOT NULL,
  role ROLE NOT NULL DEFAULT 'user'
);

CREATE TABLE languages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(5) NOT NULL
);

CREATE INDEX idx_languages_code ON languages(code);

CREATE TABLE countries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(2) NOT NULL
);

CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  vendor VARCHAR(255),
  kcal INT NOT NULL,
  fat INT,
  saturated_fat INT,
  ingedients TEXT NOT NULL DEFAULT '',
  carbs INT,
  sugar INT,
  protein INT,
  salt INT,
  country_id UUID REFERENCES countries (id)
);

CREATE INDEX idx_products_country_id ON products(country_id);

CREATE TABLE product_translations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  ingredients TEXT NOT NULL DEFAULT '',
  barcode VARCHAR(13) NOT NULL DEFAULT '',
  package_size INT NOT NULL,
  unit UNIT NOT NULL,
  product_id UUID REFERENCES products (id) ON DELETE CASCADE,
  language_id UUID REFERENCES languages (id) ON DELETE CASCADE
);

CREATE INDEX idx_product_translations_product_id ON product_translations(product_id);

CREATE INDEX idx_product_translations_language_id ON product_translations(language_id);

CREATE INDEX idx_product_translations_product_language ON product_translations(product_id, language_id);

CREATE TABLE recipes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  country_id UUID REFERENCES countries (id),
  author_id UUID REFERENCES users (id) ON DELETE
  SET NULL
);

CREATE TABLE ingredients (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  amount INT NOT NULL,
  unit UNIT NOT NULL,
  recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE,
  product_id UUID REFERENCES products (id)
);

-- Insert default values
INSERT INTO languages (name, code)
VALUES ('English', 'en'),
  ('Polski', 'pl');

INSERT INTO countries (name, code)
VALUES ('Polska', 'PL');

-- Create triggers
CREATE OR REPLACE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_languages_updated_at BEFORE
UPDATE ON languages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_countries_updated_at BEFORE
UPDATE ON countries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_products_updated_at BEFORE
UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_product_translations_updated_at BEFORE
UPDATE ON product_translations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_recipes_updated_at BEFORE
UPDATE ON recipes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_ingredients_updated_at BEFORE
UPDATE ON ingredients FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();