package helper

import (
	"errors"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"strings"
)

type (
	FileUpload struct {
		OriginName		string
		SaveName		string			`json:"-"`
		SavePath		string			`json:"-"`
		PublicPath		string
	}
)

func UploadFile(r *http.Request,key string,storePath string,subs ...string) (info FileUpload,err error) {
	file, handler, err := r.FormFile(key)
	if err != nil {
		err = errors.New("Wrong input name")
		return
	}
	defer file.Close()


	split := strings.Split(handler.Filename,".")
	sub := split[len(split) - 1]
	if len(subs) > 0 && IndexSliceString(subs, sub) == -1 {
		err = errors.New("File type is not allow")
		return
	}
	info.OriginName = handler.Filename
	info.SaveName = uuid.NewV1().String() + "." + sub
	info.SavePath = "./public/upload/" + storePath + info.SaveName
	info.PublicPath = "/public/upload/" + storePath + info.SaveName

	if ok,_ := exists("./public/upload/" + storePath); !ok {
		if err = os.MkdirAll("./public/upload/" + storePath, os.ModePerm);err != nil {
			return
		}
	}

	f, err := os.OpenFile(info.SavePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		//err = errors.New("Wrong path")
		return
	}

	defer f.Close()
	io.Copy(f, file)
	return
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}
func (f FileUpload) Remove() error {
	return os.Remove(f.SavePath)
}

func IndexSliceString(slice []string,str string) int {
	for k,v := range slice{
		if str == v {
			return k
		}
	}
	return -1
}

func IndexSliceInt(slice []int,str int) int {
	for k,v := range slice{
		if str == v {
			return k
		}
	}
	return -1
}
