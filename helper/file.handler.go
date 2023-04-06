package helper

import (
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/KuroNeko6666/prima-api/interfaces/form"
	"github.com/KuroNeko6666/prima-api/storage"
	"github.com/minio/minio-go"
)

func FileUpload(file *multipart.FileHeader, bucket string, user model.User) (result form.FileForm, e error) {

	client, err := storage.ConnectStorage()

	if err != nil {
		fmt.Println("connect")
		return result, err
	}

	buffer, err := file.Open()

	if err != nil {
		fmt.Println("open")

		return result, err
	}

	defer buffer.Close()
	currentTime := time.Now()
	replacer := strings.NewReplacer(":", "-", ".", "-", " ", "T")
	rawName := strings.Split(file.Filename, ".")
	ext := rawName[len(rawName)-1]
	name := rawName[0]
	fileTime := replacer.Replace(currentTime.Format("2006.01.02 15:04:05"))

	var res = form.FileForm{
		FileName:   name + "CR-" + user.ID + "-" + fileTime + "." + ext,
		FileBuffer: buffer,
		FileSize:   file.Size,
		FileType:   file.Header["Content-Type"][0],
	}

	info, err := client.PutObject(bucket, res.FileName, res.FileBuffer, res.FileSize, minio.PutObjectOptions{ContentType: res.FileType})
	if err != nil {
		fmt.Println("put")

		return result, err
	}

	log.Println(info)

	return res, nil

}

func FilesUpload(files *multipart.Form, bucket string, user model.User) ([]form.FileForm, error) {
	var result []form.FileForm
	client, err := storage.ConnectStorage()

	if err != nil {
		return result, err
	}

	for _, file := range files.File {

		for _, fileHeader := range file {
			buffer, err := fileHeader.Open()

			if err != nil {
				return result, err
			}

			defer buffer.Close()
			replacer := strings.NewReplacer(":", "-", ".", "-", "+", "-")
			rawName := strings.Split(fileHeader.Filename, ".")
			ext := rawName[len(rawName)-1]
			name := rawName[0]
			fileTime := replacer.Replace(time.Now().Format(time.RFC3339))

			var res = form.FileForm{
				FileName:   name + "-CR-" + user.ID + "-" + fileTime + "." + ext,
				FileBuffer: buffer,
				FileSize:   fileHeader.Size,
				FileType:   fileHeader.Header["Content-Type"][0],
			}

			info, err := client.PutObject(bucket, res.FileName, res.FileBuffer, res.FileSize, minio.PutObjectOptions{ContentType: res.FileType})
			if err != nil {
				return result, err
			}

			result = append(result, res)

			log.Println(info)
		}
	}

	return result, nil

}

func FileDelete(fileName string, bucket string) error {
	client, err := storage.ConnectStorage()

	if err != nil {
		return err
	}

	if err := client.RemoveObject(bucket, fileName); err != nil {
		return err
	}

	return nil
}
