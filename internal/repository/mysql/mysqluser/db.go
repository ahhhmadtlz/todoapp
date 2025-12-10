package mysqluser

import (
	"todoapp/internal/repository/mysql"
)

type DB struct {
	conn *mysql.MySQLDB
}

// GetUserByPhonenumber implements userservice.Repository.
// func (d *DB) GetUserByPhonenumber(ctx context.Context, phone string) (entity.User, error) {
// 	panic("unimplemented")
// }

// RegisterUser implements userservice.Repository.
// func (d *DB) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
// 	panic("unimplemented")
// }

func New(conn *mysql.MySQLDB) *DB {
	return &DB{conn: conn}
}
