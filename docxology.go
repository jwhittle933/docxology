package docxology

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// UnZip struct for handling extention methods on unziped files.
type UnZip struct {
	Reader *zip.Reader
	Files  []*zip.File
}

// XMLData struct for Unmarshalling xml.
type XMLData struct {
	Text string `xml:"w:t"`
}

// Document type for storing word/document.xml extracted from docx zip.
type Document struct {
	Doc *zip.File
}

// DocxOnDiscUnzip for reading Word .docx files << Entry func
// pathToFile param path to file stored on disc.
func DocxOnDiscUnzip(pathToFile string) error {
	saveLocation := fmt.Sprintf("./saves/dir-%d/unzip", uuid.New())
	zip := ExtractLocalFile(pathToFile)

	if err := zip.MapFiles(saveLocation); err != nil {
		return err
	}

	return nil
}

// ExtractLocalFile returns *UnZip
// pathToFile param path to file stored on disc.
func ExtractLocalFile(pathToFile string) *UnZip {
	file, err := os.Open(pathToFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))

	return &UnZip{
		Reader: zipReader,
		Files:  zipReader.File,
	}
}

// ExtractFileHTTP return *UnZip
func ExtractFileHTTP(fi *multipart.FileHeader) *UnZip {
	file, err := fi.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(src), fi.Size)

	return &UnZip{
		Reader: zipReader,
		Files:  zipReader.File,
	}
}

// MapFiles for iterating through zip.File slice
// and performing an operation on it.
func (f *UnZip) MapFiles(saveLocation string) error {
	if err := os.MkdirAll(saveLocation, 0755); err != nil {
		return err
	}
	for _, file := range f.Files {
		path := filepath.Join(saveLocation, file.Name)
		dirPath := filepath.Dir(path)
		os.MkdirAll(dirPath, 0777)
		_, err := os.Create(path)
		if err != nil {
			return err
		}

		if err := CopyToOS(file, path); err != nil {
			return err
		}
	}
	return nil
}

// CopyToOS for mapping over files.
func CopyToOS(file *zip.File, filePath string) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	// t := XMLData{}
	// data, err := ioutil.ReadAll(targetFile)
	// if err != nil {
	// 	return err
	// }
	// xml.Unmarshal([]byte(data), &t)
	// fmt.Println("Reading XML", t.Text)

	return nil
}

//FindDoc locates file named word/document.xml
func (f *UnZip) FindDoc() (file *Document) {
	for _, fi := range f.Files {
		if fi.Name == "word/document.xml" {
			file = &Document{
				Doc: fi,
			}
		}
	}
	return file
}

// XMLExtractText for manipulating xml
func XMLExtractText() {
	return
}
