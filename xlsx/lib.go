package xlsx

import(
		ex "github.com/xuri/excelize/v2"
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
