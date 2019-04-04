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

// Document type for storing word/document.xml extracted from docx zip.
type Document struct {
	Doc *zip.File
}

// XMLDocMacroData struct for Unmarshalling xml.
/*
!! https://www.loc.gov/preservation/digital/formats/fdd/fdd000397.shtml
!! http://officeopenxml.com/anatomyofOOXML.php
!! https://docs.microsoft.com/en-us/office/open-xml/structure-of-a-wordprocessingml-document
*
* Notes from Microsoft:
* A WordprocessingML document is organized around the concept of stories.
* A story is a region of content in a WordprocessingML document.
*
* The main document story of the simplest WordprocessingML document consists of the following XML elements:
* document – The root element
* body – The container for the collection of block-level structures
* p – A paragraph
* r – A run
* t – A range of text
*/
type XMLDocMacroData struct {
	DocumentMeta xml.Name `xml:"document"`
	Text         string   `xml:"body>p>r>t"`
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
// This method must include dirs, i.e., word/document/xml, word/theme/theme1.xml
func (f *UnZip) FindDoc(searchDoc string) (file *Document) {
	for _, fi := range f.Files {
		if fi.Name == searchDoc {
			file = &Document{
				Doc: fi,
			}
		}
	}
	return file
}

// XMLExtractText for manipulating xml
func (d *Document) XMLExtractText() {
	var doc XMLDocMacroData
	file, err := d.Doc.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err := xml.Unmarshal(data, &doc); err != nil {
		panic(err)
	}

	return
}
