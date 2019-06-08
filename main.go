package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cpllbstr/gogrn/grn"

	"github.com/cpllbstr/gogrn/slv"
	"github.com/gin-gonic/gin"
)

type login struct {
	Name    string
	Address string
}

func main() {
	filsys := make(map[string]bool)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("./www/index.html")
	})
	r.GET("/sim", func(c *gin.Context) {
		c.File("./www/sim.html")
	})
	r.GET("/dat", func(c *gin.Context) {
		t, err := template.ParseFiles("./www/dat.html")
		if err != nil {
			panic(err)
		}

		err = t.Execute(c.Writer, filsys)
		if err != nil {
			panic(err)
		}
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

		c.Redirect(301, "/dat")
		avg := c.Request.FormValue("avg")

		body := slv.NewThreeBodyModel(m, k)
		//y, mon, d := time.Now().Date()
		filename := fmt.Sprintf("sym%v.dat", time.Now().Format("15:04:05_02Jan06"))
		file, err := os.Create(fmt.Sprint("./dat/", filename)) //файл в который пишутся результаты расчетов
		if err != nil {
			log.Println("Cannot create file!  :", err)
			c.File("./www/error.html")
			return
		}

		machine := slv.StateMachFromModel(body, v, l, step, 0, fint)
		if avg == "on" {
			go func(machine grn.StateMachine, body grn.ThreeBodyModel, file *os.File, filename string) {
				filsys[filename] = false
				slv.SimulateAv(machine, body, file)
				filsys[filename] = true
			}(machine, body, file, filename)

		} else {
			go func(machine grn.StateMachine, file *os.File, filename string) {
				filsys[filename] = false
				time.Sleep(30 * time.Second)
				slv.Simulate(machine, file)
				filsys[filename] = true
			}(machine, file, filename)
		}
	})
	r.POST("/2d", func(c *gin.Context) {
		c.File("./www/error.html")
	})
	r.POST("/3d", func(c *gin.Context) {
		c.File("./www/error.html")
	})
	r.POST("/dat", func(c *gin.Context) {
		filename := c.Request.FormValue("file")
		c.FileAttachment(fmt.Sprintf("./dat/%v", filename), filename)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
