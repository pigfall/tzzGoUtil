package xlsx

import(
		ex "github.com/xuri/excelize/v2"
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
		err := file.SetSheetRow(sheetname,fmt.Sprintf("A%d",i+1),&v)
		if err != nil{
			return err
		}
	}
	return file.SaveAs(filename)
}
