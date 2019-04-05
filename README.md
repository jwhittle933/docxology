# Docxology

Golang Word Doc (.docx) file extractor and manipulator.

_In progress_ and open to contributions, suggestions, etc.

[![Build Status](https://travis-ci.com/jwhittle933/docxology.svg?branch=master)](https://travis-ci.com/jwhittle933/docxology)
[![GoDoc](https://godoc.org/github.com/jwhittle933/docxology?status.svg)](https://godoc.org/github.com/jwhittle933/docxology)

## How To

### Info

.docx files are really just "application/zip" made of XML files. This package is intended to assist in extracting the XML and manipulating the data as you need. Go has everything you need built-in to handle this type of functionality, so this package aims, not to replace that funcationality, but to make it more immediately user-friendly.

### Get

```bash
go get github.com/jwhittle933/docxology
```
