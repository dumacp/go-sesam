package sam

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dumacp/smartcard"
	"github.com/dumacp/smartcard/nxp/mifare/samav2"
	"github.com/dumacp/smartcard/pcsc"
)

func Test_enableKeys(t *testing.T) {

	ctx, err := pcsc.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	rs, err := ctx.ListReaders()
	if err != nil {
		t.Fatal(err)
	}
	var r pcsc.Reader

	for _, v := range rs {
		if strings.Contains(v, "SAM") {
			r = pcsc.NewReader(ctx, v)
			break
		}
	}
	if r == nil {
		t.Fatalf("reader SAM not found")
	}

	card, err := r.ConnectSamCard()
	if err != nil {
		t.Fatal(err)
	}

	sam := samav2.SamAV2(card)

	keyMaster := make([]byte, 16)
	if _, err := sam.AuthHostAV2(keyMaster, 0, 0, 0); err != nil {
		t.Fatal(err)
	}
	t.Log("Auth!!!")

	type args struct {
		se samav2.SamAv2
	}
	tests := []struct {
		name    string
		args    args
		want    map[int]int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				se: sam,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := enableKeys(tt.args.se)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("enableKeys() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_samAv2_Decrypt(t *testing.T) {
	type fields struct {
		dev  samav2.SamAv2
		card smartcard.ICard
	}
	type args struct {
		data     []byte
		iv       []byte
		divInput []byte
		slot     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &samAv2{
				dev:  tt.fields.dev,
				card: tt.fields.card,
			}
			got, err := s.Decrypt(tt.args.data, tt.args.iv, tt.args.divInput, tt.args.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("samAv2.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("samAv2.Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
