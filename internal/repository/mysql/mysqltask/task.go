package mysqltask

import (
	"context"
	"database/sql"
	"todoapp/internal/entity"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"
)

const (
	OpCreateTask           = richerror.Op("mysqltask.CreateTask")
	OpGetAllTasks          = richerror.Op("mysqltask.GetAllTasks")
	OpGetTaskByID          = richerror.Op("mysqltask.GetTaskById")
	OpGetTasksByCategory = richerror.Op("mysqltask.GetTasksByCategory")
	OpUpdateTask           = richerror.Op("mysqltask.UpdateTask")
	OpDeleteTask           = richerror.Op("mysqltask.DeleteTask")
)



func (d *DB) CreateTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	const op = OpCreateTask

	query := `INSERT INTO tasks (user_id, category_id, title, description, due_date, priority, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := d.conn.Conn().ExecContext(ctx, query,
		task.UserID,
		task.CategoryID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Priority.String(),
		task.Status.String(),
	)

	if err != nil {
		return entity.Task{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to insert task").
			WithKind(richerror.KindUnexpected)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return entity.Task{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get last insert ID").
			WithKind(richerror.KindUnexpected)
	}

	

	return d.GetTaskByID(ctx, uint(id), task.UserID)
}

func (d *DB) GetTaskByID(ctx context.Context, id uint, userID uint) (entity.Task, error) {
	const op = OpGetTaskByID

	query := `SELECT id, user_id, category_id, title, description, due_date, priority, status, created_at, updated_at
	          FROM tasks WHERE id = ? AND user_id = ?`

	var t entity.Task
	var description sql.NullString
	var dueDate sql.NullTime
	var priorityStr string
	var statusStr string

	err := d.conn.Conn().QueryRowContext(ctx, query, id, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.CategoryID,
		&t.Title,
		&description,
		&dueDate,
		&priorityStr,
		&statusStr,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Task{}, richerror.New(op).
				WithMessage(errmsg.ErrorMsgNotFound).
				WithKind(richerror.KindNotFound)
		}
		return entity.Task{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get task").
			WithKind(richerror.KindUnexpected)
	}

	if description.Valid {
		t.Description = description.String
	}

	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}

	t.Priority = entity.MapToPriorityEntity(priorityStr)
	t.Status = entity.MapToStatusEntity(statusStr)

	return t, nil
}

func (d *DB) GetAllTasks(ctx context.Context, userID uint) ([]entity.Task, error) {
	const op = OpGetAllTasks

	query := `SELECT id, user_id, category_id, title, description, due_date, priority, status, created_at, updated_at 
	          FROM tasks WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := d.conn.Conn().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get tasks").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var t entity.Task
		var description sql.NullString
		var dueDate sql.NullTime
		var priorityStr string
		var statusStr string

		err := rows.Scan(
			&t.ID, &t.UserID, &t.CategoryID, &t.Title, &description, &dueDate,
			&priorityStr, &statusStr, &t.CreatedAt, &t.UpdatedAt,
		)

		if err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
		}

		if description.Valid {
			t.Description = description.String
		}

		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}

		t.Priority = entity.MapToPriorityEntity(priorityStr)
		t.Status = entity.MapToStatusEntity(statusStr)

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (d *DB) GetTasksByCategory(ctx context.Context, userID uint, categoryID uint) ([]entity.Task, error) {
	const op = OpGetTasksByCategory

	query := `SELECT id, user_id, category_id, title, description, due_date, priority, status, created_at, updated_at
	          FROM tasks WHERE user_id = ? AND category_id = ? ORDER BY created_at DESC`

	rows, err := d.conn.Conn().QueryContext(ctx, query, userID, categoryID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get tasks by category").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var t entity.Task
		var description sql.NullString
		var dueDate sql.NullTime
		var priorityStr string
		var statusStr string

		err := rows.Scan(
			&t.ID, &t.UserID, &t.CategoryID, &t.Title, &description, &dueDate,
			&priorityStr, &statusStr, &t.CreatedAt, &t.UpdatedAt,
		)

		if err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
		}

		if description.Valid {
			t.Description = description.String
		}

		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}

		t.Priority = entity.MapToPriorityEntity(priorityStr)
		t.Status = entity.MapToStatusEntity(statusStr)

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (d *DB) UpdateTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	const op = OpUpdateTask

	query := `UPDATE tasks SET category_id = ?, title = ?, description = ?, due_date = ?, priority = ?, status = ?
	          WHERE id = ? AND user_id = ?`

	res, err := d.conn.Conn().ExecContext(ctx, query,
		task.CategoryID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Priority.String(),
		task.Status.String(),
		task.ID,
		task.UserID,
	)

	if err != nil {
		return entity.Task{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to update task").
			WithKind(richerror.KindUnexpected)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Task{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected").
			WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		return entity.Task{}, richerror.New(op).
			WithMessage("task not found or you don't have permission").
			WithKind(richerror.KindNotFound)
	}

	return d.GetTaskByID(ctx, task.ID, task.UserID)
}

func (d *DB) DeleteTask(ctx context.Context, id uint, userID uint) error {
	const op = OpDeleteTask

	query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`

	res, err := d.conn.Conn().ExecContext(ctx, query, id, userID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete task").
			WithKind(richerror.KindUnexpected)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected").
			WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("task not found or you don't have permission").
			WithKind(richerror.KindNotFound)
	}

	return nil
}