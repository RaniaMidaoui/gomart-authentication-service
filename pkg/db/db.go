package db

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/models"
	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{})

	return Handler{db}
}

func Mock() Handler {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       mockDb,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	columns := []string{"id", "email", "password"}
	mock.ExpectQuery("SELECT (.+)").
		WithArgs("test@admida0ui.tech").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "test@admida0ui.tech", utils.HashPassword("P@ssw0rd")))

	mock.ExpectBegin()
	mock.ExpectRollback()

	return Handler{db}
}
