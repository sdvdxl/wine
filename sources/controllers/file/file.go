package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(c *gin.Context) {
	req := c.Request
	req.ParseMultipartForm(32 << 20)

	file, handler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		log.Logger.Error(err)
		return
	}

	defer file.Close()
	//	fmt.Fprintf(w, "%v", handler.Header)
	fileName := "./test/" + handler.Filename
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Logger.Error(err)
		return
	}

	io.Copy(f, file)
	f.Close()

	f, err = os.Open(fileName)
	if err != nil {
		log.Logger.Error(err)
		return
	}

	bytes, _ := ioutil.ReadAll(f)
	md5Value := util.MD5(string(bytes))

	c.JSON(http.StatusOK, md5Value+" "+c.Request.FormValue("test"))

}
