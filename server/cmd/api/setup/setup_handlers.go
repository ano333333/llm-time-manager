package setuphandlers

import (
	"database/sql"
	"net/http"

	"github.com/ano333333/llm-time-manager/server/internal/handler"
	"github.com/ano333333/llm-time-manager/server/internal/store"
)

func SetupHandlers(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// リポジトリ
	captureScheduleStore := store.DefaultCaptureScheduleStore{DB: db}
	goalStore := store.DefaultGoalStore{DB: db}
	transactionStore := store.DefaultTransactionStore{DB: db}

	// ハンドラ
	mux.Handle("/capture/schedule", &handler.CaptureScheduleHandler{
		CaptureScheduleStore: &captureScheduleStore,
		TransactionStore:     &transactionStore,
	})
	mux.Handle("/goal", &handler.GoalHandler{
		GoalStore:        &goalStore,
		TransactionStore: &transactionStore,
	})

	return mux
}
