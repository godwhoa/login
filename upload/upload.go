package upload

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func IsImage(filename string) bool {
	log.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Image open err")
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println("Image read err")
	}

	filetype := http.DetectContentType(buff)
	return strings.Contains(filetype, "image")

}

//They'll get filename for storing in db &
//results which is neg or pos which will be passed on to client
type Result struct {
	Filename string
	Res      string
}

func Upload(r *http.Request) Result {
	neg := Result{Filename: "", Res: "neg"}
	err := r.ParseMultipartForm(100000)
	if err != nil {
		log.Printf("Multipart err: %s", err.Error())
		return neg
	}
	m := r.MultipartForm
	var filename string

	//get the *fileheaders
	files := m.File["file"]
	log.Println(files)
	for i, _ := range files {
		filename = "./public/pics/" + files[i].Filename
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Printf("Fileheader err: %s\n", err.Error())
			return neg
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create(filename)
		defer dst.Close()
		if err != nil {
			log.Printf("File creation err: %s\n", err.Error())
			return neg
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("Copy err: %s\n", err.Error())
			return neg
		}

		if IsImage(filename) {
			//Give client the ok sign
			return Result{Filename: filename, Res: "pos"}
		} else {
			//If not image then it fails form validation
			//Delete the file we copied for validation
			_ = os.Remove(filename)
			return neg
			break
		}

	}
	if len(files) == 0 {
		return Result{Filename: "", Res: "pos"}
	} else {
		return neg
	}
}
