package terra

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/websocket"
	"github.com/terra-money/terra.go/client"
	"github.com/terra-money/terra.go/key"
	"github.com/terra-money/terra.go/msg"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tidwall/gjson"
)

type WsConn struct {
	mu   *sync.Mutex
	conn *websocket.Conn
}

type RpcRequest struct {
	Jsonrpc string
	Method  string
	Params  []interface{}
	ID      string
}

type RpcMessage struct {
	Error string
	Data  []byte
}

type Client struct {
	close     chan struct{}
	codec     *codec.LegacyAmino
	httpURL   string
	chainID   string
	wsURL     string
	ws        WsConn
	wsStarted bool
	subs      map[string]subscription
	subUnsub  map[string]chan<- Events
	queryCh   chan RpcMessage
	txCh      chan RpcMessage
	Height    uint64
	Log       Logger
}

func NewClient(spec OCR2Spec, lggr Logger) (Client, error) {
	return Client{
		close:    make(chan struct{}),
		codec:    codec.NewLegacyAmino(),
		chainID:  spec.ChainID,
		httpURL:  spec.NodeEndpointHTTP,
		wsURL:    spec.NodeEndpointWS,
		subs:     make(map[string]subscription),
		subUnsub: make(map[string]chan<- Events),
		queryCh:  make(chan RpcMessage),
		txCh:     make(chan RpcMessage),
		Log:      lggr,
	}, nil
}

func (c Client) LCD(gasPrice msg.DecCoin, gasAdjustment msg.Dec, signer key.PrivKey, timeout time.Duration) *client.LCDClient {
	return client.NewLCDClient(c.httpURL, c.chainID, gasPrice, gasAdjustment, signer, timeout)
}

type BroadcastTxMethod string

const (
	// Returns with response from CheckTx, does not wait for DeliverTx
	BroadcastSync BroadcastTxMethod = "broadcast_tx_sync"
	// Returns right away with no response
	BroadcastAsync BroadcastTxMethod = "broadcast_tx_async"
	// Returns with response from CheckTx and DeliverTx
	BroadcastBlock BroadcastTxMethod = "broadcast_tx_commit"
)

// Send Terra transaction to blockchain
func (c Client) Send(ctx context.Context, txBytes []byte, method BroadcastTxMethod) (*BroadcastTxResponse, error) {
	err := c.ensureWsConnection()
	if err != nil {
		return nil, err
	}

	c.Log.Infof("[broadcast transaction] %s,", string(method))

	payload := RpcRequest{
		Jsonrpc: "2.0",
		Method:  string(method),
		Params:  []interface{}{txBytes},
		ID:      string(method),
	}

	if err := c.ws.conn.WriteJSON(payload); err != nil {
		return nil, err
	}

	message := <-c.txCh
	if message.Error != "" {
		return nil, fmt.Errorf("Broadcast error: %s", message.Error)
	}

	var tx BroadcastTxResponse
	if err = json.Unmarshal(message.Data, &tx); err != nil {
		return nil, fmt.Errorf("Broadcast error: %s", err)
	}

	return &tx, nil
}

type QueryType string

type ABCIQueryParams struct {
	ContractID string
	Msg        []byte
}

func NewAbciQueryParams(contractID string, msg []byte) ABCIQueryParams {
	return ABCIQueryParams{contractID, msg}
}

const (
	TX   QueryType = "tx_search"
	ABCI QueryType = "abci_query"
)

var defaultAbciQueryParameters = []interface{}{"0", false}
var defaultTxQueryParameters = []interface{}{false, "1", "30", "desc"}

func (c Client) parseParameters(method QueryType, params []interface{}) ([]interface{}, error) {
	paramsLen := len(params)
	// check min params and append default parameters
	if method == ABCI {
		// 2 is minimum parameters that should be passed for abci queries
		if paramsLen < 2 {
			return nil, fmt.Errorf("Query error: not enough query parameters were passed")
		}
		// 4 is the length of required parameters for an abci query
		if paramsLen < 4 {
			// calculate how many default parameters to append, in case if more than required are passed
			params = append(params, defaultAbciQueryParameters[paramsLen-2:]...)
		}

		// use amino codec to encode abci parameters
		bz, err := c.codec.MarshalJSON(params[1])
		if err != nil {
			return nil, fmt.Errorf("Query error: %s", err)
		}
		params[1] = hex.EncodeToString(bz)
	}

	// check min params and append default parameters
	if method == TX {
		// 1 is minimum parameters should be passed for tx queries
		if paramsLen < 1 {
			return nil, fmt.Errorf("Query error: not enough query parameters were passed")
		}
		// 5 is the length of required parameters for a tx query
		if paramsLen < 5 {
			params = append(params, defaultTxQueryParameters[paramsLen-1:]...)
		}

	}

	return params, nil
}

