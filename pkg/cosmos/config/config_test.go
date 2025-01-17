package config

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
)

func Test_sdkDecFromDecimal(t *testing.T) {
	tests := []string{
		"0.0",
		"0.1",
		"1.0",
		"0.000000000000000001",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			val := decimal.RequireFromString(tt)
			exp := sdk.MustNewDecFromStr(tt)
			assert.Equal(t, exp, sdkDecFromDecimal(&val))
		})
	}
}

func TestCosmosConfig_GetNode(t *testing.T) {
	type fields struct {
		ChainID *string
		Nodes   Nodes
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Node
		wantErr bool
	}{
		{
			name: "not found",
			args: args{
				name: "not a node",
			},
			fields:  fields{Nodes: Nodes{}},
			want:    db.Node{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				name: "node",
			},
			fields: fields{
				ChainID: ptr("chainID"),
				Nodes: []*Node{
					&Node{
						Name:          ptr("node"),
						TendermintURL: &config.URL{},
					},
				}},
			want: db.Node{
				CosmosChainID: "chainID",
				Name:          "node",
				TendermintURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TOMLConfig{
				Nodes:   tt.fields.Nodes,
				ChainID: tt.fields.ChainID,
			}
			got, err := c.GetNode(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("TOMLConfig.GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TOMLConfig.GetNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
