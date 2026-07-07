package message

import (
	"testing"

	"github.com/free5gc/nas/message"
)

/*
	This test and testdata/* are generated and patched from https://github.com/free5gc/nas/pull/18

	Notice:
    Compared to free5gc nas which used raw bytes as encoding input,
    this fork uses semantic params which rely on manual implementation,
    so fuzz test here can't be as strict as that of original PR.
    For example, if a fuzz test case contains unimplemented ies, marshal/parse may
    have error raised by unimplementation.
    Therefore, encoding/decoding error is not check here in this fuzz-test. As long as the procedure
    doesn't cause panic, the test is treated as pass.
*/

// TODO: to fully test with fuzz-test, remove the defer func() part in Parse() and run with
// 'go test -fuzz=FuzzNAS' to see if panic happens.
func FuzzNAS(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		m, err := message.Parse(d, nil)
		if err == nil {
			// test re-encoding/re-decoding if the fuzz data is in valid format
			buf, _ := m.MarshalBinary() //nolint
			// require.NoError(t, err, "Re-encoding failed") // see Notice.
			_, _ = message.Parse(buf, nil) //nolint
			// require.NoError(t, err, "Re-decoding failed") // see Notice.

			// TODO: uncomment this to check for panic
			// if strings.HasPrefix(err.Error(), "Parse(): panic") {
			// 	require.NoError(t, err)
			// }
		} // TODO: uncomment this to check for panic
		// else {
		// TODO: uncomment this to check for panic
		// if strings.HasPrefix(err.Error(), "Parse(): panic") {
		// 	require.NoError(t, err)
		// }
		// }
	})
}

// TODO: to fully test with fuzz-test, remove the defer func() part in ParseGMM() and run with
// 'go test -fuzz=FuzzGmmMessageDecode' to see if panic happens.
func FuzzGmmMessageDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		m, err := message.ParseGMM(d)
		if err == nil {
			// test re-encoding/re-decoding if the fuzz data is in valid format
			buf, _ := m.MarshalBinary() //nolint
			// require.NoError(t, err, "Re-encoding failed") // see Notice.
			_, _ = message.ParseGMM(buf) //nolint
			// require.NoError(t, err, "Re-decoding failed") // see Notice.

			// TODO: uncomment this to check for panic
			// if strings.HasPrefix(err.Error(), "ParseGMM(): panic") {
			// 	require.NoError(t, err)
			// }
		} // TODO: uncomment this to check for panic
		// else {
		// TODO: uncomment this to check for panic
		// if strings.HasPrefix(err.Error(), "ParseGMM(): panic") {
		// 	require.NoError(t, err)
		// }
		// }
	})
}

// TODO: to fully test with fuzz-test, remove the defer func() part in ParseGSM() and run with
// 'go test -fuzz=FuzzGsmMessageDecode' to see if panic happens.
func FuzzGsmMessageDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		m, err := message.ParseGSM(d)
		if err == nil {
			// test re-encoding/re-decoding if the fuzz data is in valid format
			buf, _ := m.MarshalBinary() //nolint
			// require.NoError(t, err, "Re-encoding failed") // see Notice.
			_, _ = message.ParseGMM(buf) //nolint
			// require.NoError(t, err, "Re-decoding failed") // see Notice.

			// TODO: uncomment this to check for panic
			// if strings.HasPrefix(err.Error(), "ParseGSM(): panic") {
			// 	require.NoError(t, err)
			// }
		} // TODO: uncomment this to check for panic
		//else {
		// if strings.HasPrefix(err.Error(), "ParseGSM(): panic") {
		// 	require.NoError(t, err)
		// }
		//}
	})
}
