package gorm

import "context"

type CreateUserParams struct {
	Username       string `json:"username"`
	PasswordHashed string `json:"password_hashed"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (store *GormStore) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	user := User{
		Username:       arg.Username,
		PasswordHashed: arg.PasswordHashed,
		FullName:       arg.FullName,
		Email:          arg.Email,
	}

	result := store.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

func (store *GormStore) GetUser(ctx context.Context, username string) (User, error) {
	var user User
	result := store.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}
