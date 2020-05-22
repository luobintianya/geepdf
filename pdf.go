package main

import (
	"errors"
	"fmt"
	"github.com/unidoc/unipdf/v3/contentstream"
	"github.com/unidoc/unipdf/v3/core"
	pdf "github.com/unidoc/unipdf/v3/model"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	var inputPath = "G:\\2\\iro 24 - 副本.pdf"
	outputPath := "G:\\2\\"
	rmWaterMark(inputPath, outputPath)
	if !folderExitOrnot(outputPath) {
		os.MkdirAll(outputPath, 777)
	}
	//err := splitPdf(inputPath, outputPath, 12)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//mergePdf(outputPath,outputPath)
	//fmt.Printf("Complete, see output file: %s\n", outputPath)
	println(inputPath)
}

func folderExitOrnot(fielPath string) bool {
	_, err := os.Stat(fielPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func rmWaterMark(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()
	//common.SetLogger(common.ConsoleLogger{LogLevel: 5})
	pdfWriter := pdf.NewPdfWriter()
	pdfReader, err := pdf.NewPdfReaderLazy(f)
	if err != nil {
		return err
	}
	fileExt := filepath.Ext(f.Name())
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
	//pdfWriter := pdf.NewPdfWriter()

	for i := 0; i < numPages; i++ {
		pdfPage, _ := pdfReader.GetPage(i + 1)

		fmt.Println("Start prcess page " + strconv.Itoa(i+1))
		if pdfPage.Contents != nil {
			typeOf := reflect.TypeOf(pdfPage.Contents)
			typeOfEl := typeOf.Elem()
			if typeOfEl == reflect.TypeOf(core.PdfObjectArray{}) {
				arrays := pdfPage.Contents.(*core.PdfObjectArray)

				if arrays.Len() > 0 {

					//	fterx := arrays.Get(0).(*core.PdfObjectReference).Resolve().(core.PdfObject).(*core.PdfObjectStream)
					for _, norContent := range arrays.Elements() {
						fterValue := norContent.(*core.PdfObjectReference).Resolve().(core.PdfObject).(*core.PdfObjectStream)
						filterTj(fterValue, false)
					}
				}
			}
			if typeOfEl == reflect.TypeOf(core.PdfObjectReference{}) {
				filter := pdfPage.Contents.(*core.PdfObjectReference).Resolve().(*core.PdfObjectStream)

				filterTj(filter, false)

			}

		}
		pdfWriter.AddPage(pdfPage)
	}
	outFull := outputPath + "rmwaterMarker" + fileExt
	fmt.Println(outFull)
	fWrite, err := os.Create(outFull)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)

	return nil
}

func filterTj(content *core.PdfObjectStream, left bool) {

	fla, _ := core.DecodeStream(content)
	fmt.Println(string(fla))
	fmt.Println("proc==================================")
	cStreamParser := contentstream.NewContentStreamParser(string(fla))
	parsed, _ := cStreamParser.Parse()
	needreset := false

	startrm := false
	alltype := ""
	for id, i2 := range *parsed {
		if !strings.Contains(alltype, i2.Operand) {
			alltype = alltype + " " + i2.Operand
		}
		//cCKkSQmcmlwreTfjBDCrigsTjMJdDoG
		if i2.Params != nil && !(strings.Contains("BDC k ri gs f EMC q W n  Q BT Tf ET K w d S G Do Td", i2.Operand)) {
			//
			//if i2.Operand == "TJ" {
			//	pfoa := i2.Params[0].(*core.PdfObjectArray)
			//	for _, ix := range pfoa.Elements() {
			//		switch vv := ix.(type) {
			//		case *core.PdfObjectString:
			//			decoded := vv.Decoded()
			//			if strings.ContainsAny(decoded, "FãFþGe") || strings.ContainsAny(decoded, "		-	") {
			//				pfoa.Clear()
			//				needreset = true
			//			}
			//
			//		}
			//	}
			//
			//} else

			if i2.Operand == "Tm" {
				startrm = false
				position := len(i2.Params) - 1
				x := i2.Params[position]
				y := i2.Params[position-1]
				typex := reflect.TypeOf(x)
				typey := reflect.TypeOf(y)
				if typex.Elem() == reflect.TypeOf(core.PdfObjectFloat(0)) && (typey.Elem() == reflect.TypeOf(core.PdfObjectFloat(0))) {

					if (*x.(*core.PdfObjectFloat) < 16 && *y.(*core.PdfObjectFloat) > 250 && *y.(*core.PdfObjectFloat) < 608) || validatekaKa(x.(*core.PdfObjectFloat), y.(*core.PdfObjectFloat)) {
						(*parsed)[id] = &contentstream.ContentStreamOperation{}
						needreset = true
						startrm = true
						fmt.Println("=============")
						fmt.Println(i2)
					}
				}

			} else if i2.Operand == "cm" {
				startrm = false
				position := len(i2.Params) - 1
				x := i2.Params[position]
				y := i2.Params[position-1]
				typex := reflect.TypeOf(x)
				typey := reflect.TypeOf(y)
				if typex.Elem() == reflect.TypeOf(core.PdfObjectFloat(0)) && (typey.Elem() == reflect.TypeOf(core.PdfObjectFloat(0))) {

					if validatekother(x.(*core.PdfObjectFloat), y.(*core.PdfObjectFloat)) {
						(*parsed)[id] = &contentstream.ContentStreamOperation{}
						needreset = true
						startrm = true
						fmt.Println("=============")
						fmt.Println(i2)
					}

				}

			} else if i2.Operand == "l" || i2.Operand == "c" {
				if startrm {
					(*parsed)[id] = &contentstream.ContentStreamOperation{}
					startrm = true
					needreset = true
					fmt.Println("=============")
					fmt.Println(i2)
				}

				//fmt.Println(strconv.Itoa(i) )
			} else if i2.Operand == "h" || i2.Operand == "re" || i2.Operand == "TJ" || i2.Operand == "Tj" {
				if startrm {
					(*parsed)[id] = &contentstream.ContentStreamOperation{}
					fmt.Println("=============")
					fmt.Println(i2)
				}
			} else if i2.Operand == "m" {
				startrm = false
				position := len(i2.Params) - 1
				x := i2.Params[position]
				y := i2.Params[position-1]
				typex := reflect.TypeOf(x)
				typey := reflect.TypeOf(y)
				if typex.Elem() == reflect.TypeOf(core.PdfObjectInteger(0)) && (typey.Elem() == reflect.TypeOf(core.PdfObjectInteger(0))) {
					if *x.(*core.PdfObjectInteger) == 0 && *y.(*core.PdfObjectInteger) == 0 {
						(*parsed)[id] = &contentstream.ContentStreamOperation{}
						needreset = true
						startrm = true
						fmt.Println("=============")
						fmt.Println(i2)
					}
				}
			}
		}

	}
	fmt.Println(alltype)
	if needreset {
		fmt.Println(string(parsed.Bytes()))
		content.Stream, _ = core.NewFlateEncoder().EncodeBytes(parsed.Bytes())
	}

	//fmt.Println(xk)
}

func validatekaKa(y, x *core.PdfObjectFloat) bool {
	if *y < 23.5 && *x > 12 { //x坐标小于12 y坐标大于30    6 0 0 6 24.5576019(x) 30.1304016(y) Tm这样的就不处理
		return true
	}
	return false
}
func validatekother(y, x *core.PdfObjectFloat) bool {
	if ((*y < 378.5 && *y > 370) || *y < 18) && *x > 568 && *x < 608 {
		return true
	}
	return false
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
	fileExt := filepath.Ext(f.Name())
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

	prefilePages := int(math.Ceil(float64(numPages) / float64(splitfiles)))
	println(strconv.Itoa(numPages) + " " + strconv.Itoa(prefilePages))

	for i := 0; i < splitfiles; i++ {

		pdfWriter := pdf.NewPdfWriter()

		for y := i * prefilePages; y < numPages && y < (i+1)*prefilePages; y++ {

			pageNum := y + 1
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
	allFiles, err := ioutil.ReadDir(inputFolder)

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

	fWrite, err := os.Create(outputPath + "merged" + fileExt)
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
