package psqlUser

import (
	"database/sql"
	"errors"
	models2 "github.com/TeenBanner/Inventory_system/User/Domain/model"
	"github.com/TeenBanner/Inventory_system/pkg/database"
	"log"
	"time"
)

// UserStorage it's used for interact with DB
type userStorage struct {
	db *sql.DB
}

// NewUserStorage contructure for UserStorage
func NewPsqlUser(DB *sql.DB) *userStorage {
	return &userStorage{
		db: DB,
	}
}

// User methods
func (u *userStorage) PsqlCreateUser(user models2.User) error {
	stmt, err := u.db.Prepare(SqlCreateUserQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()

	UserNullTime := database.TimeToNull(user.UpdatedAt)

	_, err = stmt.Exec(
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		UserNullTime,
	)

	if err != nil {
		return err
	}

	log.Println("Usuario creado")
	return nil
}

// GetUser get info from a user
func (u *userStorage) PsqlGetUserByEmail(email string) (models2.User, error) {
	stmt, err := u.db.Prepare(SqlGetUser)
	if err != nil {
		return models2.User{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)
	if row == nil {
		return models2.User{}, errors.New("user does not exist")
	}
	user := models2.User{}
	nulltime := sql.NullTime{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &nulltime)
	user.UpdatedAt = nulltime.Time
	if err != nil {
		return models2.User{}, err
	}

	return user, nil
}

func (U *userStorage) PsqlGetUserName(email string) (string, error) {
	stmt, err := U.db.Prepare(SqlGetUserName)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	row, err := stmt.Query(email)
	if err != nil {
		return "", err
	}

	var name string
	err = row.Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (u *userStorage) PsqlUpdateUserName(email, name string) error {
	stmt, err := u.db.Prepare(SqlUpdateUserName)
	if err != nil {
		return err
	}

	defer stmt.Close()
	update_time := time.Now()
	_, err = stmt.Exec(name, update_time, email)
	if err != nil {
		return err
	}
	return nil
}

func (u *userStorage) PsqlUpdateUserEmail(ActualEmail, NewEmail string) error {
	stmt, err := u.db.Prepare(SqlUpdateUserEmail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	update_time := time.Now()
	_, err = stmt.Exec(NewEmail, update_time, ActualEmail)
	if err != nil {
		return err
	}

	return nil
}

func (U *userStorage) PsqlUpdateUserPassword(Email, NewPassword string) error {
	stmt, err := U.db.Prepare(SqlUpdateUserPassword)
	if err != nil {
		return err
	}

	defer stmt.Close()
	update_time := time.Now()
	_, err = stmt.Exec(NewPassword, update_time, Email)
	if err != nil {
		return err
	}

	return nil
}

// AdminMethods

func (u *userStorage) PsqlGetAllUsers() ([]models2.User, error) {
	stmt, err := u.db.Prepare(SqlAdminGetAllUsers)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	users := []models2.User{}
	for rows.Next() {
		user := models2.User{}
		nullTime := sql.NullTime{}

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &nullTime)
		user.UpdatedAt = nullTime.Time

		users = append(users, user)
		if err != nil {
			return nil, err
		}
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (U *userStorage) PsqlGetUserByName(name string) (models2.User, error) {
	stmt, err := U.db.Prepare(SqlGetUserByName)
	if err != nil {
		return models2.User{}, err
	}
	defer stmt.Close()

	user := models2.User{}

	row := stmt.QueryRow(name)

	if row == nil {
		return models2.User{}, err
	}

	nulltime := sql.NullTime{}

	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &nulltime)
	if err != nil {
		return models2.User{}, err
	}

	user.Password = ""
	user.UpdatedAt = nulltime.Time

	return user, nil
}

func (U *userStorage) PsqlLoginGetEmail(email string) (string, error) {
	stmt, err := U.db.Prepare(SqlLoginCompareEmails)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	var DbEmail string
	row := stmt.QueryRow(email)
	err = row.Scan(&DbEmail)
	if err != nil && DbEmail == "" {
		return "", err
	}

	return DbEmail, nil
}

func (U *userStorage) PsqlLoginGetPassword(email string) (string, error) {
	stmt, err := U.db.Prepare(SqlLoginGetHashdPasswordWithEmail)
	if err != nil {
		return "", err
	}

	defer stmt.Close()
	var HashPassword string

	row := stmt.QueryRow(email)
	err = row.Scan(&HashPassword)
	if err != nil && HashPassword == "" {
		return "", err
	}

	return HashPassword, nil
}

func (U *userStorage) PsqlFindUserEmailByName(name string) (string, error) {
	stmt, err := U.db.Prepare(SqlFindUserEmailByName)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	var email string

	row := stmt.QueryRow(name)

	err = row.Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func (U *userStorage) PsqlDeleteAccount(email string) error {
	stmt, err := U.db.Prepare(SqlDeleteAccount)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}

	return nil
}
