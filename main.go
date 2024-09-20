package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	dsn := "postgresql://postgres:1234@localhost/postgres?sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"uniqueIndex"`
}

func getUsers(c echo.Context) error {
	var users []User
	db.Find(&users)
	return c.JSON(200, users)
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(c echo.Context) error {
	var loginUser LoginUser
	c.Bind(&loginUser)

	var user User

	if err := db.Where("email = ?", loginUser.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, "unauthrorized")
		return err
	}

	if user.Password != loginUser.Password {
		return c.JSON(http.StatusUnauthorized, "unauthrorized")
	}

	return c.JSON(http.StatusOK, "you can login")
}

func registerUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	db.Create(&user)
	return c.JSON(http.StatusCreated, user)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/auth/register", registerUser)
	e.POST("/auth/login", login)
	e.GET("/users", getUsers)

	e.Logger.Fatal(e.Start(":8080"))

}
