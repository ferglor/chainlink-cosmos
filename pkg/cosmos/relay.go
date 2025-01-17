package cosmos

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters/cosmwasm"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters/injective"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/params"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/txm"
)

const (
	InjectivePrefix string = "inj"
)

// ErrMsgUnsupported is returned when an unsupported type of message is encountered.
// Deprecated: use txm.ErrMsgUnsupported
type ErrMsgUnsupported = txm.ErrMsgUnsupported

var _ types.Relayer = &Relayer{} //nolint:staticcheck

type Relayer struct {
	lggr   logger.Logger
	chain  adapters.Chain
	ctx    context.Context
	cancel func()
}

// Note: constructed in core
func NewRelayer(lggr logger.Logger, chain adapters.Chain) *Relayer {
	ctx, cancel := context.WithCancel(context.Background())

	bech32Prefix := chain.Config().Bech32Prefix()
	gasToken := chain.Config().GasToken()
	params.InitCosmosSdk(
		bech32Prefix,
		gasToken,
	)

	return &Relayer{
		lggr:   lggr,
		chain:  chain,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
}

// Start starts the relayer respecting the given context.
func (r *Relayer) Start(context.Context) error {
	if r.chain == nil {
		return errors.New("Cosmos unavailable")
	}
	return nil
}

func (r *Relayer) Close() error {
	r.cancel()
	return nil
}

func (r *Relayer) Ready() error {
	return r.chain.Ready()
}

func (r *Relayer) HealthReport() map[string]error {
	return r.chain.HealthReport()
}

func (r *Relayer) NewMercuryProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	return nil, errors.New("mercury is not supported for cosmos")
}

func (r *Relayer) NewFunctionsProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.FunctionsProvider, error) {
	return nil, errors.New("functions are not supported for cosmos")
}

func (r *Relayer) NewConfigProvider(args types.RelayArgs) (types.ConfigProvider, error) {
	var configProvider types.ConfigProvider
	var err error
	if r.chain.Config().Bech32Prefix() == InjectivePrefix {
		configProvider, err = injective.NewConfigProvider(r.ctx, r.lggr, r.chain, args)
	} else {
		// Default to cosmwasm adapter
		configProvider, err = cosmwasm.NewConfigProvider(r.ctx, r.lggr, r.chain, args)
	}
	if err != nil {
		return nil, err
	}

	return configProvider, err
}

func (r *Relayer) NewMedianProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	configProvider, err := cosmwasm.NewMedianProvider(r.ctx, r.lggr, r.chain, rargs, pargs)
	if err != nil {
		return nil, err
	}
	return configProvider, err
}
