package mysqluser

import (
	"context"
	"database/sql"
	"strings"
	"todoapp/internal/entity"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"
	"todoapp/internal/repository/mysql"
)


const (
	OpRegister=richerror.Op("mysqluser.Register")
	OpGetUserByID          = richerror.Op("mysql.GetUserByID")
	OpGetUserByPhoneNumber = richerror.Op("mysql.GetUserByPhoneNumber")
	OpIsPhoneNumberUnique=richerror.Op("mysqluser.IsPhoneNumberUnique")
)




func (d *DB) IsPhoneNumberUnique(ctx context.Context,phonenumber string) (bool, error) {

	const op=OpIsPhoneNumberUnique

	var exists bool

		err :=d.conn.Conn().QueryRowContext(ctx,`SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = ?)`,phonenumber).Scan(&exists)


	if err!=nil{
		return false,richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return !exists, nil
}


func (d *DB)RegisterUser(ctx context.Context ,user entity.User)(entity.User,error){
	const op=OpRegister

	query := `INSERT INTO users (name, phone_number, password, role, created_at) 
	          VALUES (?, ?, ?, ?, NOW())`
						
	res, err := d.conn.Conn().ExecContext(ctx, query,
			user.Name,
			user.PhoneNumber,
			user.Password,
			user.Role.String(), // Convert role to string for ENUM
		)


		if err!=nil{
			if isDuplicateKeyError(err) {
					return entity.User{}, richerror.New(op).
							WithErr(err).
							WithMessage("phone number already exists").
							WithKind(richerror.KindInvalid)
				}

			return entity.User{},richerror.New(op).WithErr(err).WithMessage("cant excute command").WithKind(richerror.KindUnexpected)
	}



  id, err := res.LastInsertId()
  if err != nil {
      return entity.User{}, richerror.New(op).
          WithErr(err).
          WithMessage("failed to get inserted id").
          WithKind(richerror.KindUnexpected)
  }
    
    user.ID = uint(id)

    return user, nil
}

func (d *DB) GetUserByPhoneNumber(ctx context.Context, phoneNumber string)(entity.User,error){
	const op=OpGetUserByPhoneNumber

	row:=d.conn.Conn().QueryRowContext(ctx,`select * from users where phone_number = ?`,phoneNumber)

	user,err:=scanUser(row)

	if err!=nil{
		if err==sql.ErrNoRows{
			return entity.User{},richerror.New(op).
			WithErr(err).
			WithMessage(errmsg.ErrorMsgNotFound).
			WithKind(richerror.KindNotFound)
		}
		return entity.User{}, richerror.New(op).
			WithErr(err).
				WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
	}
return user,nil
}


func (d *DB) GetUserByID(ctx context.Context,userID uint)(entity.User,error){
	const op=OpGetUserByID

	row :=d.conn.Conn().QueryRowContext(ctx,`select * from users where id = ? `,userID)

	user,err:=scanUser(row)
	
	if err !=nil {
		if err==sql.ErrNoRows {
			return entity.User{},richerror.New(op).
			WithErr(err).
			WithMessage(errmsg.ErrorMsgNotFound).
			WithKind(richerror.KindNotFound)
		}
		return entity.User{},richerror.New(op).
		WithMessage(errmsg.ErrorMsgCantScanQueryResult).
		WithKind(richerror.KindUnexpected)
	}

	return user ,nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User
	var createdAt []uint8

	var roleStr string

	err := scanner.Scan(
      &user.ID,           // id
      &user.Name,         // name
      &user.PhoneNumber,  // phone_number
      &roleStr,           // role (ENUM stored as string)
      &user.Password,     // password
      &createdAt,         // created_at
    )
		
	if err != nil {
		return entity.User{}, err
	}
    
	user.Role=entity.MapToRoleEntity(roleStr)

	return user,err
}

func isDuplicateKeyError(err error) bool {
    return strings.Contains(err.Error(), "Duplicate entry") ||
           strings.Contains(err.Error(), "1062")
}