package rulestree

import ()

// ResultFunction...
// e.g. hookstage.ChangeSet[hookstage.BidderRequestPayload]
type ResultFunction[T any] interface {
	Call(payload *T) error
}
