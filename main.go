package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var path = "data.json"

func createFile() {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}

	fmt.Println("File Created Successfully", path)
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

type Transaction struct {
	Amount    float32
	Timestamp string
}

func main() {
	r := gin.Default()
	createFile()
	r.GET("/statistics", func(c *gin.Context) {
		max := 0
		min := 0
		sum := 0
		avg := 0
		count := 0
		jsonFile, err := os.Open("data.json")
		if err != nil {
			fmt.Println(err)
		}
		currentTime := time.Now().UTC()

		var trans = []Transaction{}

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &trans)
		count = len(trans)
		if len(trans) > 0 {
			min = int(trans[0].Amount)
		}
		for i := 0; i < len(trans); i++ {
			date, error := time.Parse("2006-01-02T15:04:05.000Z", trans[i].Timestamp)
			difference := date.Sub(currentTime)
			if difference.Seconds() > 60 {
				if error != nil {
					fmt.Println(error)
				}
				if max < int(trans[i].Amount) {
					max = int(trans[i].Amount)
				}
				if min > int(trans[i].Amount) {
					min = int(trans[i].Amount)
				}
				sum = sum + int(trans[i].Amount)
			}
		}
		avg = sum / count
		c.JSON(http.StatusOK, gin.H{"max": max, "min": min, "sum": sum, "avg": avg, "count": count})
	})

	r.POST("/transactions", func(c *gin.Context) {
		var input Transaction
		currentTime := time.Now().UTC()
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//read all from file
		jsonFile, err := os.Open("data.json")

		if err != nil {
			fmt.Println(err)
		}
		date, error := time.Parse("2006-01-02T15:04:05.000Z", input.Timestamp)
		fmt.Println(date)
		if error != nil {
			fmt.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{"error": error})
			return
		}

		difference := date.Sub(currentTime)
		if difference.Seconds() > 60 {
			fmt.Println(difference.Seconds())

			c.JSON(http.StatusNoContent, gin.H{"error": "date time should be lesss than 60 sec"})
			return
		}
		if difference.Hours() < -1 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "date time should be lesss than 60 sec"})
			return
		}

		var trans = []Transaction{}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &trans)

		//write to file
		inputs := append([]Transaction{input}, trans...)

		content, err := json.Marshal(inputs)
		if err != nil {
			// Handle error
			fmt.Println(err)
		}

		err = ioutil.WriteFile("data.json", content, 0644)

		c.JSON(http.StatusCreated, gin.H{"data": &input})

	})
	r.DELETE("/transactions", func(c *gin.Context) {
		var trans = []Transaction{}
		content, err := json.Marshal(trans)
		if err != nil {
			// Handle error
			c.JSON(http.StatusBadRequest, gin.H{"error": err})

		}

		err = ioutil.WriteFile("data.json", content, 0644)
		c.JSON(http.StatusNoContent, gin.H{})

	})
	r.Run()
}
