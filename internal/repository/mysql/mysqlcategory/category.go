package mysqlcategory

import (
	"context"
	"database/sql"
	"strings"
	"todoapp/internal/entity"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"
)

const (
	OpCreateCategory = richerror.Op("mysqlcategory.CreateCategory")
	OpGetCategoryByID=richerror.Op("mysqlcategory.GetCategoryByID")
	OpGetAllCategories=richerror.Op("mysqlcategory.GetAllCategories")
	OpUpdateCategory=richerror.Op("mysqlcategory.UpdateCategory")
	OpDeleteCategory=richerror.Op("mysqlcategory.DeleteCategory")
	 OpGetCategoryByName=richerror.Op("mysqlcategory.GetCategoryByName")
)


func (d *DB)CreateCategory(ctx context.Context,category entity.Category)(entity.Category,error){
	const op =OpCreateCategory

	query := `
		INSERT INTO categories (user_id, name, description, created_at)
		VALUES (?, ?, ?, NOW())
	`
	res,err:=d.conn.Conn().ExecContext(ctx,query,category.UserID,category.Name,category.Description)

	if err !=nil {
		if isDuplicateError(err) {
			return entity.Category{}, richerror.New(op).
				WithErr(err).
				WithMessage("category name already exists").
				WithKind(richerror.KindInvalid).
				WithMeta("name", category.Name)
		}
		return entity.Category{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create category").
			WithKind(richerror.KindUnexpected)
		
	}
	id,err :=res.LastInsertId()
	if err != nil {
	 return entity.Category{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get inserted id").
			WithKind(richerror.KindUnexpected)
	}
	category.ID=uint(id)
	return category,nil
}

func (d *DB) GetCategoryByID(ctx context.Context,id uint,userID uint)(entity.Category,error){
	const op=OpGetCategoryByID
	query:=`SELECT id, user_id, name , description , created_at FROM categories WHERE id= ? AND user_id=?`

	var category entity.Category
	var createdAt []uint8
	var description sql.NullString

	err:=d.conn.Conn().QueryRowContext(ctx,query,id,userID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&description,
		&createdAt,
	)

	if err !=nil{
		if err ==sql.ErrNoRows {
			return entity.Category{},richerror.New(op).WithErr(err).WithMessage("category not found").WithKind(richerror.KindNotFound)
		}

		return entity.Category{},richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)

	}

	if description.Valid {
		category.Description=description.String
	}

	return category,nil
}

func (d *DB) GetAllCategories(ctx context.Context, userID uint)([]entity.Category,error){
	const op=OpGetAllCategories

	query := `SELECT id ,user_id,name ,description, created_at  FROM  categories  WHERE user_id = ? ORDER BY name ASC`

	rows,err:=d.conn.Conn().QueryContext(ctx,query,userID)

	if err!=nil{
		return nil,richerror.New(op).WithErr(err).WithMessage("failed to get categories").WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var categories []entity.Category

	for rows.Next(){
		var category entity.Category
		var createdAt []uint8
		var description sql.NullString

		err:=rows.Scan(&category.ID,&category.UserID,&category.Name,&description,&createdAt)

		if  err !=nil {
			return  nil,richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
		}
		if description.Valid {
			category.Description=description.String
		}
		categories=append(categories, category)
	}
	return  categories,nil
}


func (d *DB) GetCategoryByName(ctx context.Context, userID uint,name string)(entity.Category,error){
	const op=OpGetCategoryByName

	query:= `SELECT id , user_id,name,description, created_at  FROM categories WHERE user_id= ? AND name =? `

	var category entity.Category
	var createdAt []uint8
	var description sql.NullString

	err:=d.conn.Conn().QueryRowContext(ctx,query,userID,name).Scan(&category.ID,&category.UserID,&category.Name,&description,&createdAt)

	if err!=nil{
		if err==sql.ErrNoRows{
			return entity.Category{},richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}
		return entity.Category{},richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}
	if description.Valid {
		category.Description = description.String
	}

	return category, nil
}

func (d *DB) UpdateCategory(ctx context.Context,category entity.Category)(entity.Category,error){
	const op=OpUpdateCategory

	query:=`UPDATE categories SET name=?,description=? WHERE id = ? AND user_id = ?`

	res , err :=d.conn.Conn().ExecContext(ctx,query,category.Name,category.Description,category.ID,category.UserID)
	
	if err != nil {
	if isDuplicateError(err) {
		return entity.Category{}, richerror.New(op).
			WithErr(err).
			WithMessage("category name already exists").
			WithKind(richerror.KindInvalid).
			WithMeta("name", category.Name)
		}
		return entity.Category{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to update category").
			WithKind(richerror.KindUnexpected)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Category{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected").
			WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		return entity.Category{}, richerror.New(op).
			WithMessage("category not found or you don't have permission").
			WithKind(richerror.KindNotFound)
	}

	return  category,nil

}

func (d *DB) DeleteCategory(ctx context.Context, id uint, userID uint) error {
	const op = OpDeleteCategory

	query := `DELETE FROM categories WHERE id = ? AND user_id = ?`

	res, err := d.conn.Conn().ExecContext(ctx, query, id, userID)

	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete category").
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
			WithMessage("category not found or you don't have permission").
			WithKind(richerror.KindNotFound)
	}

	return nil
}


func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "1062")
}