package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
	// "log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"image-store-service/utils"
	// "time"
	// "os"
	// "image/png"
)


type albumData struct{
	AlbumName string 
}
func main() {
  r := gin.Default()
  r.Static("/albums", "./albums")
  
// ALBUM API -----------------------------------------

//get All Album
r.GET("/album", getAllAlbumData)

//create album
r.POST("/album/create", createAlbum)

//Delete Album
r.DELETE("/album/:albumName", deleteAlbum)

// //Rename Album
r.PUT("/album/:albumName", renameAlbum)


// IMAGE APIS -------------------------------------------
// get single Image
 r.GET("/image/album/:albumName/image/:imageName", getImage)

// get All Image in album
  r.GET("/image/album/:albumName/all", getAllImages)

//   createImage api 
  r.POST("/image/album/:albumName", utils.CreateImageMain)

//   Delete Image API
  r.DELETE("/image/album/:albumName/image/:imageName", deleteImage)

// //   Delete ALL Image API
  r.DELETE("/image/album/:albumName/delete", deleteAllImages)

// //   upload Image API
  r.POST("/image/album/:albumName/upload", uploadImage)




  r.Run() // listen and serve on 0.0.0.0:8080
}

func getAllAlbumData(c *gin.Context){
	albums := utils.GetAllAlbums()
	c.JSON(http.StatusOK, gin.H{
		"data": albums,
	  })
}

func createAlbum(c *gin.Context){
	newData, err := ioutil.ReadAll(c.Request.Body)
	var jsonData albumData
	json.Unmarshal(newData, &jsonData)
	if err != nil {
		fmt.Println("Error")	
	}
	givenAlbumName := jsonData.AlbumName 
	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(givenAlbumName,albums)
 if exists {
	c.JSON(http.StatusNotAcceptable , gin.H{
		"message": "album with name "+ givenAlbumName +"already exists",
	  })
 }
	status := utils.CreateNewAlbum(givenAlbumName)
	if status {
		c.JSON(http.StatusOK, gin.H{
			"message": "album "+givenAlbumName +" created successfully",
		  })
	}
}
// //delete and force delete
func deleteAlbum(c *gin.Context){
	albumName := c.Param("albumName");

	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists or already deleted",
		 })
	} else { 
		message, statusCode := utils.DeleteAlbum(albumName)
	   c.JSON(statusCode, gin.H{
		"message": message,
	  })
	}
	
}

func renameAlbum(c *gin.Context){
	oldAlbumName := c.Param("albumName");
	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(oldAlbumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ oldAlbumName +"does not exists",
		 })
		return
	}
	newData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error")	
	}

	var jsonData albumData
	json.Unmarshal(newData, &jsonData)
	newAlbumName := jsonData.AlbumName
	message, status := utils.RenameAlbum(oldAlbumName, newAlbumName)
	c.JSON(status , gin.H{
		"message": message,
	  })
}


//Images
func getImage(c *gin.Context){ 
	albumName := c.Param("albumName");

	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists",
		 })
		return
	}

	imageName := c.Param("imageName");

	c.Redirect(http.StatusFound, "/albums/"+albumName+"/"+imageName)
}


// // album name
func getAllImages(c *gin.Context){
	albumName := c.Param("albumName");

	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists",
		 })
		return
	}

	allImages := utils.GetAllImages(albumName)
	c.JSON(http.StatusOK, gin.H{
		"data": allImages,
	  })

}


// func createImage(c *gin.Context){
// 	albumName := c.Param("albumName");
// 	albums := utils.GetAllAlbums()
// 	exists := utils.CheckDataInArray(albumName,albums)
// 	if !exists {
// 	   c.JSON(http.StatusNotAcceptable , gin.H{
// 		   "message": "album with name "+ albumName +"does not exists",
// 		 })
// 		return
// 	}

// 	newData, err := ioutil.ReadAll(c.Request.Body)
// 	var data jsonData
// 	json.Unmarshal(newData, &data)
// 	if err != nil {
// 		fmt.Println("Error")	
// 	}
// 	var colour string

// 	if data.Colour == "" {
// 		data.Colour = "#34eb52"
// 	}

// 	if data.size == nil {
// 		data.size = 400
// 	}
// }

func deleteImage(c *gin.Context){
	albumName := c.Param("albumName");
    
	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists",
		 })
		return
	}

	imageName := c.Param("imageName");
	images := utils.GetAllImages(albumName)
	exists2 := utils.CheckDataInArray(imageName,images)
	if !exists2 {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": "image "+ imageName +" does not exists in "+albumName,
		  })
		 return
	 }

	 message, statusCode  := utils.DeleteImage(albumName,imageName)
	   c.JSON(statusCode, gin.H{
		"message": message,
	  })
}

func deleteAllImages(c *gin.Context){
	albumName := c.Param("albumName");

	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists",
		 })
		return
	}
	message, statusCode  := utils.DeleteAllImage(albumName)
	c.JSON(statusCode, gin.H{
	 "message": message,
   })
}

func uploadImage(c *gin.Context) {
	albumName := c.Param("albumName");

	albums := utils.GetAllAlbums()
	exists := utils.CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists",
		 })
		return
	}

	file, _ := c.FormFile("file")
		// Upload the file to specific dst.
		c.SaveUploadedFile(file, "./albums/"+albumName+"/"+file.Filename)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