func (c Client) Query(method QueryType, params []interface{}) ([]byte, error) {
	err := c.ensureWsConnection()
	if err != nil {
		return nil, err
	}
	// Should find a way to remove it
	// but for now solves: panic: concurrent write to websocket connection
	c.ws.mu.Lock()
	defer c.ws.mu.Unlock()
	c.Log.Infof("[query] %s, %s", string(method), params)

	params, err = c.parseParameters(method, params)
	if err != nil {
		return nil, err
	}
	payload := RpcRequest{
		Jsonrpc: "2.0",
		Method:  string(method),
		Params:  params,
		ID:      string(method),
	}

	if err := c.ws.conn.WriteJSON(payload); err != nil {
		return nil, err
	}

	// wait for a query response
	message := <-c.queryCh

	if message.Error != "" {
		return nil, fmt.Errorf("Query error: %s", message.Error)
	}

	return message.Data, nil
}

type subscription struct {
	Channel chan<- Events
	Payload request
}

type request struct {
	Jsonrpc string
	Method  string
	Params  []string
	ID      string
}

func (c *Client) HeadTracker() error {
	fq := []string{"tm.event='NewBlock'"}

	// create new block subscription
	channel := make(chan Events)
	if err := c.subscribe(context.TODO(), "head-tracker", fq, channel); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case height := <-channel:
				// update block number in the client
				c.Log.Infof("[head-tracker] Block height %d", height.Block)
				c.Height = height.Block
			case <-c.close:
				return
			}
		}
	}()

	c.Log.Infof("[head-tracker] Subscription started")
	return nil
}

// Subscribe to Terra events (address specifics)
func (c *Client) Subscribe(ctx context.Context, jobID string, address types.AccAddress, channel chan Events) error {
	fq := []string{fmt.Sprintf("tm.event='Tx' AND execute_contract.contract_address='%s'", address)}

	if err := c.subscribe(ctx, jobID, fq, channel); err != nil {
		return err
	}
	c.Log.Infof("[%s] Subscription created for %s", jobID, address)
	return nil
}

type Events struct {
	SubErr string   // error string if encountered during sub + unsub
	Block  uint64   // handle block number event
	Events []string // handle tx events
}

// create websocket connection and read if not created yet
func (c *Client) ensureWsConnection() error {
	if !c.wsStarted {
		ws, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
		if err != nil {
			return err
		}

		c.ws = WsConn{
			conn: ws,
			mu:   &sync.Mutex{},
		}
		// TODO: does the websocket connection need a close handler to resubscribe?
		// https://github.com/smartcontractkit/chainlink-terra/issues/24

		// start listening
		go c.listen()
		c.Log.Infof("Websocket connection opened to %s", c.wsURL)
		c.wsStarted = true
	}
	return nil
}

// base level subscribe function
func (c *Client) subscribe(ctx context.Context, jobID string, filterQuery []string, channel chan Events) error {
	// check if job id already exists (return error if it does)
	if _, exists := c.subs[jobID]; exists {
		return fmt.Errorf("[%s] Subscription already exists", jobID)
	}

	err := c.ensureWsConnection()
	if err != nil {
		return err
	}

	// send subscribe payload
	payload := request{
		Jsonrpc: "2.0",
		Method:  "subscribe",
		Params:  filterQuery,
		ID:      jobID,
	}
	// save response channel to subscription
	c.subs[jobID] = subscription{Channel: channel, Payload: payload}

	// set up chanel to catch first response
	subUnsub := make(chan Events)
	c.subUnsub[jobID] = subUnsub

	// send payload
	if err := c.ws.conn.WriteJSON(payload); err != nil {
		return err
	}

	// wait for successful subscription message
	msg := <-subUnsub
	if msg.SubErr != "" {
		return fmt.Errorf("[%s] WS error: %s", jobID, msg.SubErr)
	}
	return nil
}

