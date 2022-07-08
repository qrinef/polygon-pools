package pools

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type logs struct {
	Result []struct {
		Address     common.Address
		Data        string
		Topics      []string
		BlockNumber string
	}
}

func (p *Pools) getPools() (err error) {
	totalNewPools := 0

	for {
		_logs, _fromBlock, _err := p.getLogs(p.pools.LastBlockNumber)
		if _err != nil {
			return _err
		}

		countNewPools := p.handlerLogs(_logs)
		totalNewPools += countNewPools

		if countNewPools < 1 {
			break
		}

		fmt.Printf(">> Get pools from block %v of ~3000000\n", _fromBlock)
	}

	fmt.Printf("> Total pools: %v (new: %v)\n\n", len(p.pools.Pools), totalNewPools)

	return err
}

func (p *Pools) getLogs(fromBlockHex string) (logs logs, fromBlock int64, err error) {
	fromBlock, err = strconv.ParseInt(fromBlockHex[2:], 16, 64)
	if err != nil {
		return logs, fromBlock, err
	}

	client := req.C().SetTimeout(30 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"apiKey":    "AVGK24UMQVX1E9RMEVARPBWW7JT8XJ2KZY", // public key, everything is fine
			"module":    "logs",
			"action":    "getLogs",
			"topic0":    "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9",
			"fromBlock": strconv.FormatInt(fromBlock, 10),
		}).
		SetRetryCount(5).
		SetRetryFixedInterval(1 * time.Second).
		SetResult(&logs).
		Get("https://api.polygonscan.com/api")

	if !resp.IsSuccess() {
		return logs, fromBlock, errors.Errorf("retrieving logs failed")
	}

	return logs, fromBlock, err
}

func (p *Pools) handlerLogs(logs logs) (countNewPools int) {
	for _, log := range logs.Result {
		pool := common.HexToAddress(log.Data[26:66])

		if _, ok := p.pools.Pools[pool]; !ok {
			p.pools.Pools[pool] = Pool{
				Factory: log.Address,
				Token0:  common.HexToAddress(log.Topics[1]),
				Token1:  common.HexToAddress(log.Topics[2]),
			}

			countNewPools++
		}
	}

	p.pools.LastBlockNumber = logs.Result[len(logs.Result)-1].BlockNumber

	return countNewPools
}
