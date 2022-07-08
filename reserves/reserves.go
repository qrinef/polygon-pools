package reserves

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"polygon-pools/common/syncmap"
	"polygon-pools/pools"
)

type Reserves struct {
	newPools  *pools.Pools
	emulation emulation
	pools     *syncmap.Map[common.Address, pool]
}

type pool struct {
	token0    common.Address
	token1    common.Address
	reserve0  *big.Int
	reserve1  *big.Int
	lastBlock uint64
	lastIndex uint
}

type poolJSON struct {
	Pool     common.Address
	Token0   common.Address
	Token1   common.Address
	Reserve0 string
	Reserve1 string
}

func NewReserves(newPools *pools.Pools) *Reserves {
	return &Reserves{
		newPools: newPools,
		pools:    syncmap.New[common.Address, pool](),
	}
}

func (r *Reserves) Start() (err error) {
	r.setPools()
	r.setSubscribes()

	return r.setEmulation()
}

func (r *Reserves) GetReserves() ([]byte, error) {
	var _pools []poolJSON

	for _address, _pool := range r.pools.Items() {
		if _pool.reserve0 == nil || _pool.reserve1 == nil {
			continue
		}

		_pools = append(_pools, poolJSON{
			Pool:     _address,
			Token0:   _pool.token0,
			Token1:   _pool.token1,
			Reserve0: _pool.reserve0.String(),
			Reserve1: _pool.reserve1.String(),
		})
	}

	return json.MarshalIndent(_pools, "", " ")
}

func (r *Reserves) setPools() {
	for address, _pool := range r.newPools.GetPools() {
		r.pools.Store(address, pool{
			token0: _pool.Token0,
			token1: _pool.Token1,
		})
	}
}
