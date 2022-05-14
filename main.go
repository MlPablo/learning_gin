package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type album struct {
	ID     string  `json:"id" binding:"required"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

type error struct {
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
		c.IndentedJSON(http.StatusBadRequest, error{"bad_request"})
		return
	}

	albums = append(albums, NewAlbum)
	c.IndentedJSON(http.StatusCreated, NewAlbum)

}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, error{"not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range albums {
		if a.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusNoContent, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, error{"not found"})
}

func updateAlbumById(c *gin.Context) {
	id := c.Param("id")

	for i, a := range albums {
		if a.ID == id {
			c.BindJSON(&a)
			albums[i] = a
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, error{"not found"})

}

func getAlbums(c *gin.Context) {
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
