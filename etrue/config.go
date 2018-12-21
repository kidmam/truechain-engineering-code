// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package etrue

import (
	"math/big"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethash "github.com/truechain/truechain-engineering-code/consensus/minerva"
	"github.com/truechain/truechain-engineering-code/core"
	"github.com/truechain/truechain-engineering-code/core/snailchain"
	"github.com/truechain/truechain-engineering-code/etrue/downloader"
	"github.com/truechain/truechain-engineering-code/etrue/gasprice"
	"github.com/truechain/truechain-engineering-code/params"
)

// DefaultConfig contains default settings for use on the Truechain main net.
var DefaultConfig = Config{
	SyncMode: downloader.SnapShotSync,
	Ethash: ethash.Config{
		CacheDir:       "minerva",
		CachesInMem:    2,
		CachesOnDisk:   3,
		DatasetsInMem:  1,
		DatasetsOnDisk: 2,
	},
	NetworkId:     1,
	LightPeers:    100,
	DatabaseCache: 768,
	TrieCache:     256,
	TrieTimeout:   60 * time.Minute,
	GasPrice:      big.NewInt(18 * params.Shannon),

	TxPool:    core.DefaultTxPoolConfig,
	SnailPool: snailchain.DefaultSnailPoolConfig,
	GPO: gasprice.Config{
		Blocks:     20,
		Percentile: 60,
	},
	MinerThreads: 2,
	Port:         30310,
	StandbyPort:  30311,
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
	if runtime.GOOS == "windows" {
		DefaultConfig.Ethash.DatasetDir = filepath.Join(home, "AppData", "Minerva")
	} else {
		DefaultConfig.Ethash.DatasetDir = filepath.Join(home, ".minerva")
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Truechain main net block is used.
	Genesis *core.Genesis
	// FastGenesis  *fastchain.Genesis
	// SnailGenesis *snailchain.Genesis

	// Protocol options
	NetworkId    uint64 // Network ID to use for selecting peers to connect to
	SyncMode     downloader.SyncMode
	NoPruning    bool
	DeletedState bool

	// Light client options
	LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	// election options

	EnableElection bool `toml:",omitempty"`
	// CommitteeKey is the ECDSA private key for committee member.
	// If this filed is empty, can't be a committee member.
	CommitteeKey []byte `toml:",omitempty"`

	PrivateKey *ecdsa.PrivateKey `toml:"-"`

	// Host is the host interface on which to start the pbft server. If this
	// field is empty, can't be a committee member.
	Host string `toml:",omitempty"`

	// Port is the TCP port number on which to start the pbft server.
	Port int `toml:",omitempty"`

	// StandByPort is the TCP port number on which to start the pbft server.
	StandbyPort int `toml:",omitempty"`

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	TrieCache          int
	TrieTimeout        time.Duration

	// Mining-related options
	Etherbase    common.Address `toml:",omitempty"`
	MinerThreads int            `toml:",omitempty"`
	ExtraData    []byte         `toml:",omitempty"`
	GasPrice     *big.Int

	// Ethash options
	Ethash ethash.Config

	// Transaction pool options
	TxPool core.TxPoolConfig

	//fruit pool options
	SnailPool snailchain.SnailPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`

	// true indicate singlenode start
	NodeType bool `toml:",omitempty"`

	//true indicate only mine fruit
	MineFruit bool `toml:",omitempty"`

	//start for old pbft server
	OldTbft bool `toml:",omitempty"`
}

func (c *Config) GetNodeType() bool {
	return c.NodeType
}

type configMarshaling struct {
	ExtraData hexutil.Bytes
}
