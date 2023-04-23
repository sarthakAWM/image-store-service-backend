package utils

import (
	"io/ioutil"
	"fmt"
    "image"
	"image/color"
	"image/draw"
	"image/png"
	"encoding/json"
	"log"
	"net/http"
	// "strconv"
	// "time"
	"os"
    "github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)


type jsonData struct {
	ImageName string
	Size int `default:400` 
	Initals string
	Colour string `default:"#34eb52"` 
  }

func GetAllImages(albumName string) []string {

	var images = make([]string, 0, 1)

	files, err := ioutil.ReadDir("./albums/"+albumName+"/")
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        // fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir(){
			images = append(images, file.Name())
		}
		
    }
	return images
}

func DeleteImage(albumName string, imageName string) (string, int){
	err := os.RemoveAll("./albums/"+albumName+"/"+imageName)
    if err != nil {
      return string("error in deleting image "+imageName+" from album"+albumName), 400
    }
	return "deleted image "+imageName+" from album"+albumName, 200
}

func DeleteAllImage(albumName string) (string, int){
	err := os.RemoveAll("./albums/"+albumName+"/")
    if err != nil {
      return string("error in deleting all images from album "+albumName), 400
    }
	 CreateNewAlbum(albumName)
	return "deleted all image from album "+albumName, 200
}


func CreateImageMain(c *gin.Context){

	albumName := c.Param("albumName");
	albums := GetAllAlbums()
	exists := CheckDataInArray(albumName,albums)
	if !exists {
	   c.JSON(http.StatusNotAcceptable , gin.H{
		   "message": "album with name "+ albumName +"does not exists",
		 })
		return
	}

	newData, err := ioutil.ReadAll(c.Request.Body)
	var data jsonData
	json.Unmarshal(newData, &data)
	if err != nil {
		fmt.Println("Error")	
	}

	if data.Colour == "" {
		data.Colour = "#34eb52"
	}

	if data.Size == 0 {
		data.Size = 400
	}

	imageName := data.ImageName
	images := GetAllImages(albumName)
	exists2 := CheckDataInArray(imageName,images)
	if exists2 {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": "image "+ imageName +" does already exists in "+albumName,
		  })
		 return
	 }

	avatar, err := CreateAvatar(data.Size, data.Initals, data.Colour)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("albums/"+albumName+"/"+imageName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": "error in creating file",
		  })
		 return
	}
	png.Encode(file, avatar)
	c.JSON(200 , gin.H{
		"message": "Image Successfully created",
	  })
}

func CreateAvatar(size int, initials string, colour string) (*image.RGBA, error) {
	width, height := size, size
	// bgColor, err := hexToRGBA("#764abc")
	bgColor, err := hexToRGBA(colour)
	if err != nil {
		
		log.Fatal(err)
	}

	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: bgColor},
		image.Point{}, draw.Src)

	drawText(background, initials)
	return background, err
}

func drawText(canvas *image.RGBA, text string) error {
	var (
		fgColor  image.Image
		fontFace *truetype.Font
		err      error

		fontSize = 128.0
	)
	fgColor = image.White
	fontFace, err = freetype.ParseFont(goregular.TTF)
	fontDrawer := &font.Drawer{
		Dst: canvas,
		Src: fgColor,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}
	textBounds, _ := fontDrawer.BoundString(text)
	xPosition := (fixed.I(canvas.Rect.Max.X) - fontDrawer.MeasureString(text)) / 2
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	yPosition := fixed.I((canvas.Rect.Max.Y)-textHeight.Ceil())/2 + fixed.I(textHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: xPosition,
		Y: yPosition,
	}
	fontDrawer.DrawString(text)
	return err
}

func hexToRGBA(hex string) (color.RGBA, error) {
	var (
		rgba color.RGBA
		err  error
	)
	rgba.A = 0xFF
	switch len(hex) {
	case 7:
		_, err = fmt.Sscanf(hex, "#%02x%02x%02x", &rgba.R, &rgba.G, &rgba.B)
	case 4:
		_, err = fmt.Sscanf(hex, "#%1x%1x%1x", &rgba.R, &rgba.G, &rgba.B)
		rgba.R *= 17
		rgba.G *= 17
		rgba.B *= 17
	default:
		err = fmt.Errorf("invalid hex code")
	}
	return rgba, err
}
