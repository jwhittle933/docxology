package docxology

import (
	"fmt"
	"testing"
)

func TestExtractLocalDocx(t *testing.T) {
	fmt.Println("Running test TestExtractLocalDocx...")
	pathToFile := "./util/testfiles/hebrew_test.docx"
	fmt.Printf("Locating file at %s", pathToFile)
	unzip, err := ExtractLocalDocx(pathToFile)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("No errors.")
	if unzip == nil {
		t.Fatal("The file could not be found.")
	}
	if len(unzip.Files) < 1 {
		t.Fatalf("Was expecting multiple files, found %d", len(unzip.Files))
	}
}

func TestExtractFileHTTP(t *testing.T) {
	fmt.Println("Running test TestExtractFileHTTP...")
	//
}

func TestMapFiles(t *testing.T) {
	fmt.Println("Running test TestMapFiles...")

	//
}

func TestCopyToOS(t *testing.T) {
	fmt.Println("Running test TestCopyToOS...")
	//
}

func TestFindDoc(t *testing.T) {
	fmt.Println("Running test TestFindDoc...")
	//
}

func TestXMLExtractText(t *testing.T) {
	fmt.Println("Running test TestXMLExtractText...")
	//
}
