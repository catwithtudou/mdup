package api

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"io"
	"io/ioutil"
	"log"
	"mdup/config"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
)

/**
 * user: ZY
 * Date: 2020/4/7 16:56
 */

const(
	StaticViewPath = "/static/view/"
	TempFilePath = "./static/temp/"
)


func FileUpload(c *gin.Context){
	c.Redirect(http.StatusFound,StaticViewPath+"index.html")
}

func DoFileUpload(c *gin.Context){
	file,header,err:=c.Request.FormFile("file")
	if err!=nil{
		log.Println(err)
		ParamError(c)
		return
	}
	defer file.Close()

	filePath,err:=fileHandle(file,header)


	c.Header("Content-Disposition", "attachment; filename="+header.Filename)
	c.Header("Content-Type","application/file")
	c.File(filePath)

	//文件返回完成后要删除该文件
	_=os.Remove(filePath)
}



func fileHandle(file multipart.File,header *multipart.FileHeader)(filePath string,err error){

	data,err:=ioutil.ReadAll(file)
	if err!=nil{
		log.Println(err)
		return
	}

	//文件内容
	content:=string(data)

	//根据文件内容匹配正则
	regImg:=`(.*?)\!\[(.*?)\]\((.*?)\)`
	compile:=regexp.MustCompile(regImg)
	imgRow:=compile.FindAllStringSubmatch(content,-1)

	imgRows:=make([]string,len(imgRow))
	for k,v:=range imgRow{
		imgRows[k]=v[3]
	}

	reImgRows,err:=fileOperation(imgRows)
	if err!=nil{
		log.Println(err)
		return
	}


	err=uploadPhotoBed(imgRows,reImgRows)
	if err!=nil{
		log.Println(err)
		return
	}

	//上传成功后将原文件值进行单个替换
	for k,v:=range imgRows{
		regImgChange:=regexp.MustCompile(strings.ReplaceAll(v,`\`,`\\`))
		content=regImgChange.ReplaceAllString(content,config.Domain+reImgRows[k])
	}

	//最后将文件内容复制给一个新文件，然后最终返回
	//fmt.Println(content)

	//创建临时文件
	tempFilePath:=TempFilePath+header.Filename
	newFile,err:=os.Create(tempFilePath)
	if err!=nil{
		log.Println(err)
		return
	}
	defer newFile.Close()

	_,err=newFile.WriteString(content)
	if err!=nil{
		log.Println(err)
	}

	filePath=tempFilePath

	return
}


func uploadPhotoBed(imgRows []string,reImgRows []string)(err error){
	putPolicy:=storage.PutPolicy{
		Scope:config.Bucket,
	}
	//TODO:凭证应该存入redis保存
	mac:=qbox.NewMac(config.AccessKey,config.SecretKey)
	upToken:=putPolicy.UploadToken(mac)

	cfg:=storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: true,
	}

	formUploader:=storage.NewFormUploader(&cfg)
	ret:=storage.PutRet{}

	putExtra:=storage.PutExtra{
		Params: map[string]string{
			"x:name":"mdup",
		},
	}


	for k,v:=range imgRows{
		err=formUploader.PutFile(context.Background(),&ret,upToken,reImgRows[k],v,&putExtra)
		if err!=nil{
			log.Println(err)
			return
		}
	}

	return
}

//fileOperation 判断数组内图片文件是否存在，若存在则值换成hash值，若有一个不存在则返回错误
func fileOperation(imgRows []string)(result []string,err error){
	result=make([]string,len(imgRows))

	for k,v:=range imgRows{
		//判断每个文件是否存在如果不存在则直接失败
		//这里为绝对路径
		file,err:=os.Open(v)
		if err!=nil{
			return nil,err
		}
		//若存在则计算其hash值
		result[k]=fileSha1(file)
		file.Close()
	}
	return
}

//fileSha1 计算FileHash值
func fileSha1(file *os.File) string{
	_sha1:=sha1.New()
	io.Copy(_sha1,file)
	return hex.EncodeToString(_sha1.Sum(nil))
}


