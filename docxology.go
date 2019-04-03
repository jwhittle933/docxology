/*
 * Package for unzipping, reading, manipulating, and storing docx files in their expanded, xml format.
 */

package docxology

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const defaultSaveLocation = "~/"

// UnZip struct for handling extention methods on unziped files.
type UnZip struct {
	Reader *zip.Reader
	Files  []*zip.File
}

// UnZipedFile struct for extending *zip.File
type UnZipedFile struct {
	File *zip.File
}

// XMLData struct for Unmarshalling xml.
// TODO: Determine necessary fields
type XMLData struct {
	Text    string   `xml:"w:t"`
	Name    string   `xml:"FullName"`
	XMLName xml.Name `xml:"Person"`
}

// Document type for storing word/document.xml extracted from docx zip.
type Document struct {
	Doc *zip.File
}

// Callback func type
type Callback func(string) error

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
		fi := Document{file}
		path := filepath.Join(saveLocation, file.Name)
		dirPath := filepath.Dir(path)
		os.MkdirAll(dirPath, 0777)
		_, err := os.Create(path)
		if err != nil {
			return err
		}

		if err := fi.CopyToOS(path); err != nil {
			return err
		}
	}
	return nil
}

// CopyToOS for writing contents of d *Document to disc.
func (d *Document) CopyToOS(filePath string) error {
	fileReader, err := d.Doc.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, d.Doc.Mode())
	if err != nil {
		return err
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	return nil
}

//FindDoc locates file by filename and returns Document
func (f *UnZip) FindDoc(searchDoc string) (file Document) {
	for _, fi := range f.Files {
		if fi.Name == searchDoc {
			file = Document{
				Doc: fi,
			}
		}
	}
	return file
}

// XMLExtractText for manipulating xml
func (d *Document) XMLExtractText() {
	file, err := d.Doc.Open()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	s := XMLData{}
	if err := xml.Unmarshal(data, &s); err != nil {
		panic(err)
	}

	return
}
