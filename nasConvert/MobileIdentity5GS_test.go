package nasConvert

import (
	"reflect"
	"testing"

	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
)

func TestSuciToStringWithError(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name       string
		args       args
		wantSuci   string
		wantPlmnId string
		wantErr    bool
	}{
		{
			name: "SUSI-null",
			args: args{
				buf: []byte{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf1},
			},
			wantSuci:   "suci-0-208-93-0-0-0-0000001",
			wantPlmnId: "20893",
			wantErr:    false,
		},
		{
			name: "SUSI-nonnull",
			args: args{
				buf: []byte{0x01, 0x02, 0x58, 0x39, 0xf0, 0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10},
			},
			wantSuci:   "suci-0-208-935-0-1-0-00000010",
			wantPlmnId: "208935",
			wantErr:    false,
		},
		{
			name: "SUSI-short",
			args: args{
				buf: []byte{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00},
			},
			wantSuci:   "suci-0-208-93-0-0-0-00",
			wantPlmnId: "20893",
			wantErr:    false,
		},
		{
			name: "SUSI-too-short",
			args: args{
				buf: []byte{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00},
			},
			wantErr: true,
		},
		{
			name: "SUSI-nil",
			args: args{
				buf: nil,
			},
			wantErr: true,
		},
		{
			name: "TS23003-28.7.3-Examples-SUCI-NAI-IMSI",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid678.schid0.userid0999999999@5gc.mnc015.mcc234.3gppnetwork.org")...),
			},
			wantSuci:   "suci-0-234-15-678-0-0-0999999999",
			wantPlmnId: "23415",
			wantErr:    false,
		},
		{
			name: "TS23003-28.7.3-Examples-SUCI-NAI-NSI",
			args: args{
				buf: append([]byte{0x11}, []byte("type1.rid678.schid0.useriduser17@example.com")...),
			},
			wantSuci:   "suci-1-example.com-678-0-0-user17",
			wantPlmnId: "",
			wantErr:    false,
		},
		{
			name: "SUCI-NAI-too-short",
			args: args{
				buf: []byte{0x11},
			},
			wantErr: true,
		},
		{
			name: "SUCI-NAI-invalid-username-format",
			args: args{
				buf: append([]byte{0x11}, []byte("username@example.com")...),
			},
			wantErr: true,
		},
		{
			name: "SUCI-NAI-invalid-missing-@",
			args: args{
				buf: append([]byte{0x11}, []byte("username.example.com")...),
			},
			wantErr: true,
		},
		{
			name: "SUCI-NAI-invalid-protection-scheme",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid678.schid3.username@example.com")...),
			},
			wantErr: true,
		},
		{
			name: "TS23003-28.7.6-trusted-non-3GPP-access",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid678.schid0.userid0999999999@nai.5gc.mnc001.mcc001.3gppnetwork.org")...),
			},
			wantSuci:   "suci-0-001-01-678-0-0-0999999999",
			wantPlmnId: "00101",
			wantErr:    false,
		},
		{
			name: "TS23003-28.7.7-trusted-non-3GPP-access-N5CW",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid678.schid0.userid0999999999@nai.5gc-nn.mnc001.mcc001.3gppnetwork.org")...),
			},
			wantSuci:   "suci-0-001-01-678-0-0-0999999999",
			wantPlmnId: "00101",
			wantErr:    false,
		},
		{
			name: "TS23003-28.7.12-NSWO",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid678.schid0.userid0999999999@5gc-nswo.nid1234.mnc001.mcc001.3gppnetwork.org")...),
			},
			wantSuci:   "suci-0-001-01-678-0-0-0999999999",
			wantPlmnId: "00101",
			wantErr:    false,
		},
		{
			name: "TS29503-AnnexC-Example1",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid012.schid0.userid0123456789@5gc.mnc045.mcc123.3gppnetwork.org")...),
			},
			wantSuci:   "suci-0-123-45-012-0-0-0123456789",
			wantPlmnId: "12345",
			wantErr:    false,
		},
		{
			name: "TS29503-AnnexC-Example2",
			args: args{
				buf: append([]byte{0x11}, []byte("type0.rid0002.schid1.hnkey17.ecckeye9b9916c911f448d8792e6b2f387f85d3ecab9040049427d9edbb5431b0bc711.cip023be6a057.macb45d936238aebeb7@5gc.mnc045.mcc123.3gppnetwork.org")...),
			},
			wantSuci:   "suci-0-123-45-0002-1-17-e9b9916c911f448d8792e6b2f387f85d3ecab9040049427d9edbb5431b0bc711023be6a057b45d936238aebeb7",
			wantPlmnId: "12345",
			wantErr:    false,
		},
		{
			name: "TS29503-AnnexC-Example3",
			args: args{
				buf: append([]byte{0x11}, []byte("type1.rid84.schid2.hnkey250.ecckeye9b9916c911f448d8792e6b2f387f85d3ecab9040049427d9edbb5431b0bc71195.cip023be6a057.macb45d936238aebeb7@example.com")...),
			},
			wantSuci:   "suci-1-example.com-84-2-250-e9b9916c911f448d8792e6b2f387f85d3ecab9040049427d9edbb5431b0bc71195023be6a057b45d936238aebeb7",
			wantPlmnId: "",
			wantErr:    false,
		},
		{
			name: "TS29503-AnnexC-Example4",
			args: args{
				buf: append([]byte{0x11}, []byte("type3.rid012.schid0.userid00-00-5E-00-53-00@operator.com")...),
			},
			wantSuci:   "suci-3-operator.com-012-0-0-00-00-5E-00-53-00",
			wantPlmnId: "",
			wantErr:    false,
		},
		{
			name: "TS29503-AnnexC-Example5",
			args: args{
				buf: append([]byte{0x11}, []byte("type1.rid3456.schid0.useridanonymous@operator.com")...),
			},
			wantSuci:   "suci-1-operator.com-3456-0-0-anonymous",
			wantPlmnId: "",
			wantErr:    false,
		},
		{
			name: "TS29503-AnnexC-Example5-Alternative",
			args: args{
				buf: append([]byte{0x11}, []byte("type1.rid3456.schid0.userid@operator.com")...),
			},
			wantSuci:   "suci-1-operator.com-3456-0-0-",
			wantPlmnId: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotSuci, gotPlmnId, err := SuciToStringWithError(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuciToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSuci != tt.wantSuci {
				t.Errorf("SuciToString() gotSuci = %v, want %v", gotSuci, tt.wantSuci)
			}
			if gotPlmnId != tt.wantPlmnId {
				t.Errorf("SuciToString() gotPlmnId = %v, want %v", gotPlmnId, tt.wantPlmnId)
			}
		})
	}
}

