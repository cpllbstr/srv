package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cpllbstr/gogrn/slv"
	"github.com/gin-gonic/gin"
)

type login struct {
	Name    string
	Address string
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("./www/index.html")
	})
	r.GET("/sim", func(c *gin.Context) {
		c.File("./www/sim.html")
	})
	r.GET("/2d", func(c *gin.Context) {
		c.File("./www/calc2d.html")
	})
	r.GET("/3d", func(c *gin.Context) {
		c.File("./www/calc3d.html")
	})

	r.POST("/sim", func(c *gin.Context) {

		var m [3]float64
		m[0], _ = strconv.ParseFloat(c.PostForm("m1"), 64)
		m[1], _ = strconv.ParseFloat(c.PostForm("m2"), 64)
		m[2], _ = strconv.ParseFloat(c.PostForm("m3"), 64)

		var k [3]float64
		k[0], _ = strconv.ParseFloat(c.PostForm("k1"), 64)
		k[1], _ = strconv.ParseFloat(c.PostForm("k2"), 64)
		k[2], _ = strconv.ParseFloat(c.PostForm("k3"), 64)

		v, _ := strconv.ParseFloat(c.PostForm("vel"), 64)
		l, _ := strconv.ParseFloat(c.PostForm("len"), 64)

		step, _ := strconv.ParseFloat(c.PostForm("stpt"), 64)
		fint, _ := strconv.ParseFloat(c.PostForm("fint"), 64)
		c.JSON(200, m)
		c.JSON(200, k)
		c.JSON(200, v)
		c.JSON(200, l)
		avg := c.Request.FormValue("avg")

		body := slv.NewThreeBodyModel(m, k)
		y, mon, d := time.Now().Date()
		file, err := os.Create(fmt.Sprintf("./dat/sym%v:%v_%v-%v-%v.dat", time.Now().Hour(), time.Now().Minute(), d, mon, y)) //файл в который пишутся результаты расчетов
		if err != nil {
			log.Println("Cannot create file!  :", err)
			c.File("./www/error.html")
			return
		}
		machine := slv.StateMachFromModel(body, v, l, step, 0, fint)
		if avg == "on" {
			slv.SimulateAv(machine, body, file)
		} else {
			slv.Simulate(machine, file)
		}
	})
	r.POST("/2d", func(c *gin.Context) {
		//вот тут надо запускать функцию симуляции и парсить что пришло из формы
		log := login{
			Name:    c.PostForm("name"),
			Address: c.PostForm("address"),
		}
		c.JSON(200, gin.H{"status": log})
	})
	r.POST("/3d", func(c *gin.Context) {
		//вот тут надо запускать функцию симуляции и парсить что пришло из формы
		log := login{
			Name:    c.PostForm("name"),
			Address: c.PostForm("address"),
		}
		c.JSON(200, gin.H{"status": log})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
