package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ungame/command-time-track/app/ioext"
	"github.com/ungame/command-time-track/app/models"
)

type ActivitiesRepository interface {
	Create(ctx context.Context, activity *models.Activity) (int64, error)
	Update(ctx context.Context, activity *models.Activity) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Get(ctx context.Context, id int64) (*models.Activity, error)
	GetAll(ctx context.Context) ([]*models.Activity, error)
	Search(ctx context.Context, term string) ([]*models.Activity, error)
	GetByStatus(ctx context.Context, status models.Status) ([]*models.Activity, error)
}

type activitiesRepository struct {
	conn       *sql.DB
	createStmt *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
}

func NewActivitiesRepository(ctx context.Context, conn *sql.DB) ActivitiesRepository {
	return &activitiesRepository{
		conn:       conn,
		createStmt: mustCreateStmt(ctx, conn, insertActivityQuery),
		updateStmt: mustCreateStmt(ctx, conn, updateActivityQuery),
		deleteStmt: mustCreateStmt(ctx, conn, deleteActivityQuery),
	}
}

func (r *activitiesRepository) Close() {
	ioext.Close(r.createStmt)
	ioext.Close(r.updateStmt)
	ioext.Close(r.deleteStmt)
}

func (r *activitiesRepository) Create(ctx context.Context, activity *models.Activity) (int64, error) {
	result, err := r.createStmt.ExecContext(
		ctx,
		activity.Category,
		activity.Description,
		activity.StartedAt,
		activity.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *activitiesRepository) Update(ctx context.Context, activity *models.Activity) (int64, error) {
	result, err := r.updateStmt.ExecContext(
		ctx,
		activity.Category,
		activity.Description,
		activity.Status,
		activity.UpdatedAt,
		activity.FinishedAt,
		activity.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *activitiesRepository) Delete(ctx context.Context, id int64) (int64, error) {
	result, err := r.deleteStmt.ExecContext(ctx, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *activitiesRepository) Get(ctx context.Context, id int64) (*models.Activity, error) {
	var (
		activity = new(models.Activity)
		query    = `select * from activities where id = ?`
		row      = r.conn.QueryRowContext(ctx, query, id)
	)
	err := row.Scan(
		&activity.ID,
		&activity.Category,
		&activity.Description,
		&activity.Status,
		&activity.StartedAt,
		&activity.UpdatedAt,
		&activity.FinishedAt,
	)
	return activity, err
}

func (r *activitiesRepository) GetAll(ctx context.Context) ([]*models.Activity, error) {
	query := `select * from activities`
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer ioext.Close(rows)
	activities := make([]*models.Activity, 0, 10)
	for rows.Next() {
		var activity models.Activity
		err = rows.Scan(
			&activity.ID,
			&activity.Category,
			&activity.Description,
			&activity.Status,
			&activity.StartedAt,
			&activity.UpdatedAt,
			&activity.FinishedAt,
		)
		if err != nil {
			return activities, err
		}
		activities = append(activities, &activity)
	}
	return activities, nil
}

func (r *activitiesRepository) GetByStatus(ctx context.Context, status models.Status) ([]*models.Activity, error) {
	query := `select * from activities where status = ?`
	rows, err := r.conn.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer ioext.Close(rows)
	activities := make([]*models.Activity, 0, 10)
	for rows.Next() {
		var activity models.Activity
		err = rows.Scan(
			&activity.ID,
			&activity.Category,
			&activity.Description,
			&activity.Status,
			&activity.StartedAt,
			&activity.UpdatedAt,
			&activity.FinishedAt,
		)
		if err != nil {
			return activities, err
		}
		activities = append(activities, &activity)
	}
	return activities, nil
}

func (r *activitiesRepository) Search(ctx context.Context, term string) ([]*models.Activity, error) {
	query := `select * from activities where category like ? or description like ?`
	rows, err := r.conn.QueryContext(ctx, query, like(term), like(term))
	if err != nil {
		return nil, err
	}
	defer ioext.Close(rows)
	activities := make([]*models.Activity, 0, 10)
	for rows.Next() {
		var activity models.Activity
		err = rows.Scan(
			&activity.ID,
			&activity.Category,
			&activity.Description,
			&activity.Status,
			&activity.StartedAt,
			&activity.UpdatedAt,
			&activity.FinishedAt,
		)
		if err != nil {
			return activities, err
		}
		activities = append(activities, &activity)
	}
	return activities, nil
}

func like(s string) string {
	return fmt.Sprintf(`%%%s%%`, s)
}
