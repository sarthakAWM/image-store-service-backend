package utils

import(
	"io/ioutil"
	"log"
	"os"
)


func GetAllAlbums() []string {

	var albums = make([]string, 0, 1)

	files, err := ioutil.ReadDir("./albums")
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        // fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			albums = append(albums, file.Name())
		}
    }
	return albums

}

func CreateNewAlbum(albumName string) bool{
	if err := os.Mkdir("albums/"+albumName, os.ModePerm); err != nil {
        // log.Fatal(err)
		return false
    }
	return true
}

func CheckDataInArray(searchString string, newSlice []string) bool {
	found := false

	for _, v := range newSlice {
		if v == searchString {
			found = true
			break
		}
	}
	return found
}

func DeleteAlbum(albumName string) (string, int) {
	err := os.RemoveAll("./albums/"+albumName)
    if err != nil {
      return string("error in deleting album"), 400
    }
	return "deleted album with name "+ albumName, 200
}

func RenameAlbum(oldAlbumName string, newAlbumName string ) (string, int) {
	err := os.Rename("./albums/"+oldAlbumName, "./albums/"+newAlbumName)
    if err != nil {
		return string("error in renaming album"), 400
    }
	return "album with name "+ oldAlbumName +" renamed as "+newAlbumName, 200
}