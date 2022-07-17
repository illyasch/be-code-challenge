// Package calc implements calculations of Ethereum gas prices per hour.
package calc

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Calc contains the database for storing transactions.
type Calc struct {
	DB *sqlx.DB
}

type Hour struct {
	Timestamp float64 `json:"t"`
	Amount    float64 `json:"v"`
}

// New constructs a new Calc.
func New(db *sqlx.DB) Calc {
	return Calc{DB: db}
}

// Hourly calculates transaction fees per hour.
func (e Calc) Hourly(ctx context.Context) ([]Hour, error) {
	const sql = `SELECT extract(epoch from date_trunc('hour', "block_time")) AS "timestamp", sum(gas_used * gas_price) / 1E9 AS "amount" 
    FROM transactions WHERE
        "to" <> '0x0000000000000000000000000000000000000000' AND
        "from" <> '0x0000000000000000000000000000000000000000' AND
        "to" NOT IN (SELECT address FROM contracts) AND
        "from" NOT IN (SELECT address FROM contracts)
    GROUP BY 1
    ORDER BY 1 ASC`

	rows, err := e.DB.QueryxContext(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", sql, err)
	}
	defer func() { _ = rows.Close() }()

	var hh []Hour
	for rows.Next() {
		h := Hour{}
		if err := rows.StructScan(&h); err != nil {
			return nil, fmt.Errorf("StructScan: %w", err)
		}
		hh = append(hh, h)
	}

	return hh, nil
}
