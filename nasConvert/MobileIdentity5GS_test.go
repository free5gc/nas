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
			name: "SUSI-NAI",
			args: args{
				buf: []byte{0x11, 0x02, 0x58, 0x39, 0xf0, 0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10},
			},
			wantSuci:   "nai-1-025839f0ff010000000010",
			wantPlmnId: "",
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
			name: "SUSI-NAI-short",
			args: args{
				buf: []byte{0x11, 0x02},
			},
			wantSuci:   "nai-1-02",
			wantPlmnId: "",
			wantErr:    false,
		},
		{
			name: "SUSI-NAI-too-short",
			args: args{
				buf: []byte{0x11},
			},
			wantErr: true,
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
				PlmnId: &models.PlmnId{
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
				PlmnId: &models.PlmnId{
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
			name: "PEI-IMEI-even",
			args: args{
				buf: []byte{0x3, 0xf1},
			},
			want:    "imei-01",
			wantErr: false,
		},
		{
			name: "PEI-IMEISV-odd",
			args: args{
				buf: []byte{0xd, 0x21},
			},
			want:    "imeisv-012",
			wantErr: false,
		},
		{
			name:    "PEI-nil",
			wantErr: true,
		},
		{
			name: "PEI-IMEI-len1",
			args: args{
				buf: []byte{0xb},
			},
			want:    "imei-0",
			wantErr: false,
		},
		{
			name: "PEI-IMEI-len0",
			args: args{
				buf: []byte{0x3},
			},
			want:    "imei-",
			wantErr: false,
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
