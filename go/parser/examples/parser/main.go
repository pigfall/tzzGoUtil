package main

import(
	"log"
	"go/token"
	"go/ast"
	goparser "go/parser"
	"os"
	"io/ioutil"
	"github.com/pigfall/tzzGoUtil/go/parser"
)

func init(){
	log.SetFlags(log.Llongfile)

}

const fileContent =`
package main

var nameIdent string = "tzz"
` 

func main() {
	// < write example file
	filePath:="/tmp/tmp_go.go"
	err := ioutil.WriteFile(filePath,[]byte(fileContent),os.ModePerm)
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	// >

	fileset := token.NewFileSet()
	astFile,err := parser.ParseFile(fileset,filePath,goparser.ParseComments)
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	astFile.Inspect(func(node ast.Node)bool{
		declIfce,ok:= node.(ast.Decl)
		if !ok {
			log.Printf("%+v\n",node)
			return true
		}
		// < is declaretion
		generalDecl,ok := declIfce.(*ast.GenDecl)
		if !ok {
			return false
		}
		for _,spec := range generalDecl.Specs{
			valueSpec := spec.(*ast.ValueSpec)
			log.Printf("ident name %+v , type %+v \n",valueSpec.Names,valueSpec.Type)
		}
		log.Printf("%+v\n",generalDecl)
		// >
		return true
	})
}
