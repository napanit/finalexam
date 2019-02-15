package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
//	"github.com/napanit/myapi/customer"
//	"strconv"
//	"fmt"
	"database/sql"
	"os"
	"log"
	_"github.com/lib/pq"
)
type Item struct {
	ID     int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Status string`json:"status"`
}

var Items []Item
var DB *sql.DB

func main () {
	
	var err error 
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil{
		log.Fatal(err)
	}
    createDB()
	r := gin.Default()
	r.Use(loginMiddleware)
	r.POST("/customers",postCustomerHandler)
	r.GET("/customers/:id",getCustomerByIDHandler)
	r.GET("/customers",getCustomerAllHandler)	
	r.PUT("/customers/:id",putCustomerHandler)
	r.DELETE("/customers/:id",delCustomerIDHandler)
//	r.D
	r.Run(":2019")
}

func postCustomerHandler(c *gin.Context) {
	 var newItem Item
	 err := c.ShouldBindJSON(&newItem)
	 if err != nil{
		 c.String(http.StatusBadRequest, err.Error())
		 return
	 }
	 row := DB.QueryRow("INSERT into customers(name, email, status) values ($1,$2,$3) RETURNING id",newItem.Name, newItem.Email, newItem.Status)
	 var id int
	 err = row.Scan(&id)
	 if err != nil{
		 c.String(http.StatusInternalServerError, err.Error())
		 return
	 }
	 newItem.ID = id
	 c.JSON(201, newItem)
	 }

func getCustomerByIDHandler(c *gin.Context){
	 id := c.Param("id")
	 stmt, err := DB.Prepare("SELECT id, name, email, status FROM customers where id = $1 ")
	 if err != nil {
		 c.String(http.StatusInternalServerError, err.Error())
		 return
	 }
	 row := stmt.QueryRow(id)
	 var item Item
	 err = row.Scan(&item.ID, &item.Name, &item.Email, &item.Status)
	 if err != nil {
		 c.String(http.StatusInternalServerError, err.Error())
		 return
	 }
	 c.JSON(200,item)

}

func getCustomerAllHandler(c *gin.Context) {
	stmt, err := DB.Prepare("SELECT id, name, email, status FROM Customers")
	if err != nil{
		c.String(http.StatusInternalServerError, err.Error())
		return 
	}
	rows, err := stmt.Query() 
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} 
	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.Status)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, item)
	}
	c.JSON(200,items)
}

func putCustomerHandler(c *gin.Context){
	 id := c.Param("id")
	 var newItem Item
	 err := c.ShouldBindJSON(&newItem)
	 if err != nil{
		 c.String(http.StatusBadRequest, err.Error())
		 return
	 }
	 stmt, err := DB.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1")
	 if err != nil{
		 c.String(http.StatusBadRequest, err.Error())
		 return 
	 }
	 _, err = stmt.Exec(id, newItem.Name, newItem.Email, newItem.Status)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return

}
stmt, err = DB.Prepare("SELECT id, name, email, status FROM customers where id = $1 ")
if err != nil {
	c.String(http.StatusInternalServerError, err.Error())
	return
}
row := stmt.QueryRow(id)
var item Item
err = row.Scan(&item.ID, &item.Name, &item.Email, &item.Status)
if err != nil {
	c.String(http.StatusInternalServerError, err.Error())
	return
}	
     c.JSON(200,item)
}

func delCustomerIDHandler(c *gin.Context) {
	id := c.Param("id")

	stmt, err := DB.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "customer deleted",
	})

}

func loginMiddleware(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Status code is 401 Unauthorized")
		c.Abort()
		return
	}
	c.Next()
	log.Println("ending middleware")
}

func createDB() {
var err  error
createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
    ` 
	_,err = DB.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table",err)
	}	
}