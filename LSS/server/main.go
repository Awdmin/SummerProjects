package main

import(
	"fmt"
	"net/http"
	"database/sql"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_"github.com/glebarez/go-sqlite"
)

const ntfUrl string = "https://ntfy.sh/<my-topic>"

func signUpNtf(name string, date string) error {
	d, _ := time.Parse("2006-01-02", date)
	msg := fmt.Sprintf("%s se je najavil/a na kosilo za dan %s", name, d.Format("02.01.2006"))
	return sendNotification(msg, "Najava")
}

func deleteNtf(name string, date string) error {
	d, _ := time.Parse("2006-01-02T15:04:05Z", date)
	msg := fmt.Sprintf("%s se je odjavil/a od kosilo za dan %s", name, d.Format("02.01.2006"))
	return sendNotification(msg, "Odjava")
}

func sendNotification(message string, t string) error {
    req, err := http.NewRequest("POST", ntfUrl, strings.NewReader(message))
    if err != nil {
        return err
    }

    req.Header.Set("Title", t)
    req.Header.Set("Priority", "high")       // min, low, default, high, urgent
	req.Header.Set("Icon", "https://cdn-icons-png.flaticon.com/512/857/857681.png")

    _, err = http.DefaultClient.Do(req)
    return err
}

type entry struct {
	ID		string	`json:"id"`
	Name	string	`json:"name"`
	Date	string	`json:"date"`
}

type Handler struct {
	DB 		*sql.DB
}

func (h *Handler) initDB(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS entrys (
		id		INTEGER 	PRIMARY KEY AUTOINCREMENT,
		name	TEXT 		NOT NULL,
		date	DATE		NOT NULL
	);`

	return h.DB.Exec(sql)
}


func (h *Handler) getEntries(c *gin.Context) {
	sql := `SELECT * FROM entrys;`

	rows, err := h.DB.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var entries []entry
	for rows.Next() {
		e := &entry{}
		err := rows.Scan(&e.ID, &e.Name, &e.Date)
		if err != nil {
			fmt.Println(err)
			return 
		}

		entries = append(entries, *e)
	}

	//fmt.Println(entries)

	c.IndentedJSON(http.StatusOK, entries)
}

func (h *Handler) postEntry(c *gin.Context) {
	var newEntry entry

	if err := c.BindJSON(&newEntry); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	sql := `INSERT INTO entrys (name, date)
		VALUES (?, ?);`
	_, err := h.DB.Exec(sql, newEntry.Name, newEntry.Date)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	fmt.Println("Successfully inserted a new entry.")

	c.IndentedJSON(http.StatusOK, newEntry)

	signUpNtf(newEntry.Name, newEntry.Date)
}

func (h *Handler) deleteEntry(c *gin.Context) {
	var delEntry entry

	if err := c.BindJSON(&delEntry); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	sql := `DELETE FROM entrys 
		WHERE id = ?;`
	_, err := h.DB.Exec(sql, delEntry.ID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	fmt.Println("Successfully inserted a new entry.")

	c.IndentedJSON(http.StatusOK, delEntry)

	deleteNtf(delEntry.Name, delEntry.Date)
}


func main() {

	db, err := sql.Open("sqlite", "./my.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	fmt.Println("Connection to the database was successful.")
	h := &Handler{DB: db}

	_, err = h.initDB(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/entries", h.getEntries)
	router.POST("/entry", h.postEntry)
	router.DELETE("/entry", h.deleteEntry)

	router.Run(":8111")

}
