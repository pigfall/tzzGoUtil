package parser


import(
		"go/parser"
		"go/token"
)

func ParseFile(fileSet *token.FileSet,filename string,mode parser.Mode)(*ASTFile,error){
	astFile,err := parser.ParseFile(fileSet,filename,nil,mode)
	if err != nil{
		return nil,err
	}
	return ASTFileNew(astFile),nil
}


