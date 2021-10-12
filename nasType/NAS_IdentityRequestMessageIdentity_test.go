package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasType"
)

type nasTypeIdentityRequestMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeIdentityRequestMessageIdentityTable = []nasTypeIdentityRequestMessageIdentityData{
	{nas.MsgTypeIdentityRequest, nas.MsgTypeIdentityRequest},
}

func TestNasTypeNewIdentityRequestMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityRequestMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeIdentityRequestMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityRequestMessageIdentity()
	for _, table := range nasTypeIdentityRequestMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
