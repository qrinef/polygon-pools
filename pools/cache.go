package pools

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (p *Pools) setCache() (err error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	p.cachePath = filepath.Join(dirname, "/.polygon-pools/pools.json")

	if err = os.MkdirAll(filepath.Dir(p.cachePath), os.ModePerm); os.IsPermission(err) {
		return err
	}

	return p.readCache()
}

func (p *Pools) readCache() (err error) {
	cache, err := os.OpenFile(p.cachePath, os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(cache).Decode(&p.pools); err != nil && err != io.EOF {
		return err
	}

	if len(p.pools.Pools) < 1 {
		p.pools.LastBlockNumber = "0x0"
	}

	return cache.Close()
}

func (p *Pools) writeCache() (err error) {
	cache, err := json.MarshalIndent(p.pools, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.cachePath, cache, 0644)
}
