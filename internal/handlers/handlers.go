package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
    "gorm.io/driver/postgres"
	log "github.com/sirupsen/logrus"

	"avito/cmd/config"
)

type App struct {
	Cfg     config.Config
}

func NewApp() *App {
	cfg := config.ConfigSetup()
	return &App{
		Cfg:     cfg,
	}
}


type User struct {  
	gorm.Model
	ID int
	Name string `json:"name"`
	Surname string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age int `json:"age,omitempty"`
	Gender string `json:"gender,omitempty"`
	Country []string `json:"country,omitempty"`
}

type GenderResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type NationResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func (newUser *User) AddParams() {
	var genderResponse GenderResponse
	var ageResponse AgeResponse
	var nationResponse NationResponse

	// запрос возраста
	response, err := http.Get("https://api.agify.io/?name="+newUser.Name)
	if err != nil {
		log.Fatal(err)
	}
	// читаем ответ
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	// распарсим возраст

	json.Unmarshal([]byte(body), &genderResponse)

	newUser.Gender = genderResponse.Gender
	// запрос пола
	response, err = http.Get("https://api.genderize.io/?name="+newUser.Name)
	if err != nil {
		log.Fatal(err)
	}
	// читаем ответ
	body, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	// распарсим пол

	json.Unmarshal([]byte(body), &ageResponse)

	newUser.Age = ageResponse.Age
	// запрос национальности
	response, err = http.Get("https://api.nationalize.io/?name="+newUser.Name)
	if err != nil {
		log.Fatal(err)
	}
	// читаем ответ
	body, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	// распарсим национальности

	json.Unmarshal([]byte(body), &nationResponse)

	for _, nation := range nationResponse.Country {
		newUser.Country = append(newUser.Country, nation.CountryID)
	}
}

func (newApp *App) AddUser(c *gin.Context) {
	// создаем экземпляр пользователя
	var newUser User
	// распарсиваем входные данные нового пользователя
	log.Debug("decoding request")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	if err := json.Unmarshal([]byte(body), &newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	// обогащаем нового пользователя
	
	newUser.AddParams()
	
	//добавляем нового пользователя в таблицу
	log.Debug("adding to the DB")
	db, err := gorm.Open(postgres.Open(newApp.Cfg.PostgresDBAddr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})

	result := db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"запись не добавлена в базу данных": result.Error,
		})
	} else {
		//отправляем айдишник добавленного сегмента
		c.JSON(http.StatusCreated, gin.H{
			"id added": newUser.ID,
			"name added": newUser.Name,
		})
	}
	}
	
func (newApp *App) RemoveUser(c *gin.Context) {
	// создаем экземпляр пользователя
	var newUser User
	// распарсиваем данные для удаления
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	//убираем сегмент из таблицы
	db, err := gorm.Open(postgres.Open(newApp.Cfg.PostgresDBAddr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	result := db.Delete(&newUser, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"запись не удалена": result.Error,
		})
	} else {
		//отправляем имя удаленного
		c.JSON(http.StatusOK, gin.H{
			"name deleted": newUser.Name,
			"surname deleted": newUser.Surname,
		})
	}
}

func (newApp *App) UpdateUser(c *gin.Context) {
	
}

func (newApp *App) GetUsers(c *gin.Context) {
	// забираем из базы всех пользователей
	var users []User
	db, err := gorm.Open(postgres.Open(newApp.Cfg.PostgresDBAddr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	// установим пагинацию
	page := 1
	pageSize := 5
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	//отправляем ответ
	c.JSON(http.StatusOK, &users)
}

