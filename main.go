
package utils

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chainCode struct{}

func (t *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)

}
func (t *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func TestNamespacePutState(t *testing.T) {
	var value = []byte("value")

	mock := shim.NewMockStub("hyperledger-fabric-evmcc", &chainCode{})
	stub := NewNamespaceChainStubDecorator("prefix", mock)

	mock.MockTransactionStart("tx1")
	stub.PutState("key1", value)
	mock.MockTransactionEnd("tx1")

	valFromState, exists := mock.State["prefix:key1"]

	assert.True(t, exists)
	assert.Equal(t, value, valFromState)
}

func TestNamespaceGetState(t *testing.T) {
	var value = []byte("value")

	mock := shim.NewMockStub("hyperledger-fabric-evmcc", &chainCode{})
	stub := NewNamespaceChainStubDecorator("prefix", mock)

	mock.State["prefix:key1"] = value

	valFromState, err := stub.GetState("key1")

	assert.NoError(t, err)
	assert.Equal(t, value, valFromState)
}

func TestNamespaceCreateAndSplitCompositeKey(t *testing.T) {
	var srcArgs = []string{"a1", "a2"}

	mock := shim.NewMockStub("hyperledger-fabric-evmcc", &chainCode{})
	stub := NewNamespaceChainStubDecorator("p", mock)

	compKey, err := stub.CreateCompositeKey("key1", srcArgs)

	assert.NoError(t, err)
	assert.Equal(t, string([]byte{0, 0x70, 0x3a, 0x6b, 0x65, 0x79, 0x31, 0}), compKey[:8])

	objType, args, err := stub.SplitCompositeKey(compKey)
	assert.NoError(t, err)
	assert.Equal(t, "key1", objType)
	assert.Equal(t, srcArgs, args)
}

func TestNamespaceRangeIterator(t *testing.T) {
	mock := shim.NewMockStub("hyperledger-fabric-evmcc", &chainCode{})
	stub := NewNamespaceChainStubDecorator("prefix", mock)

	mock.MockTransactionStart("tx1")
	stub.PutState("key1", []byte{1})
	stub.PutState("key2", []byte{2})
	stub.PutState("key3", []byte{3})
	mock.MockTransactionEnd("tx1")

	valFromState := mock.State["prefix:key1"]
	assert.Equal(t, []byte{1}, valFromState)
	valFromState = mock.State["prefix:key2"]
	assert.Equal(t, []byte{2}, valFromState)
	valFromState = mock.State["prefix:key3"]
	assert.Equal(t, []byte{3}, valFromState)

	var index = 1

	iter, err := stub.GetStateByRange("key1", "key3")
	require.NoError(t, err)
	defer iter.Close()
	for iter.HasNext() {
		kv, err := iter.Next()
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf("key%d", index), kv.GetKey())
		assert.Equal(t, uint8(index), kv.GetValue()[0])

		index++
	}
	assert.Equal(t, index, 4)
}

