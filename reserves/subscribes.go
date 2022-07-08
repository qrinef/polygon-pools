package reserves

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

var endpoints = []string{
	"wss://ws-matic-mainnet.chainstacklabs.com",
	"wss://rpc-mainnet.matic.quiknode.pro",
}

func (r *Reserves) setSubscribes() {
	fmt.Println("> Get subscribes...")

	channel := make(chan types.Log, 1000)

	for _, endpoint := range endpoints {
		if err := r.subscribe(endpoint, channel); err != nil {
			fmt.Printf(">> Subscribe (ER): %v\n", endpoint)
		} else {
			fmt.Printf(">> Subscribe (OK): %v\n", endpoint)
		}
	}

	go r.handler(channel)
	fmt.Println()
}

func (r *Reserves) subscribe(endpoint string, channel chan types.Log) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		return err
	}

	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{{
			common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"),
		}},
	}

	if _, err = client.SubscribeFilterLogs(context.Background(), query, channel); err != nil {
		return err
	}

	return err
}

func (r *Reserves) handler(logs chan types.Log) {
	for {
		select {
		case log := <-logs:
			if p, ok := r.pools.Load(log.Address); ok {
				if log.BlockNumber > p.lastBlock || (log.BlockNumber == p.lastBlock && log.Index > p.lastIndex) {
					p.reserve0 = new(big.Int).SetBytes(log.Data[:32])
					p.reserve1 = new(big.Int).SetBytes(log.Data[32:])
					p.lastBlock = log.BlockNumber
					p.lastIndex = log.Index

					r.pools.Store(log.Address, p)
				}
			}
		}
	}
}
