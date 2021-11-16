package parser

import(
		"go/ast"
)


type ASTFile struct{
	*ast.File
}

func ASTFileNew(astFile *ast.File)*ASTFile{
	return &ASTFile{
		File:astFile,
	}
}

func (this *ASTFile) Inspect(inspect func( node ast.Node)bool){
	ast.Inspect(this.File,inspect)
}
