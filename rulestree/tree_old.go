package rulestree

import (
)

// variable types:
// schema functions: payload (e.g. requestWrapper)
// result functions: payload (e.g. changeset of certain type, some arbitrary function)
//

// type Tree [T1 any, T2 any] struct {
// 	Root *Node[T1, T2]
// }

// type Node [T1 any, T2 any] struct {
// 	SchemaFunction  SchemaFunction[T1]
// 	ResultFunctions []ResultFunction[T2]
// 	Children        map[string]*Node[T1, T2]
// }

// This can either return all of the result functions OR 
// give it a pointer to the thing the result functions should modify as a parameter
// func Execute [T1 any, T2 any](r *Tree[T1, T2], rw *T1) (*hookstage.ChangeSet[hookstage.BidderRequestPayload], error) {
func Execute[T1 any, T2 any](r *Tree[T1, T2], rw *T1, result *T2) error {
	currNode := r.Root

	for len(currNode.Children) > 0 {
		res, err := currNode.SchemaFunction.Call(rw)
		if err != nil {
			return err
		}

		next := currNode.Children[res]

		if next != nil {
			currNode = next
		}
	}

	for _, rf := range currNode.ResultFunctions {
		err := rf.Call(result)
		if err != nil {
			return err
		}
	}
	return nil
}

// func BuildRulesTree[T1 any, T2 any](conf structs.ModelGroup) (*Tree[T1, T2], error) {

// 	currNode := &Node[T1, T2]{}
// 	rules := Tree[T1, T2]{Root: currNode}

// 	for _, rule := range conf.Rule {
// 		for ci, condition := range rule.Conditions {

// 			if len(currNode.Children) == 0 {
// 				currNode.Children = make(map[string]*Node[T1, T2], 0)
// 				f, err := NewSchemaFunctionFactory(conf.Schema[ci].Func, conf.Schema[ci].Args)
// 				if err != nil {
// 					return nil, err
// 				}
// 				currNode.SchemaFunction = f
// 			}

// 			_, ok := currNode.Children[condition]
// 			if ok {
// 				currNode = currNode.Children[condition]
// 			} else {
// 				nextNode := &Node[T1, T2]{}
// 				currNode.Children[condition] = nextNode
// 				currNode = nextNode
// 			}
// 		}

// 		// array of func
// 		for _, res := range rule.Results {
// 			resFunc, err := NewResultFunctionFactory(res.Func, res.Args)
// 			if err != nil {
// 				return nil, err
// 			}
// 			currNode.ResultFunctions = append(currNode.ResultFunctions, resFunc)
// 		}

// 		currNode = rules.Root

// 	}

// 	return &rules, nil
// }
