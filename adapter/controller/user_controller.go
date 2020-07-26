package controller

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"trackpump/auth"
	"trackpump/usecase"
	"trackpump/usecase/exception"

	"github.com/labstack/echo"
)

const (
	templatesPath = "adapter/controller/templates/"
)

type userController struct {
	useCases    usecase.UseCases
	authService *auth.Auth
}

// UserController defines the controllers
type UserController interface {
	// API methods
	Create(c echo.Context) error

	Login(c echo.Context) error

	WeeklyReport(c echo.Context) error

	// Frontend methods
	HomePage(c echo.Context) error

	SignUp(c echo.Context) error

	ProcessSignUp(c echo.Context) error

	Admin(c echo.Context) error

	LoginPage(c echo.Context) error

	ProcessLogin(c echo.Context) error

	MeasurementPage(c echo.Context) error

	ProcessMeasurement(c echo.Context) error
}

// NewUsersController returns a new donors controller
func NewUsersController(useCases usecase.UseCases, authService *auth.Auth) UserController {
	return &userController{
		useCases:    useCases,
		authService: authService,
	}
}

func (u *userController) Create(c echo.Context) error {
	in := usecase.CreateAccountInput{}
	if err := c.Bind(&in); err != nil {
		return c.String(http.StatusInternalServerError, "invalid payload")
	}
	res, err := u.useCases.CreateAccount(&in)
	if err != nil {
		var e *exception.Error
		if errors.As(err, &e) {
			log.Println(e.Err)
			return c.JSON(e.Code, e)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	requestAuth := auth.RequestAuth{
		ID:    res.ID,
		Email: res.Email,
	}
	authorization, err := u.authService.GetToken(&requestAuth)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	c.Response().Header().Add("Authorization", authorization)
	return c.JSON(http.StatusCreated, res)
}

func (u *userController) Login(c echo.Context) error {
	in := usecase.LoginInput{}
	if err := c.Bind(&in); err != nil {
		return c.String(http.StatusInternalServerError, "invalid payload")
	}
	res, err := u.useCases.Login(&in)
	if err != nil {
		var e *exception.Error
		if errors.As(err, &e) {
			log.Println(e.Err)
			return c.JSON(e.Code, e)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	requestAuth := auth.RequestAuth{
		ID:    res.ID,
		Email: res.Email,
	}
	authorization, err := u.authService.GetToken(&requestAuth)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	c.Response().Header().Add("Authorization", authorization)
	return c.JSON(http.StatusCreated, res)
}

func (u *userController) WeeklyReport(c echo.Context) error {
	if err := u.useCases.RequestWeeklyReport(); err != nil {
		var e *exception.Error
		if errors.As(err, &e) {
			log.Println(e.Err)
			return c.JSON(e.Code, e)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Println("processing reports")
	return c.String(http.StatusOK, "processing")
}

func (u *userController) HomePage(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles(templatesPath + "home.html"))
	var html bytes.Buffer
	err := tmpl.Execute(&html, nil)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func (u *userController) SignUp(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles(templatesPath + "signUp.html"))
	var html bytes.Buffer
	err := tmpl.Execute(&html, nil)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func (u *userController) ProcessSignUp(c echo.Context) error {
	request := c.Request()
	email := request.FormValue("email")
	password := request.FormValue("password")
	name := request.FormValue("name")
	gender := request.FormValue("gender")
	birth := request.FormValue("birth")
	height := request.FormValue("height")
	genderAsNumber, err := strconv.Atoi(gender)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error on parsing gender to number: %s</h1>", err.Error()))
	}
	heightAsNumber, err := strconv.Atoi(height)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error on parsing height to number: %s</h1>", err.Error()))
	}
	createAccountInput := usecase.CreateAccountInput{
		Name:     name,
		Email:    email,
		Password: password,
		Birth:    birth,
		Height:   heightAsNumber,
		Gender:   genderAsNumber,
	}
	res, err := u.useCases.CreateAccount(&createAccountInput)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error on createing account: %s</h1>", err.Error()))
	}
	requestAuth := auth.RequestAuth{
		ID:    res.ID,
		Email: res.Email,
	}
	authorization, err := u.authService.GetToken(&requestAuth)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error on getting token: %s</h1>", err.Error()))
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/admin?authorization=%s", authorization))
}

func (u *userController) Admin(c echo.Context) error {
	token := c.Request().URL.Query().Get("authorization")
	claims, err := auth.GetClaims(token)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error on auth: %s</h1>", err.Error()))
	}
	state := struct {
		Email         string
		Authorization string
	}{
		claims["email"],
		token,
	}
	tmpl := template.Must(template.ParseFiles(templatesPath + "admin.html"))
	var html bytes.Buffer
	err = tmpl.Execute(&html, state)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error: %s</h1>", err.Error()))
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func (u *userController) LoginPage(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles(templatesPath + "login.html"))
	var html bytes.Buffer
	err := tmpl.Execute(&html, nil)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func (u *userController) ProcessLogin(c echo.Context) error {
	request := c.Request()
	email := request.FormValue("email")
	password := request.FormValue("password")
	loginInput := usecase.LoginInput{
		Email:    email,
		Password: password,
	}
	res, err := u.useCases.Login(&loginInput)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error: %s</h1>", err.Error()))
	}
	requestAuth := auth.RequestAuth{
		ID:    res.ID,
		Email: res.Email,
	}
	authorization, err := u.authService.GetToken(&requestAuth)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error on getting token: %s</h1>", err.Error()))
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/admin?authorization=%s", authorization))

}

func (u *userController) MeasurementPage(c echo.Context) error {
	token := c.Request().FormValue("token")
	tmpl := template.Must(template.ParseFiles(templatesPath + "newMeasurement.html"))
	var html bytes.Buffer
	state := struct {
		Authorization string
	}{
		token,
	}
	err := tmpl.Execute(&html, state)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func (u *userController) ProcessMeasurement(c echo.Context) error {
	token := c.FormValue("token")
	claims, err := auth.GetClaims(token)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error %s</h1>", err.Error()))
	}
	id := claims["id"]
	request := c.Request()
	weight, err := strconv.Atoi(request.FormValue("weight"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert weight: %s</h1>", err.Error()))
	}
	abdominalCircunference, err := strconv.Atoi(request.FormValue("abdominalCircunference"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert abdominal circunference: %s</h1>", err.Error()))
	}
	arm, err := strconv.Atoi(request.FormValue("arm"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert arm: %s</h1>", err.Error()))
	}
	forearm, err := strconv.Atoi(request.FormValue("forearm"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert arm: %s</h1>", err.Error()))
	}
	calf, err := strconv.Atoi(request.FormValue("calf"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert calf: %s</h1>", err.Error()))
	}
	neck, err := strconv.Atoi(request.FormValue("neck"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert neck: %s</h1>", err.Error()))
	}
	hip, err := strconv.Atoi(request.FormValue("hip"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert hip: %s</h1>", err.Error()))
	}
	thigh, err := strconv.Atoi(request.FormValue("thigh"))
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on convert thigh: %s</h1>", err.Error()))
	}
	fontalFile, _, err := request.FormFile("frontalPicture")
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on reading frontalPicture: %s</h1>", err.Error()))
	}
	fontalFileFileBytes, err := ioutil.ReadAll(fontalFile)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on reading bytes of frontalPicture: %s</h1>", err.Error()))
	}
	fontalFileBuffer := new(bytes.Buffer)
	if _, err := fontalFileBuffer.Write(fontalFileFileBytes); err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on writing bytes of frontal picture: %s</h1>", err.Error()))
	}
	fontalFileContent, err := ioutil.ReadAll(fontalFileBuffer)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on reading bytes of frontal picture: %s</h1>", err.Error()))
	}
	frontalPictureAsBase64 := base64.StdEncoding.EncodeToString(fontalFileContent)
	sideFile, _, err := request.FormFile("sidePicture")
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on reading frontalPicture: %s</h1>", err.Error()))
	}
	sideFileFileBytes, err := ioutil.ReadAll(sideFile)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on reading bytes of frontalPicture: %s</h1>", err.Error()))
	}
	sideFileBuffer := new(bytes.Buffer)
	if _, err := sideFileBuffer.Write(sideFileFileBytes); err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on writing bytes of frontal picture: %s</h1>", err.Error()))
	}
	sideFileContent, err := ioutil.ReadAll(sideFileBuffer)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Error on reading bytes of frontal picture: %s</h1>", err.Error()))
	}
	sidePictureAsBase64 := base64.StdEncoding.EncodeToString(sideFileContent)
	in := usecase.RegisterMeasurementInput{
		ID:                     id,
		Weight:                 weight,
		AbdominalCircunference: abdominalCircunference,
		Arm:                    arm,
		Forearm:                forearm,
		Calf:                   calf,
		Neck:                   neck,
		Hip:                    hip,
		Thigh:                  thigh,
		FrontalPicture:         frontalPictureAsBase64,
		SidePicture:            sidePictureAsBase64,
	}
	err = u.useCases.RegisterMeasurement(&in)
	if err != nil {
		log.Println(err)
		var e *exception.Error
		if errors.As(err, &e) {
			log.Println(e.Err)
			return c.JSON(e.Code, e)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/admin?authorization=%s", token))
}
