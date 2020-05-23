package main

import (
	"container/list"
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
	var inputPath = "G:\\2\\hiragana-sigoto2012-1.pdf"
	outputPath := "G:\\2\\"
	//allFiles := listFiles(inputPath, ".pdf")
	//if !folderExitOrnot(outputPath) {
	//	os.MkdirAll(outputPath, 777)
	//}
	//fmt.Println(allFiles.Back())
	//var index int
	//for inputPath := allFiles.Front(); inputPath != nil; inputPath = inputPath.Next() {
	//	index = index+1
	// rmWaterMark(inputPath.Value.(string), outputPath+"out\\", "-rmwaterMaker"+strconv.Itoa(index))
	rmWaterMark(inputPath, outputPath, "rmwaterMaker")
	//
	//}
	//splitPdf(inputPath, outputPath+"\\cai\\", 7)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	os.Exit(1)
	//}
	////
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

func rmWaterMark(inputPath, outputPath, outFileSuffix string) error {
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
	outFull := outputPath + fileName + outFileSuffix + fileExt
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
	fmt.Println("====ORG==================================")
	fmt.Println(string(fla))
	fmt.Println("====END ORG============ST======================")
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
				//fmt.Println(typex.Elem().Name(), typex.Elem().Kind())
				//fmt.Println(&x)
				//fmt.Println(x)
				//if validateTm(x, y) {
				//	(*parsed)[id] = &contentstream.ContentStreamOperation{}
				//	needreset = true
				//	startrm = true
				//	printRomved(i2)
				//}
				if typex.Elem() == reflect.TypeOf(core.PdfObjectFloat(0)) && (typey.Elem() == reflect.TypeOf(core.PdfObjectFloat(0))) {

					if validatekaKa(x.(*core.PdfObjectFloat), y.(*core.PdfObjectFloat)) {
						(*parsed)[id] = &contentstream.ContentStreamOperation{}
						needreset = true
						startrm = true

						printRomved(i2, " Tm ")

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

						startrm = true
						printRomved(i2, " CM ")
						needreset = false
					}

				}

			} else if i2.Operand == "l" || i2.Operand == "c" {
				if startrm {
					(*parsed)[id] = &contentstream.ContentStreamOperation{}

					needreset = true
					printRomved(i2, " l OR c ")
					startrm = true
				}

				//fmt.Println(strconv.Itoa(i) )
			} else if i2.Operand == "h" || i2.Operand == "re" || i2.Operand == "TJ" || i2.Operand == "Tj" {
				if startrm {
					(*parsed)[id] = &contentstream.ContentStreamOperation{}
					printRomved(i2, "TJ OR RE OR Tj OR h")
					//startrm=false
				}
			} else if i2.Operand == "m" {

				position := len(i2.Params) - 1
				x := i2.Params[position]
				y := i2.Params[position-1]
				typex := reflect.TypeOf(x)
				typey := reflect.TypeOf(y)
				if typex.Elem() == reflect.TypeOf(core.PdfObjectFloat(0)) && (typey.Elem() == reflect.TypeOf(core.PdfObjectFloat(0))) {

					if *x.(*core.PdfObjectFloat) < 16 && *x.(*core.PdfObjectFloat) > -16 && *y.(*core.PdfObjectFloat) > 365 && *y.(*core.PdfObjectFloat) < 370 {
						(*parsed)[id] = &contentstream.ContentStreamOperation{}
						needreset = true
						startrm = true
						printRomved(i2, "m")
					}
				}
			}
		}

	}
	fmt.Println(alltype)
	if needreset {
		//fmt.Println(string(parsed.Bytes()))
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
func validateTm(y, x core.PdfObject) bool {
	var vy, vx float64

	typex := reflect.TypeOf(x)
	typey := reflect.TypeOf(y)
	fmt.Println(typex.Elem().Name())
	fmt.Println(typex.Elem().Kind())
	if typex.Elem() == reflect.TypeOf(core.PdfObjectFloat(0)) {
		vx = reflect.ValueOf(x).Elem().Float()
	}
	if typex.Elem() == reflect.TypeOf(core.PdfObjectInteger(0)) {
		vx = float64(reflect.ValueOf(x).Elem().Int())
	}
	if typey.Elem() == reflect.TypeOf(core.PdfObjectFloat(0)) {
		vy = reflect.ValueOf(y).Elem().Float()
	}
	if typey.Elem() == reflect.TypeOf(core.PdfObjectInteger(0)) {
		vy = float64(reflect.ValueOf(y).Elem().Int())
	}

	if vx > 225 && vy < 22 {
		return true
	}

	fmt.Println(vx, vy)

	return false
}

func printRomved(rmObj *contentstream.ContentStreamOperation, tag string) {
	fmt.Println("=====RM==BY==" + tag + "======")
	fmt.Println(rmObj)

}

func validatekother(y, x *core.PdfObjectFloat) bool {
	//(*y>37 && *y<366 && *x>568 && *x<573.8) 处理竖条

	// && *x!=569.791 && *x!=570.1025) //处理hiraganahyo-a4-1 的单杠
	if ((*y < 378.5 && *y > 370) || *y < 18) && *x > 568 && *x < 577 || (*y > 37 && *y < 366 && *x > 568 && *x < 573.8 && *x != 568.5332 && *x != 569.791 && *x != 570.1025) || (*y < 376 && *x > 570 && *x < 578.68 && *x != 575.6201 && *x != 575.6797 && *x != 575.3105) {
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

func listFiles(inputFolder, fileExt string) (all *list.List) {
	all = list.New()
	allFiles, _ := ioutil.ReadDir(inputFolder)
	for _, file := range allFiles {
		if !file.IsDir() {
			fileFullPath := inputFolder + file.Name()
			if fileExt == filepath.Ext(fileFullPath) {
				all.PushFront(fileFullPath)
			}
			println(fileFullPath)
		}
	}
	return
}

func mergePdf(inputFolder string, outputPath, fileExt string) error {

	var inputListFiles = listFiles(inputFolder, fileExt)

	pdfWriter := pdf.NewPdfWriter()

	for inputPath := inputListFiles.Front(); inputPath != nil; inputPath = inputPath.Next() {
		f, err := os.Open(inputPath.Value.(string))
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
