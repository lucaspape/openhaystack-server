package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

type Config struct {
	ApiPrefix   string
	HaystackDir string
	ApiKey      string
}

var config Config

func main() {
	f, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := io.ReadAll(f)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(b, &config)

	if err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	registerHandlers(r)

	err = r.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func registerHandlers(r *gin.Engine) {
	r.GET(config.ApiPrefix+"/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
	})

	r.GET(config.ApiPrefix+"/accessories", func(c *gin.Context) {
		k := c.Query("key")

		if k != config.ApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong api key"})
			return
		}

		files, err := os.ReadDir(config.HaystackDir)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		var r []map[string]interface{}

		for _, file := range files {
			f, err := os.Open(config.HaystackDir + file.Name())

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			b, err := io.ReadAll(f)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			m := make(map[string]interface{})

			err = json.Unmarshal(b, &m)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			r = append(r, m)
		}

		c.JSON(http.StatusOK, gin.H{"results": r})
	})
}
