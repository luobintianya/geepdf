package main

import (
	"errors"
	"fmt"
	pdf "github.com/unidoc/unipdf/v3/model"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var inputPath = "C:\\Users\\lixin\\Desktop\\PMP全真模拟试题一.pdf"

	outputPath := "C:\\Users\\lixin\\Desktop\\outputpdf\\"

	err := splitPdf(inputPath, outputPath, 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	mergePdf(outputPath,outputPath)
	fmt.Printf("Complete, see output file: %s\n", outputPath)
	println(inputPath)
}

func splitPdf(inputPath string, outputPath string, splitfiles int) error {

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReaderLazy(f)
	if err != nil {
		return err
	}
	fileExt:=filepath.Ext(f.Name())
	fileName := strings.TrimSuffix(filepath.Base(f.Name()), fileExt)
	println(fileName)
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	prefilePages := int(math.Ceil(float64(numPages)/ float64(splitfiles)))
	println(strconv.Itoa(numPages)+" "+strconv.Itoa(prefilePages))

	for i := 0; i < splitfiles; i++ {

		pdfWriter := pdf.NewPdfWriter()

		for y := i * prefilePages; y < numPages && y < (i+1)*prefilePages; y++ {

			pageNum := y+1
			println(pageNum)
			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}

		}
		outFile := outputPath + fileName + strconv.Itoa(i) + fileExt
		println(outFile)
		fWrite, err := os.Create(outFile)
		if err != nil {
			return err
		}

		err = pdfWriter.Write(fWrite)
		fWrite.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func mergePdf(inputFolder string, outputPath string) error {

	var inputPaths []string
	allFiles,err := ioutil.ReadDir(inputFolder)

	var fileExt string
	for _, file := range allFiles {
		if !file.IsDir() {
			fileFullPath := inputFolder + file.Name()
			fileExt = filepath.Ext(fileFullPath)
			inputPaths = append(inputPaths, fileFullPath)
			println(fileFullPath)
		}

	}
	pdfWriter := pdf.NewPdfWriter()

	for _, inputPath := range inputPaths {
		f, err := os.Open(inputPath)
		if err != nil {
			return err
		}

		defer f.Close()

		pdfReader, err := pdf.NewPdfReader(f)
		if err != nil {
			return err
		}

		isEncrypted, err := pdfReader.IsEncrypted()
		if err != nil {
			return err
		}

		if isEncrypted {
			auth, err := pdfReader.Decrypt([]byte(""))
			if err != nil {
				return err
			}
			if !auth {
				return errors.New("Cannot merge encrypted, password protected document")
			}
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
	}


	fWrite, err := os.Create(outputPath+"merged"+fileExt)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}