package repository

import (
	"context"
	"database/sql"
	"fmt"
	databasesql "github.com/jamiewhitney/fairways-core/pkg/database/mysql"
)

type PricingRepository interface {
	GetPriceRule(t int, d int, c int64, golfers int) (float64, error)
	GetBasePrice(c int64) (float64, error)
	LookupBand(courseId string) (int64, error)
	UpdateBasePrice(ctx context.Context, courseId int64, year int64, price float64) error
}

type MySQLPricingRepository struct {
	db *databasesql.DB
}

func NewMySQLPricingRepository(db *databasesql.DB) PricingRepository {
	return &MySQLPricingRepository{
		db: db,
	}
}

func (this *MySQLPricingRepository) GetPriceRule(week int, hour int, band int64, golfers int) (float64, error) {
	var price float64

	golfersMap := map[int]string{
		1: "golfers_1_modifier",
		2: "golfers_2_modifier",
		3: "golfers_3_modifier",
		4: "golfers_4_modifier",
	}

	query := fmt.Sprintf(`SELECT %s FROM pricing.course_pricing_rules WHERE WEEK = ? AND HOUR = ? AND price_band = ?;`, golfersMap[golfers])

	if err := this.db.Pool.QueryRow(query, week, band, hour).Scan(&price); err != nil {
		return 0, err
	}

	//if err := this.db.InTx(context.Background(), func(tx *sql.Tx) error {
	//	var err error
	//	query := fmt.Sprintf(`SELECT %s FROM pricing.course_pricing_rules WHERE WEEK = ? AND HOUR = ? AND price_band = ?;`, golfersMap[golfers])
	//	price, err = getPriceRule(context.Background(), tx, query, week, hour, band)
	//	return err
	//}); err != nil {
	//	return 0, err
	//}

	return price, nil
}

func getPriceRule(ctx context.Context, fn *sql.Tx, query string, args ...interface{}) (float64, error) {
	row := fn.QueryRowContext(ctx, query, []interface{}{args})

	var price float64
	if err := row.Scan(&price); err != nil {
		return 0, err
	}
	return price, nil
}

func (this *MySQLPricingRepository) GetBasePrice(c int64) (float64, error) {
	result := this.db.Pool.QueryRow(`SELECT price FROM pricing.course_base_prices WHERE course_id = ?`, c)
	if result.Err() != nil {
		return 0, result.Err()
	}

	var d float64
	if err := result.Scan(&d); err != nil {
		return 0, err
	}

	return d, nil
}

func (this *MySQLPricingRepository) LookupBand(courseId string) (int64, error) {
	result := this.db.Pool.QueryRow(`SELECT band FROM pricing.course_price_bands WHERE course_id = ?`, courseId)
	if result.Err() != nil {
		return 0, result.Err()
	}
	var r int64
	if err := result.Scan(&r); err != nil {
		return 0, err
	}
	return r, nil
}

func (this *MySQLPricingRepository) UpdateBasePrice(ctx context.Context, courseId int64, year int64, price float64) error {
	//_, err := this.db.ExecContext(ctx, `INSERT INTO pricing.course_base_prices (course_id, year, price) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE price = ?`, courseId, year, price, price)
	//if err != nil {
	//	return err
	//}
	//return nil
	return nil
}
