package pools

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type Pools struct {
	pools     pools
	cachePath string
}

type pools struct {
	LastBlockNumber string
	Pools           map[common.Address]Pool
}

type Pool struct {
	Factory common.Address
	Token0  common.Address
	Token1  common.Address
}

func NewPools() *Pools {
	return &Pools{
		pools: pools{
			Pools: map[common.Address]Pool{},
		},
	}
}

func (p *Pools) Start() (err error) {
	fmt.Println("> Get pools...")

	if err = p.setCache(); err != nil {
		return err
	}

	if err = p.getPools(); err != nil {
		return err
	}

	if err = p.writeCache(); err != nil {
		return err
	}

	return err
}

func (p *Pools) GetPools() (pools map[common.Address]Pool) {
	return p.pools.Pools
}