func TestGutiToStringWithError(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name      string
		args      args
		wantGuami models.Guami
		wantGuti  string
		wantErr   bool
	}{
		{
			name: "GUTI-MNC2",
			args: args{
				buf: []byte{0xf2, 0x02, 0xf8, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23},
			},
			wantGuami: models.Guami{
				PlmnId: &models.PlmnIdNid{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "012345",
			},
			wantGuti: "2089301234567890123",
			wantErr:  false,
		},
		{
			name: "GUTI-MNC3",
			args: args{
				buf: []byte{0xf2, 0x02, 0x58, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23},
			},
			wantGuami: models.Guami{
				PlmnId: &models.PlmnIdNid{
					Mcc: "208",
					Mnc: "935",
				},
				AmfId: "012345",
			},
			wantGuti: "20893501234567890123",
			wantErr:  false,
		},
		{
			name: "GUTI-too-long",
			args: args{
				buf: []byte{0xf2, 0x02, 0xf8, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45},
			},
			wantErr: true,
		},
		{
			name: "GUTI-too-short",
			args: args{
				buf: []byte{0xf2, 0x02, 0xf8, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotGuami, gotGuti, err := GutiToStringWithError(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("GutiToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGuami, tt.wantGuami) {
				t.Errorf("GutiToString() gotGuami = %v, want %v", gotGuami, tt.wantGuami)
			}
			if gotGuti != tt.wantGuti {
				t.Errorf("GutiToString() gotGuti = %v, want %v", gotGuti, tt.wantGuti)
			}
		})
	}
}

func TestGutiToNasWithError(t *testing.T) {
	type args struct {
		guti string
	}
	tests := []struct {
		name    string
		args    args
		want    nasType.GUTI5G
		wantErr bool
	}{
		{
			name: "GUTI-MNC2",
			args: args{
				guti: "2089301234567890123",
			},
			want: nasType.GUTI5G{
				Iei:   0,
				Len:   11,
				Octet: [11]uint8{0xf2, 0x02, 0xf8, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23},
			},
			wantErr: false,
		},
		{
			name: "GUTI-MNC3",
			args: args{
				guti: "20893501234567890123",
			},
			want: nasType.GUTI5G{
				Iei:   0,
				Len:   11,
				Octet: [11]uint8{0xf2, 0x02, 0x58, 0x39, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23},
			},
			wantErr: false,
		},
		{
			name: "GUTI-too-long",
			args: args{
				guti: "208935012345678901234",
			},
			wantErr: true,
		},
		{
			name: "GUTI-too-short",
			args: args{
				guti: "208930123456789012",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-MCC1",
			args: args{
				guti: "x089301234567890123",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-MCC2",
			args: args{
				guti: "2x89301234567890123",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-MCC3",
			args: args{
				guti: "20x9301234567890123",
			},
			wantErr: true,
		}, {
			name: "GUTI-bad-MNC1",
			args: args{
				guti: "208x301234567890123",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-MNC2",
			args: args{
				guti: "2089x01234567890123",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-MNC3",
			args: args{
				guti: "20893x01234567890123",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-TMSI",
			args: args{
				guti: "208930123456789012x",
			},
			wantErr: true,
		},
		{
			name: "GUTI-bad-AMFID",
			args: args{
				guti: "2089301x34567890123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GutiToNasWithError(tt.args.guti)
			if (err != nil) != tt.wantErr {
				t.Errorf("GutiToNas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GutiToNas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeiToStringWithError(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Complete-Valid-IMEI",
			args: args{
				// Example encoding for a valid 15-digit IMEI
				buf: []byte{
					0x4b, 0x09, 0x51, 0x24, 0x30, 0x32, 0x57, 0x81,
				},
			},
			want:    "imei-490154203237518",
			wantErr: false,
		},
		{
			name: "Complete-Ivalid-IMEI",
			args: args{
				// Not valid 15-digit IMEI: CD(Check Digit) not valid
				buf: []byte{
					0x4b, 0x09, 0x51, 0x24, 0x30, 0x32, 0x57, 0x82,
				},
			},
			wantErr: true,
		},
		{
			name: "Complete-Valid-IMEISV",
			args: args{
				buf: []byte{
					0x90, 0x87, 0x65, 0x43, 0x21, 0x01, 0x23, 0x45, 0x60,
				},
			},
			want:    "imeisv-9785634121032540",
			wantErr: false,
		},
		{
			name: "IMEI-TooLong",
			args: args{
				buf: []byte{
					0x4b, 0x09, 0x51, 0x24, 0x30, 0x32, 0x57, 0x81, 0x20,
				},
			},
			wantErr: true,
		},
		{
			name: "IMEI-TooShort",
			args: args{
				buf: []byte{
					0x4b, 0x09, 0x51, 0x24, 0x30, 0x32,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := PeiToStringWithError(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeiToStringWithError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PeiToStringWithError() = %v, want %v", got, tt.want)
			}
		})
	}
}
