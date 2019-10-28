// Start fileheader template
// Code generated by go generate; DO NOT EDIT.
// This file was generated by FactomGenerate robots

// Start Generated Code

package generated_test

import (
	"testing"

	"github.com/FactomProject/factomd/common"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/generated"
)

// End fileheader template

// Start accountedqueue_test generated go code

func TestAccountedQueue(t *testing.T) {
	q := new(generated.Queue_IMsg).Init(common.NilName, "Test", 10)

	if q.Dequeue() != nil {
		t.Fatal("empty dequeue return non-nil")
	}

	for i := 0; i < 10; i++ {
		q.Enqueue(new(messages.Bounce))
	}

	// commented out because it requires a modern prometheus package
	//if testutil.ToFloat64(q.TotalMetric()) != float64(10) {
	//	t.Fatal("TotalMetric fail")
	//}

	for i := 9; i >= 0; i-- {
		q.Dequeue()
		// commented out because it requires a modern prometheus package
		//if testutil.ToFloat64(q.Metric()) != float64(i) {
		//	t.Fatal("Metric fail")
		//}
	}

	if q.Dequeue() != nil {
		t.Fatal("empty dequeue return non-nil")
	}
}

//
// Start filetail template
// Code generated by go generate; DO NOT EDIT.
// at <no value>
// End filetail template
// End Generated Code
