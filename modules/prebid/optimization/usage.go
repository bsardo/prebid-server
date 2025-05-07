package optimization

import (
	"github.com/prebid/prebid-server/v3/hooks/hookstage"
	"github.com/prebid/prebid-server/v3/openrtb_ext"
	"github.com/prebid/prebid-server/v3/rulestree"
	structs "github.com/prebid/prebid-server/v3/modules/prebid/optimization/config"
	"github.com/prebid/prebid-server/v3/modules/prebid/optimization/rulesengine"
)


// func main() {
// 	rw := openrtb_ext.RequestWrapper{}
// 	root := Node[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]]{}
// 	tree := Tree[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]]{Root: &root}

// 	changeSet := hookstage.ChangeSet[hookstage.BidderRequestPayload]{}

// 	// err := Execute[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]](&tree, &rw, &changeSet)
// 	err := tree.Run(&rw, &changeSet)
// 	if err != nil {
// 		return
// 	}


// 	builderFunc := builder[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]]
// 	tree2, err := NewTree[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]](builderFunc)
// }

// func builder[T1 any, T2 any](tree *Tree[T1, T2]) error {
// 	return nil
// }


// func main2() {
// 	rw := openrtb_ext.RequestWrapper{} //T1
// 	changeSet := hookstage.ChangeSet[hookstage.BidderRequestPayload]{} //T2

// 	tree := NewTree[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]]()
// 	// BuildTree(configJSON) --> 
// 	if err := tree.Valid(); err != nil {
// 		return
// 	}
// 	if err := tree.Run(&rw, &changeSet); err != nil {
// 		return
// 	}
// }

func main3() {
	rw := openrtb_ext.RequestWrapper{} //T1
	changeSet := hookstage.ChangeSet[hookstage.BidderRequestPayload]{} //T2

	myTreeBuilder := rulesengine.MyTreeBuilder[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]]{
		Config: structs.ModelGroup{},
		Sff:    rulestree.NewRequestSchemaFunction,
		Rff:    NewResultFunction,
	}

	tree, err := rulestree.NewTree[openrtb_ext.RequestWrapper, hookstage.ChangeSet[hookstage.BidderRequestPayload]](
		&myTreeBuilder,
	)
	if err != nil {
		return
	}
	if err := tree.Run(&rw, &changeSet); err != nil {
		return
	}
}

