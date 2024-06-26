package Domain

import (
	"errors"
	"fmt"
	model2 "github.com/TeenBanner/Inventory_system/User/Domain/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UserStorage
}

func NewUser(storage UserStorage) *User {
	return &User{
		UserStorage: storage,
	}
}

func (u *User) CreateUser(user model2.User) error {
	user.CreatedAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.IsAdmin = false

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = u.UserStorage.PsqlCreateUser(user)
	if err != nil {
		return err
	}

	fmt.Println("User created")

	return nil
}

func (u *User) GetUserByEmail(email string) (model2.User, error) {
	user, err := u.UserStorage.PsqlGetUserByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) GetUserByName(name string) (model2.User, error) {
	if name == "" {
		return model2.User{}, errors.New("search name cannot be empty")
	}
	user, err := u.UserStorage.PsqlGetUserByName(name)

	if err != nil {
		return model2.User{}, err
	}

	return user, nil
}

func (u *User) UpdateUserEmail(ActualEmail, NewEmail string) error {
	if NewEmail == ActualEmail {
		return errors.New("user email can't be equal")
	}

	if NewEmail == "" {
		return errors.New("user email can't be empty")
	}

	err := u.UserStorage.PsqlUpdateUserEmail(ActualEmail, NewEmail)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateUserName(email, NewName string) error {
	if email == "" {
		return errors.New("user email can't be empty")
	}

	if NewName == "" {
		return errors.New("user name can't be empty")
	}

	err := u.UserStorage.PsqlUpdateUserName(email, NewName)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateUserPassword(email, NewPassword string) error {
	if email == "" || NewPassword == "" {
		return errors.New("please provide a valid email or password")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Error while hashing password")
	}

	err = u.UserStorage.PsqlUpdateUserPassword(email, string(hashPassword))
	if err != nil {
		return err
	}

	return nil
}

// GetAllUsers admin functiond
func (u *User) GetAllUsers() ([]model2.User, error) {
	users, err := u.PsqlGetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) FindEmailByName(name string) (string, error) {
	if name == "" {
		return "", errors.New("please provide a valid searchName")
	}

	email, err := u.UserStorage.PsqlFindUserEmailByName(name)
	if err != nil {
		return "", err
	}

	return email, err
}

func (u *User) DeleteAccount(email string) error {
	if email == "" {
		return errors.New("please provide a valid email")
	}

	err := u.UserStorage.PsqlDeleteAccount(email)
	if err != nil {
		return err
	}

	return nil
}
