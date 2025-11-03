package datamodel

import "time"

type Goal struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	KpiName     *string   `json:"kpi_name"`   // kpi未設定の場合nil
	KpiTarget   *float64  `json:"kpi_target"` // kpi未設定の場合nil
	KpiUnit     *string   `json:"kpi_unit"`   // kpi未設定の場合nil
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
