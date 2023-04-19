package config

import (
	"net/url"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	relayconfig "github.com/smartcontractkit/chainlink-relay/pkg/config"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
)

// Global defaults.
var defaultConfigSet = configSet{
	BlockRate: 6 * time.Second,
	// ~6s per block, so ~3m until we give up on the tx getting confirmed
	// Anecdotally it appears anything more than 4 blocks would be an extremely long wait,
	// In practice during the UST depegging and subsequent extreme congestion, we saw
	// ~16 block FIFO lineups.
	BlocksUntilTxTimeout:  30,
	ConfirmPollPeriod:     time.Second,
	FallbackGasPriceUAtom: sdk.MustNewDecFromStr("0.015"),
	// This is high since we simulate before signing the transaction.
	// There's a chicken and egg problem: need to sign to simulate accurately
	// but you need to specify a gas limit when signing.
	// TODO: Determine how much gas a signature adds and then
	// add that directly so we can be more accurate.
	GasLimitMultiplier: client.DefaultGasLimitMultiplier,
	// The max gas limit per block is 1_000_000_000
	// https://github.com/terra-money/core/blob/d6037b9a12c8bf6b09fe861c8ad93456aac5eebb/app/legacy/migrate.go#L69.
	// The max msg size is 10KB https://github.com/terra-money/core/blob/d6037b9a12c8bf6b09fe861c8ad93456aac5eebb/x/wasm/types/params.go#L15.
	// Our msgs are only OCR reports for now, which will not exceed that size.
	// There appears to be no gas limit per tx, only per block, so theoretically
	// we could include 1000 msgs which use up to 1M gas.
	// To be conservative and since the number of messages we'd
	// have in a batch on average roughly corresponds to the number of terra ocr jobs we're running (do not expect more than 100),
	// we can set a max msgs per batch of 100.
	MaxMsgsPerBatch:     100,
	OCR2CachePollPeriod: 4 * time.Second,
	OCR2CacheTTL:        time.Minute,
	TxMsgTimeout:        10 * time.Minute,
}

type Config interface {
	BlockRate() time.Duration
	BlocksUntilTxTimeout() int64
	ConfirmPollPeriod() time.Duration
	FallbackGasPriceUAtom() sdk.Dec
	FCDURL() url.URL
	GasLimitMultiplier() float64
	MaxMsgsPerBatch() int64
	OCR2CachePollPeriod() time.Duration
	OCR2CacheTTL() time.Duration
	TxMsgTimeout() time.Duration
}

// opt: remove
type configSet struct {
	BlockRate             time.Duration
	BlocksUntilTxTimeout  int64
	ConfirmPollPeriod     time.Duration
	FallbackGasPriceUAtom sdk.Dec
	FCDURL                url.URL
	GasLimitMultiplier    float64
	MaxMsgsPerBatch       int64
	OCR2CachePollPeriod   time.Duration
	OCR2CacheTTL          time.Duration
	TxMsgTimeout          time.Duration
}

type Chain struct {
	BlockRate             *utils.Duration
	BlocksUntilTxTimeout  *int64
	ConfirmPollPeriod     *utils.Duration
	FallbackGasPriceUAtom *decimal.Decimal
	FCDURL                *utils.URL
	GasLimitMultiplier    *decimal.Decimal
	MaxMsgsPerBatch       *int64
	OCR2CachePollPeriod   *utils.Duration
	OCR2CacheTTL          *utils.Duration
	TxMsgTimeout          *utils.Duration
}

func (c *Chain) SetDefaults() {
	if c.BlockRate == nil {
		c.BlockRate = utils.MustNewDuration(defaultConfigSet.BlockRate)
	}
	if c.BlocksUntilTxTimeout == nil {
		c.BlocksUntilTxTimeout = &defaultConfigSet.BlocksUntilTxTimeout
	}
	if c.ConfirmPollPeriod == nil {
		c.ConfirmPollPeriod = utils.MustNewDuration(defaultConfigSet.ConfirmPollPeriod)
	}
	if c.FallbackGasPriceUAtom == nil {
		d := decimal.NewFromBigInt(defaultConfigSet.FallbackGasPriceUAtom.BigInt(), -sdk.Precision)
		c.FallbackGasPriceUAtom = &d
	}
	if c.FCDURL == nil {
		c.FCDURL = (*utils.URL)(&defaultConfigSet.FCDURL)
	}
	if c.GasLimitMultiplier == nil {
		d := decimal.NewFromFloat(defaultConfigSet.GasLimitMultiplier)
		c.GasLimitMultiplier = &d
	}
	if c.MaxMsgsPerBatch == nil {
		c.MaxMsgsPerBatch = &defaultConfigSet.MaxMsgsPerBatch
	}
	if c.OCR2CachePollPeriod == nil {
		c.OCR2CachePollPeriod = utils.MustNewDuration(defaultConfigSet.OCR2CachePollPeriod)
	}
	if c.OCR2CacheTTL == nil {
		c.OCR2CacheTTL = utils.MustNewDuration(defaultConfigSet.OCR2CacheTTL)
	}
	if c.TxMsgTimeout == nil {
		c.TxMsgTimeout = utils.MustNewDuration(defaultConfigSet.TxMsgTimeout)
	}
}

type Node struct {
	Name          *string
	TendermintURL *utils.URL
}

func (n *Node) ValidateConfig() (err error) {
	if n.Name == nil {
		err = multierr.Append(err, relayconfig.ErrMissing{Name: "Name", Msg: "required for all nodes"})
	} else if *n.Name == "" {
		err = multierr.Append(err, relayconfig.ErrEmpty{Name: "Name", Msg: "required for all nodes"})
	}
	if n.TendermintURL == nil {
		err = multierr.Append(err, relayconfig.ErrMissing{Name: "TendermintURL", Msg: "required for all nodes"})
	}
	return
}