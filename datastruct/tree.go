package datastruct

import(
		"context"
		"fmt"
		"strings"
)

type Tree struct{
	rootNode *Node
}

type Node struct{
	parent *Node
	childreen []*Node
	value interface{}
}

func (this *Node) GetParent()*Node{
	return this.parent
}

func (this *Node) GetValue()interface{}{
	return this.value
}

type TreeBuilderI interface{
	GetValue(ctx context.Context) (interface{},error)
	GetChildren(ctx context.Context) ([]TreeBuilderI,error)
}


func BuildTree(ctx context.Context,builder TreeBuilderI)(*Tree,error){
	rootNode,err := buildTree(ctx,builder)
	if err != nil{
		return nil,err
	}
	return &Tree{
		rootNode:rootNode,
	},nil
}

func buildTree(ctx context.Context,builder TreeBuilderI)(*Node,error){
	value,err := builder.GetValue(ctx)
	if err != nil{
		return nil,err
	}
	childreen,err := builder.GetChildren(ctx)
	if err != nil{
		return nil,err
	}
	rootNode := &Node{
		value:value,
	}
	nodes :=make([]*Node,0,len(childreen))
	for _,c := range childreen{
		node,err := buildTree(ctx,c)
		if err != nil{
			return nil,err
		}
		node.parent = rootNode
		nodes = append(nodes,node)
	}
	rootNode.childreen =nodes
	return rootNode,nil
}

func (this *Tree) BredthFirstDo(
	ctx context.Context,
	do func(node *Node,)(error),
	onEnterNextBredth func(),
)(error){
	fmt.Printf("rootNode %p \n",this.rootNode)
	rest,err := this.rootNode.BredthFirstDo(ctx,do)
	if err != nil{
		return err
	}
	for{
		if len(rest) == 0{
			return nil
		}
		if onEnterNextBredth != nil{
			onEnterNextBredth()
		}
		var nextRest = make([]*Node,0,len(rest))
		for _,node := range rest{
			restTmp,err := node.BredthFirstDo(ctx,do)
			if err != nil{
				return err
			}
			nextRest = append(nextRest,restTmp...)
		}
		rest = nextRest
	}
}

func (this *Node) BredthFirstDo(
	ctx context.Context,
	do func(node *Node)(error),
)(rest []*Node,err error){
	err = do(this)
	if err != nil{
		return nil,err
	}
	return this.childreen,nil
}



func (this *Tree) ToString(ctx context.Context)string{
	strBuilder := strings.Builder{}
	this.BredthFirstDo(
		ctx,
		func(node *Node)error{
			strBuilder.WriteString(fmt.Sprintf("< Node %p | %+v >",node,node))
			return nil
		},
		func(){
			strBuilder.WriteString("\n")
			strBuilder.WriteString("\n")
		},
	)
	return strBuilder.String()
}

func (this *Tree) DepthFirstDo(ctx context.Context,f func(ctx context.Context,node *Node)error)error{
	return this.rootNode.DepthFirstDo(ctx,f)
}

func(this *Node) DepthFirstDo(ctx context.Context,f func(ctx context.Context,node *Node)error)error{
	err := f(ctx,this)
	if err!=nil{
		return err
	}
	for _,child := range this.childreen{
		err := child.DepthFirstDo(ctx,f)
		if err != nil{
			return err
		}
	}
	return nil
}

