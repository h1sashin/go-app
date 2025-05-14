-- name: CreateProduct :one
INSERT INTO products(
    vendor,
    kcal,
    barcode,
    fat,
    carbs,
    protein,
    country_id
  )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProductByID :one
WITH translations AS (
  SELECT pt.name,
    pt.ingredients,
    pt.barcode,
    pt.package_size,
    pt.unit,
    l.code AS language,
    pt.product_id
  FROM product_translations pt
    JOIN languages l ON l.id = pt.language_id
  WHERE pt.product_id = $1
    AND l.code = $2
)
SELECT tr.name,
  tr.ingredients,
  tr.barcode,
  tr.package_size,
  tr.unit,
  p.vendor,
  p.kcal,
  p.fat,
  p.saturated_fat,
  p.carbs,
  p.sugar,
  p.protein,
  p.salt
FROM products p
  LEFT JOIN translations tr ON tr.product_id = p.id
WHERE p.id = $1
  AND (
    $3::UUID IS NULL
    OR p.country_id = $3
  );

-- name: GetProducts :many
WITH filtered_products AS (
  SELECT p.id,
    p.created_at,
    p.vendor,
    p.kcal,
    p.fat,
    p.saturated_fat,
    p.carbs,
    p.sugar,
    p.protein,
    p.salt,
    p.country_id
  FROM products p
  WHERE (
      $1::TEXT IS NULL
      OR $1 = ''
      OR EXISTS (
        SELECT 1
        FROM product_translations pt
          JOIN languages l ON l.id = pt.language_id
        WHERE pt.product_id = p.id
          AND l.code = $4
          AND pt.name ILIKE '%' || $1 || '%'
      )
    )
    AND (
      $2::UUID IS NULL
      OR p.country_id = $2
    )
    AND (
      (
        $3::TEXT IS NULL
        OR $3 = ''
      )
      OR p.id > $3::UUID
    )
  ORDER BY p.id
  LIMIT $5
), translations AS (
  SELECT pt.name,
    pt.ingredients,
    pt.barcode,
    pt.package_size,
    pt.unit,
    l.code AS language,
    pt.product_id
  FROM product_translations pt
    JOIN languages l ON l.id = pt.language_id
  WHERE l.code = $4
    AND pt.product_id IN (
      SELECT id
      FROM filtered_products
    )
),
total_count AS (
  SELECT COUNT(*) AS total
  FROM products p
  WHERE (
      $1::TEXT IS NULL
      OR $1 = ''
      OR EXISTS (
        SELECT 1
        FROM product_translations pt
          JOIN languages l ON l.id = pt.language_id
        WHERE pt.product_id = p.id
          AND l.code = $4
          AND pt.name ILIKE '%' || $1 || '%'
      )
    )
    AND (
      $2::UUID IS NULL
      OR p.country_id = $2
    )
)
SELECT fp.id,
  fp.created_at,
  tr.name,
  tr.ingredients,
  tr.barcode,
  tr.package_size,
  tr.unit,
  fp.vendor,
  fp.kcal,
  fp.fat,
  fp.saturated_fat,
  fp.carbs,
  fp.sugar,
  fp.protein,
  fp.salt,
  (
    SELECT total
    FROM total_count
  ) AS total_count,
  (
    SELECT EXISTS(
        SELECT 1
        FROM products p2
        WHERE (
            $1::TEXT IS NULL
            OR $1 = ''
            OR EXISTS (
              SELECT 1
              FROM product_translations pt
                JOIN languages l ON l.id = pt.language_id
              WHERE pt.product_id = p2.id
                AND l.code = $4
                AND pt.name ILIKE '%' || $1 || '%'
            )
          )
          AND (
            $2::UUID IS NULL
            OR p2.country_id = $2
          )
          AND p2.id > (
            SELECT MAX(id)
            FROM filtered_products
          )
        LIMIT 1
      )
  ) AS has_next_page,
  (
    SELECT $3 IS NOT NULL
      AND $3 <> ''
  ) AS has_previous_page,
  (
    SELECT MAX(id)
    FROM filtered_products
  ) AS next_cursor
FROM filtered_products fp
  LEFT JOIN translations tr ON tr.product_id = fp.id
ORDER BY fp.id;