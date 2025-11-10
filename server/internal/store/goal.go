package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/google/uuid"
)

type GoalStore interface {
	GetGoal(tx Transaction, status []string) ([]datamodel.Goal, error)
	// goalsテーブルにinsertする。新規作成されたGoalを返す。
	//
	// idが指定されていない場合はUUIDを生成してinsertする。
	// kpi_*のnull/非nullが揃っているかチェックしない。
	CreateGoal(tx Transaction, id *string, title string, description string, startDate time.Time, endDate time.Time, kpiName *string, kpiTarget *float64, kpiUnit *string, status string) (datamodel.Goal, error)
}

type DefaultGoalStore struct {
	DB *sql.DB
}

func (s *DefaultGoalStore) GetGoal(tx Transaction, status []string) ([]datamodel.Goal, error) {
	defaultTx, ok := tx.(DefaultTransaction)
	if !ok {
		return nil, errors.New("transaction is not DefaultTransaction")
	}
	statusQuoted := make([]string, 0)
	for _, s := range status {
		statusQuoted = append(statusQuoted, fmt.Sprintf("'%s'", s))
	}
	statusJoined := strings.Join(statusQuoted, ",")
	query := fmt.Sprintf("SELECT id, title, description, start_date, end_date, kpi_name, kpi_target, kpi_unit, status, created_at, updated_at FROM goals WHERE status IN (%s) ORDER BY id ASC;", statusJoined)
	rows, err := defaultTx.Tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	goals := []datamodel.Goal{}
	for rows.Next() {
		var goal datamodel.Goal
		err = rows.Scan(&goal.ID, &goal.Title, &goal.Description, &goal.StartDate, &goal.EndDate, &goal.KpiName, &goal.KpiTarget, &goal.KpiUnit, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return goals, nil
}

func (s *DefaultGoalStore) CreateGoal(tx Transaction, id *string, title string, description string, startDate time.Time, endDate time.Time, kpiName *string, kpiTarget *float64, kpiUnit *string, status string) (datamodel.Goal, error) {
	emptyModel := datamodel.Goal{}

	defaultTx, ok := tx.(DefaultTransaction)
	if !ok {
		return emptyModel, errors.New("transaction is not DefaultTransaction")
	}

	args := []any{}
	if id != nil {
		args = append(args, *id)
	} else {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return emptyModel, err
		}
		args = append(args, uuid.String())
	}
	args = append(args, title, description, startDate, endDate)
	args = append(args, valueOrNil(kpiName), valueOrNil(kpiTarget), valueOrNil(kpiUnit))
	args = append(args, status)
	result := defaultTx.Tx.QueryRow(
		`INSERT INTO goals 
		(id, title, description, start_date, end_date, kpi_name, kpi_target, kpi_unit, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, title, description, start_date, end_date, kpi_name, kpi_target, kpi_unit, status, created_at, updated_at;`,
		args...,
	)
	var goal datamodel.Goal
	if err := result.Scan(&goal.ID, &goal.Title, &goal.Description, &goal.StartDate, &goal.EndDate, &goal.KpiName, &goal.KpiTarget, &goal.KpiUnit, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt); err != nil {
		return emptyModel, err
	}

	return goal, nil
}

// 非nilの場合は*valueを、nilの場合はnilを返す
func valueOrNil[T any](value *T) any {
	if value == nil {
		return nil
	}
	return *value
}
