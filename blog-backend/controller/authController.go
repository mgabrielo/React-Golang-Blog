package controller

import (
	"fmt"
	"github/mgabrielo/React-Golang-Blog/database"
	"github/mgabrielo/React-Golang-Blog/models"
	"github/mgabrielo/React-Golang-Blog/util"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("unable to parse body")
	}
	// check if password length is greater than 6
	if len(data["password"].(string)) <= 5 {
		c.Status(400)
		return c.JSON(fiber.Map{"message": "Password must be greater than six"})
	}
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{"message": "Invalid Email"})
	}

	//  check if email already exist
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{"message": "Email Already Exist"})
	}

	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Phone:     data["phone"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Password:  nil,
	}

	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{"user": user, "message": "Account creation Successful"})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("unable to parse body")
	}
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{"message": "Email Does not Exist"})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(403)
		return c.JSON(fiber.Map{"message": "incorrect password"})
	}

	token, err := util.GenerateJWt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.JSON(fiber.StatusInternalServerError)
		return nil
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "Login Successful", "user": user})

}

type Claim struct {
	jwt.StandardClaims
}
