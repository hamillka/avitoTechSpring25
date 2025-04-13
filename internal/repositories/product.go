package repositories

import (
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductRepository struct {
	db *sqlx.DB
}

const (
	addProduct                = "INSERT INTO products (product_type, reception_id) VALUES ($1, $2) RETURNING id, date_time, product_type, reception_id"
	getLastProduct            = "SELECT * FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1"
	deleteProduct             = "DELETE FROM products WHERE id = $1"
	getProductsByReceptionIds = `
	SELECT id, date_time, product_type, reception_id 
	FROM products 
	WHERE reception_id = ANY($1) AND date_time BETWEEN $2 AND $3
`
)

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (pr *ProductRepository) AddProduct(productType, receptionId string) (models.Product, error) {
	var product models.Product

	err := pr.db.QueryRow(addProduct, productType, receptionId).
		Scan(
			&product.Id,
			&product.DateTime,
			&product.Type,
			&product.ReceptionId,
		)
	if err != nil {
		return models.Product{}, dto.ErrDBInsert
	}

	return product, nil
}

func (pr *ProductRepository) GetLastProduct(recId string) (models.Product, error) {
	var product models.Product
	err := pr.db.QueryRow(getLastProduct, recId).
		Scan(
			&product.Id,
			&product.DateTime,
			&product.Type,
			&product.ReceptionId,
		)
	if err != nil {
		return models.Product{}, dto.ErrNoProductsInReception
	}

	return product, nil
}

func (pr *ProductRepository) DeleteProduct(prodId string) error {
	_, err := pr.db.Exec(deleteProduct, prodId)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) GetProductsByReceptionIds(recIds []string, startDate, endDate *time.Time) ([]models.Product, error) {
	start := time.Time{}
	end := time.Now()

	if startDate != nil {
		start = *startDate
	}

	if endDate != nil {
		end = *endDate
	}

	var products []models.Product

	rows, err := pr.db.Query(getProductsByReceptionIds, pq.Array(recIds), start, end)
	if err != nil {
		return products, nil
	}
	defer rows.Close()

	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(
			&product.Id,
			&product.DateTime,
			&product.Type,
			&product.ReceptionId,
		); err != nil {
			return nil, dto.ErrDBRead
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return products, nil
	}

	return products, nil
}
