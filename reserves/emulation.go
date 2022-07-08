package reserves

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strings"
	"time"
)

//go:embed contracts/emulation/abi.json
var emulationABI string

//go:embed contracts/emulation/runtime.hex
var emulationRuntime string

type emulation struct {
	clientRPC *rpc.Client
	abi       abi.ABI
	runtime   []byte
}

type reserves struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

func (r *Reserves) setEmulation() (err error) {
	fmt.Println("> Get pools reserves...")

	if err = r.setContract(); err != nil {
		return err
	}

	if err = r.runChunks(); err != nil {
		return err
	}

	fmt.Println()
	return err
}

func (r *Reserves) setContract() (err error) {
	r.emulation.clientRPC, err = rpc.Dial("https://polygon-rpc.com")
	if err != nil {
		return err
	}

	r.emulation.abi, err = abi.JSON(strings.NewReader(emulationABI))
	if err != nil {
		return err
	}

	r.emulation.runtime = common.FromHex(strings.TrimSpace(emulationRuntime))

	return err
}

func (r *Reserves) getPoolsAddresses() (pools []common.Address) {
	for address, _ := range r.pools.Items() {
		pools = append(pools, address)
	}

	return pools
}

func (r *Reserves) runChunks() (err error) {
	pools := r.getPoolsAddresses()

	length := len(pools)
	chunkSize := 1000

	for i := 0; i < length; i += chunkSize {
		end := i + chunkSize
		if end > length {
			end = length
		}

		_reserves, _err := r.executeChunk(pools[i:end])
		if _err != nil {
			return _err
		}

		if _err = r.handlerReserves(_reserves); _err != nil {
			return _err
		}

		fmt.Printf(">> Get pools reserves %v of %v\n", end, length)
	}

	return err
}

func (r *Reserves) executeChunk(pools []common.Address) (_reserves map[common.Address]reserves, err error) {
	_reserves = make(map[common.Address]reserves)
	attempts := 0

	_err := r.runEmulation(pools, _reserves)
	for _err != nil {
		if attempts > 10 {
			return _reserves, errors.New("exceeded attempts receive reserves")
		}

		attempts++
		time.Sleep(time.Second)
		_err = r.runEmulation(pools, _reserves)
	}

	return _reserves, err
}

func (r *Reserves) handlerReserves(_reserves map[common.Address]reserves) (err error) {
	for _poolAddress, _poolReserves := range _reserves {
		_pool, ok := r.pools.Load(_poolAddress)
		if !ok {
			return errors.New("pool address not found")
		}

		if _pool.reserve0 == nil && _pool.reserve1 == nil {
			_pool.reserve0 = _poolReserves.Reserve0
			_pool.reserve1 = _poolReserves.Reserve1

			r.pools.Store(_poolAddress, _pool)
		}
	}

	return err
}

func (r *Reserves) runEmulation(pool []common.Address, _reserves map[common.Address]reserves) (err error) {
	to := common.HexToAddress("0x2f347aDb46dC79C85AC05f92DBb020ee7C6b66A6")

	input, err := r.emulation.abi.Pack("start", pool)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"to":   to,
		"data": hexutil.Bytes(input),
	}

	override := map[*common.Address]struct{ Code hexutil.Bytes }{
		&to: {Code: r.emulation.runtime},
	}

	var call interface{}
	if err = r.emulation.clientRPC.Call(&call, "eth_call", data, "latest", override); err != nil {
		return err
	}

	var result struct {
		Reserves []reserves
	}
	if err = r.emulation.abi.UnpackIntoInterface(&result, "start", common.FromHex(call.(string))); err != nil {
		return err
	}

	for i, _poolReserves := range result.Reserves {
		_reserves[pool[i]] = reserves{
			Reserve0: _poolReserves.Reserve0,
			Reserve1: _poolReserves.Reserve1,
		}
	}

	return err
}
