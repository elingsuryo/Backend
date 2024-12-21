package user_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/repository/user"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestRepositoryImpl_Create(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	// Define expected user
	expectedUser := &entity.User{
		ID:                 "123",
		Email:              "test@example.com",
		Password:           "hashedpassword",
		Name:               "Test User",
		Role:               "admin",
		Status:             true,
		LastLogin:          nil,
		VerifyEmailToken:   "",
		IsVerified:         false,
		ResetPasswordToken: "",
	}

	// Mock the query for Create
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`id`,`email`,`password`,`name`,`role`,`status`,`last_login`,`verify_email_token`,`reset_password_token`,`is_verified`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(
			expectedUser.ID,
			expectedUser.Email,
			expectedUser.Password,
			expectedUser.Name,
			expectedUser.Role,
			expectedUser.Status,
			expectedUser.LastLogin,
			expectedUser.VerifyEmailToken,
			expectedUser.ResetPasswordToken,
			expectedUser.IsVerified,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	err = repo.Create(gormDB, expectedUser)

	// Assertions
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}
func TestRepositoryImpl_Update(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	// Define expected user
	expectedUser := &entity.User{
		ID:                 "123",
		Email:              "test@example.com",
		Password:           "hashedpassword",
		Name:               "Test User",
		Role:               "admin",
		Status:             true,
		LastLogin:          nil,
		VerifyEmailToken:   "",
		IsVerified:         false,
		ResetPasswordToken: "",
	}

	// Mock the query for Update
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `email`=?,`password`=?,`name`=?,`role`=?,`status`=?,`last_login`=?,`verify_email_token`=?,`reset_password_token`=?,`is_verified`=?,`created_at`=?,`updated_at`=?,`deleted_at`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs(
			expectedUser.Email,
			expectedUser.Password,
			expectedUser.Name,
			expectedUser.Role,
			expectedUser.Status,
			expectedUser.LastLogin,
			expectedUser.VerifyEmailToken,
			expectedUser.ResetPasswordToken,
			expectedUser.IsVerified,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			expectedUser.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	err = repo.Update(gormDB, expectedUser)

	// Assertions
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestRepositoryImpl_Delete(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	// Define expected user
	expectedUser := &entity.User{
		ID: "123",
	}

	// Mock the query for Delete
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), expectedUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	err = repo.Delete(gormDB, expectedUser)

	// Assertions
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestUserRepository_GetFirst(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup repository and logger
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	expectedUser := &entity.User{ID: "123", Email: "test@example.com"}
	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(expectedUser.ID, expectedUser.Email)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(1).
		WillReturnRows(rows)

	data := &entity.User{}
	err = repo.GetFirst(gormDB, data)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, data.ID)
	assert.Equal(t, expectedUser.Email, data.Email)
	mock.ExpectationsWereMet()
}

func TestUserRepository_GetByID(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup repository and logger
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	// Define expected behavior
	expectedUser := &entity.User{ID: "123", Email: "test@example.com"}
	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(expectedUser.ID, expectedUser.Email)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL LIMIT ?")).
		WithArgs("123", 1).
		WillReturnRows(rows)

	// Call the method
	data := &entity.User{}
	err = repo.GetByID(gormDB, data, "123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, data.ID)
	assert.Equal(t, expectedUser.Email, data.Email)
	mock.ExpectationsWereMet()
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup repository and logger
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	// Define expected data data
	expectedUser := &entity.User{ID: "123", Email: "test@example.com"}
	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(expectedUser.ID, expectedUser.Email)

	// Mock the query for GetByEmail
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL LIMIT ?")).
		WithArgs("test@example.com", 1).
		WillReturnRows(rows)

	// Call the method
	data := &entity.User{}
	err = repo.GetByEmail(gormDB, data, "test@example.com")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, data.ID)
	assert.Equal(t, expectedUser.Email, data.Email)
	mock.ExpectationsWereMet()
}

func TestUserRepository_CountByRole(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup repository and logger
	logger := logrus.New()
	repo := user.NewUserRepository(gormDB, logger)

	// Define expected behavior
	expectedCount := int64(10)
	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE role = ?")).
		WithArgs("buyer").
		WillReturnRows(rows)

	// Call the method
	count, err := repo.CountByRole(gormDB, "buyer")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mock.ExpectationsWereMet()
}