// listen is the message
func (c *Client) listen() {
	// TODO: Need a way to prevent "use of closed network connection" error when closing connection
	// https://github.com/smartcontractkit/chainlink-terra/issues/25
	for {
		_, messageBytes, err := c.ws.conn.ReadMessage()
		if err != nil {
			c.Log.Errorf("[client/listen/read]: %s", err)
			return
		}
		message := string(messageBytes)
		// parse job id
		jobID := gjson.Get(message, "id")

		// if in a sub or unsub state, return error string
		if c.subUnsub[jobID.Str] != nil {
			error := gjson.Get(message, "error")
			c.subUnsub[jobID.Str] <- Events{SubErr: error.Raw}

			// remove sub/unsub state
			close(c.subUnsub[jobID.Str])
			c.subUnsub[jobID.Str] = nil
			continue
		}

		// process head tracker event
		if jobID.Str == "head-tracker" {
			height := gjson.Get(message, "result.data.value.block.header.height")
			c.subs[jobID.Str].Channel <- Events{Block: height.Uint()}
			continue
		}

		if jobID.Str == string(ABCI) {
			var res abci.ResponseQuery
			// code == 0 when no error is encountered
			result := gjson.Get(message, "result.response").Raw
			err := json.Unmarshal([]byte(result), &res)
			if err != nil {
				err := fmt.Sprintf("Couldn't decode result string: %s", err)
				c.queryCh <- RpcMessage{Error: err}
				continue
			}
			if res.Code != 0 {
				// if there's an error it's the `log` field
				c.queryCh <- RpcMessage{Error: res.Log}
				continue
			}
			c.queryCh <- RpcMessage{Data: res.Value}
			continue
		}

		if jobID.Str == string(TX) {
			errorMessage := gjson.Get(message, "error")
			if errorMessage.Raw != "" {
				c.queryCh <- RpcMessage{Error: errorMessage.Raw}
				continue
			}

			response := gjson.Get(message, "result")

			c.queryCh <- RpcMessage{Data: []byte(response.Raw)}
			continue
		}

		if strings.Contains(jobID.Str, "broadcast_tx") {
			errorMessage := gjson.Get(message, "error")
			if errorMessage.Raw != "" {
				c.txCh <- RpcMessage{Error: errorMessage.Raw}
				continue
			}

			// for BrodcastBlock need to check the result of deliver_tx
			if jobID.Str == string(BroadcastBlock) {
				responseCode := gjson.Get(message, "result.deliver_tx.code").Int()
				if responseCode != 0 {
					responseLog := gjson.Get(message, "result.log")
					c.queryCh <- RpcMessage{Error: responseLog.Raw}
				}
			}

			responseCode := gjson.Get(message, "result.code").Int()
			if responseCode != 0 {
				responseLog := gjson.Get(message, "result.log")
				c.queryCh <- RpcMessage{Error: responseLog.Raw}
			}

			response := gjson.Get(message, "result")

			c.txCh <- RpcMessage{Data: []byte(response.Raw)}
			continue
		}

		// parse events but skip if nothing is found (happens for sub/unsub responses)
		eventsRaw := gjson.Get(message, "result.data.value.TxResult.result.events")
		if !eventsRaw.Exists() {
			continue
		}

		// parse events as tendermint events
		var events []types.Event
		if err := json.Unmarshal([]byte(eventsRaw.Raw), &events); err != nil {
			c.Log.Errorf("[client/listen/unmarshal]: %s %s", err, eventsRaw)
		}

		// parse data into a standard format based on events
		parsedEvents := parseEvents(events)

		// send data to job specific channel
		c.subs[jobID.Str].Channel <- Events{Events: parsedEvents}
	}
}

func parseEvents(events []types.Event) (output []string) {
	// PLACEHOLDER - just returns an array of event names (no data)
	// example: https://github.com/smartcontractkit/external-initiator/blob/84cec9a579530db29ae3ca2489819c3e54d4960c/blockchain/terra/terra.go#L131
	// TODO: implement OCR specific event filtering
	for _, event := range events {
		if strings.HasPrefix(event.Type, "wasm-") {
			output = append(output, event.Type)
		}
	}
	return
}

func (c *Client) Unsubscribe(ctx context.Context, jobID string) error {
	// check if jobID exists
	if _, ok := c.subs[jobID]; !ok {
		return fmt.Errorf("[%s] Cannot unsubscribe. Job does not exist", jobID)
	}

	// create unsubscribe message from the stored subscribe message
	payload := c.subs[jobID].Payload
	payload.Method = "unsubscribe"

	// set up chanel to catch unsub response
	subUnsub := make(chan Events)
	c.subUnsub[jobID] = subUnsub

	// send unsub payload
	if err := c.ws.conn.WriteJSON(payload); err != nil {
		return err
	}

	// wait for successful unsubscription message
	msg := <-subUnsub
	if msg.SubErr != "" {
		return fmt.Errorf("[%s] WS error: %s", jobID, msg.SubErr)
	}

	// remove saved job
	delete(c.subs, jobID)
	c.Log.Infof("[%s] Unsubscribe successful", jobID)
	return nil
}

// Close websocket connection
func (c Client) Close() error {
	if err := c.Unsubscribe(context.TODO(), "head-tracker"); err != nil {
		c.Log.Errorf("[head-tracker] %s", err)
	}

	// trigger close channel to trigger stop related services
	close(c.close)

	// return if ws client has not been initialized
	if !c.wsStarted {
		return nil
	}

	c.Log.Infof("Closing websocket connection to %s", c.wsURL)
	return c.ws.conn.Close()
}
