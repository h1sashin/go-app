-- Description: This file contains the SQL script to create the initial database schema.
-- Enum for roles
CREATE TYPE role AS ENUM ('admin', 'user');

-- Users table
CREATE TABLE
  users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role role NOT NULL DEFAULT 'user'
  );

-- Language table
CREATE TABLE
  languages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(5) NOT NULL
  );

-- Country table
CREATE TABLE
  countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(2) NOT NULL
  );

-- Product table
CREATE TABLE
  products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    vendor VARCHAR(255),
    calories INT NOT NULL,
    fat INT NOT NULL,
    carbs INT NOT NULL,
    protein INT NOT NULL,
    country_id UUID REFERENCES countries (id)
  );

-- Product translations table
CREATE TABLE
  product_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    ingredients TEXT NOT NULL DEFAULT '',
    barcode VARCHAR(255) NOT NULL DEFAULT '',
    package_size INT NOT NULL,
    product_id UUID REFERENCES products (id) ON DELETE CASCADE,
    language_id UUID REFERENCES languages (id) ON DELETE CASCADE
  );

-- Recipes table
CREATE TABLE
  recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    country_id UUID REFERENCES countries (id),
    author_id UUID REFERENCES users (id) ON DELETE SET NULL
  );

-- Enum for units
CREATE TYPE unit AS ENUM ('g', 'ml', 'tbsp', 'tsp', 'cup', 'oz', 'lb');

CREATE TABLE
  ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    unit unit NOT NULL,
    recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE,
    product_id UUID REFERENCES products (id)
  );