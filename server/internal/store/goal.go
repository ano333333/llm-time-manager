package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
)

type GoalStore interface {
	GetGoal(tx Transaction, status []string) ([]datamodel.Goal, error)
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
