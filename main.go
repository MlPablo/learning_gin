package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Stor Storage = NewStorage()

type album struct {
	ID     string  `json:"id" binding:"required"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

type HttpError struct {
	Error string `json:"error"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Hypno", Artist: "Travis Scott", Price: 120.99},
	{ID: "3", Title: "Rich da kid", Artist: "Asap Rocky", Price: 573.99},
}

func postAlbum(c *gin.Context) {
	var NewAlbum album
	if err := c.BindJSON(&NewAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpError{"bad_request"})
		return
	}
	alb, err := Stor.Add(NewAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpError{"Bad Form"})
	}
	c.IndentedJSON(http.StatusCreated, alb)

}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")
	album, err := Stor.OneRecord(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")
	err := Stor.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.IndentedJSON(http.StatusNoContent, HttpError{"deleted"})
}

func updateAlbumById(c *gin.Context) {
	id := c.Param("id")

	alb, err := Stor.OneRecord(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.BindJSON(&alb)
	Stor.Update(id, alb)
	c.IndentedJSON(http.StatusOK, alb)
}

func getAlbums(c *gin.Context) {
	albums, err := Stor.Read()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func GetRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.PUT("/albums/:id", updateAlbumById)
	router.DELETE("/albums/:id", deleteAlbumByID)
	router.POST("/albums", postAlbum)
	return router
}

func main() {
	router := GetRouter()
	router.Run(":8080")
}
