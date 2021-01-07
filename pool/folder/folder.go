package folder

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type FileInfo struct {
	Name	string	`json:"Name"`
	Size	int64	`json:"Size"`
	ModTime	time.Time	`json:"ModTime"`
	IsDir	 bool	`json:"IsDir"`
}

func GetFolderlist(folder string) []byte {
	var fileInfos []FileInfo
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		fileInfo := FileInfo{
			Name: file.Name(),
			Size:file.Size(),
			ModTime:file.ModTime(),
			IsDir:file.IsDir(),
		}
		fileInfos = append(fileInfos, fileInfo)
	}
	bytesData, _ := json.Marshal(fileInfos)

	return bytesData
}