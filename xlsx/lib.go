package xlsx

import(
	"os"
		ex "github.com/xuri/excelize/v2"
		"runtime/pprof"
		xl "github.com/tealeg/xlsx/v3"
		"fmt"
)

type Excel struct{
	file *ex.File
}

func OpenFile(filepath string)(*Excel,error){
	file,err := ex.OpenFile(filepath)
	if err != nil{
		return nil,err
	}
	return &Excel{
		file:file,
	},nil
}

func (this *Excel) GetAllRows(sheet string)([][]string,error){
	return this.file.GetRows(sheet)
}

func WriteToFile(filename string,sheetname string, values [][]string)error{
	file := ex.NewFile()
	file.NewSheet(sheetname)
	for i,v := range values{
		err :=file.SetSheetRow(sheetname,fmt.Sprintf("A%d",i+1),&v)
		if err != nil{
			return err
		}
	}
	return file.SaveAs(filename)
}

func WriteToFileWithSlices(filename string,sheetname string, values [][]interface{})error{
	file := ex.NewFile()
	file.NewSheet(sheetname)
	for i,v := range values{
		err :=file.SetSheetRow(sheetname,fmt.Sprintf("A%d",i+1),&v)
		if err != nil{
			return err
		}
	}
	return file.SaveAs(filename)
}

func StreamWriteToFile(filename string,sheetname string, values [][]interface{})error{
	file := ex.NewFile()
	sheetIndex := file.NewSheet(sheetname)
	file.SetActiveSheet(sheetIndex)
	writer, err := file.NewStreamWriter(sheetname)
	if err != nil{
		return err
	}
	for i,v := range values{
		//if i+1 % 1000 == 0 {
		//	err = writer.Flush()
		//	if err != nil {
		//		return err
		//	}
		//}
		err :=writer.SetRow(fmt.Sprintf("A%d",i+1),v)
		if err != nil{
			return err
		}
	}
	err = writer.Flush()
	if err != nil{
		return fmt.Errorf("Flush err when  %w",err)
	}
	return file.SaveAs(filename)
}

func WriteToFile_Tmp(filename string,sheetname string,values [][]string)error{
	fmt.Println("Call  WriteToFile_Tmp")
	file:= xl.NewFile(xl.UseDiskVCellStore)
	// file:= xl.NewFile()
	sh,err := file.AddSheet(sheetname)
	if err != nil{
		return err
	}

	for i,rowValues  := range values{
		fmt.Println(i," Add Row")
		row:= sh.AddRow()
		for cellValue := range rowValues{
			fmt.Println("Add Cell")
			cell := row.AddCell()
			cell.SetValue(cellValue)
		}
	}
	fmt.Println("Debug")
	debugFile,err := os.Create("write_to_file")
	if err != nil{
		panic(err)
	}
	defer debugFile.Close()
	err = pprof.WriteHeapProfile(debugFile)
	if err != nil{
		panic(err)
	}
	return file.Save(filename)
}
