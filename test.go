package main

import (
	"fmt"
	"image-store-service/utils"
)

func main(){
	a := utils.GetAllAlbums()
	fmt.Println(a)
}