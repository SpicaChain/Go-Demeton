package neblet

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof" // Register some standard stuff

	"github.com/sirupsen/logrus"

	"github.com/cyber-demeton/go-demeton/cmd/console"

	"net"

	"github.com/cyber-demeton/go-demeton/account"
	"github.com/cyber-demeton/go-demeton/consensus/dpos"
	"github.com/cyber-demeton/go-demeton/core"
	"github.com/cyber-demeton/go-demeton/core/pb"
	"github.com/cyber-demeton/go-demeton/metrics"
	"github.com/cyber-demeton/go-demeton/neblet/pb"
	nebnet "github.com/cyber-demeton/go-demeton/net"
	"github.com/cyber-demeton/go-demeton/dvm"
	"github.com/cyber-demeton/go-demeton/rpc"
	"github.com/cyber-demeton/go-demeton/storage"
	nsync "github.com/cyber-demeton/go-demeton/sync"
	"github.com/cyber-demeton/go-demeton/util/logging"
	m "github.com/rcrowley/go-metrics"
)

var (
	// ErrNebletAlreadyRunning throws when the neblet is already running.
	ErrNebletAlreadyRunning = errors.New("neblet is already running")

	// ErrIncompatibleStorageSchemeVersion throws when the storage schema has been changed
	ErrIncompatibleStorageSchemeVersion = errors.New("incompatible storage schema version, pls migrate your storage")
)

var (
	metricsNebstartGauge = m.GetOrRegisterGauge("deb.start", nil)
)

// Deblet manages ldife cycle of blockchain services.
type Deblet struct {
	config *nebletpb.Config

	genesis *corepb.Genesis

	accountManager *account.Manager

	netService nebnet.Service

	consensus core.Consensus

	storage storage.Storage

	blockChain *core.BlockChain

	syncService *nsync.Service

	rpcServer rpc.GRPCServer

	lock sync.RWMutex

	eventEmitter *core.EventEmitter

	dvm core.DVM

	running bool
}

// New returns a new neblet.
func New(config *nebletpb.Config) (*Deblet, error) {
	//var err error
	n := &Deblet{config: config}

	// try enable profile.
	n.TryStartProfiling()

	if chain := config.GetChain(); chain == nil {
		logging.CLog().Error("Failed to find chain config in config file")
		return nil, ErrConfigShouldHasChain
	}

	var err error
	n.genesis, err = core.LoadGenesisConf(config.Chain.Genesis)
	if err != nil {
		logging.CLog().Error("Failed to load genesis config")
		return nil, err
	}

	if config.Chain.Dynasty != "" {
		core.LoadDynastyConf(config.Chain.Dynasty, n.genesis)
	}

	am, err := account.NewManager(n)
	if err != nil {
		return nil, err
	}
	n.accountManager = am

	// init random seed.
	rand.Seed(time.Now().UTC().UnixNano())

	return n, nil
}

// Setup setup neblet
func (n *Deblet) Setup() {
	var err error
	logging.CLog().Info("Setuping Deblet...")

	// storage
	// n.storage, err = storage.NewDiskStorage(n.config.Chain.Datadir)
	// n.storage, err = storage.NewMemoryStorage()
	n.storage, err = storage.NewRocksStorage(n.config.Chain.Datadir)
	if err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"dir": n.config.Chain.Datadir,
			"err": err,
		}).Fatal("Failed to open disk storage.")
	}

	// net
	n.netService, err = nebnet.NewNebService(n)
	if err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to setup net service.")
	}

	// dvm
	n.dvm = dvm.NewNebulasVM()
	if err = n.dvm.CheckV8Run(); err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to setup V8.")
	}
	// core
	n.eventEmitter = core.NewEventEmitter(40960)
	n.consensus = dpos.NewDpos()
	n.blockChain, err = core.NewBlockChain(n)
	if err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to setup blockchain.")
	}

	// consensus
	if err := n.consensus.Setup(n); err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to setup consensus.")
	}
	if err := n.blockChain.Setup(n); err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to setup blockchain.")
	}

	// sync
	n.syncService = nsync.NewService(n.blockChain, n.netService)
	n.blockChain.SetSyncService(n.syncService)

	// rpc
	n.rpcServer = rpc.NewServer(n)

	logging.CLog().Info("Setuped Deblet.")
}

// StartPprof start pprof http listen
func (n *Deblet) StartPprof(listen string) error {
	if len(listen) > 0 {
		conn, err := net.DialTimeout("tcp", listen, time.Second*1)
		if err == nil {
			logging.CLog().WithFields(logrus.Fields{
				"listen": listen,
				"err":    err,
			}).Error("Failed to start pprof")
			conn.Close()
			return err
		}

		go func() {
			logging.CLog().WithFields(logrus.Fields{
				"listen": listen,
			}).Info("Starting pprof...")
			http.ListenAndServe(listen, nil)
		}()
	}
	return nil
}

// Start starts the services of the neblet.
func (n *Deblet) Start() {
	n.lock.Lock()
	defer n.lock.Unlock()

	logging.CLog().Info("Starting Deblet...")

	if n.running {
		logging.CLog().WithFields(logrus.Fields{
			"err": "neblet is already running",
		}).Fatal("Failed to start neblet.")
	}
	n.running = true

	if n.config.Stats.EnableMetrics {
		metrics.Start(n)
	}

	if err := n.netService.Start(); err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to start net service.")
	}

	if err := n.rpcServer.Start(); err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to start api server.")
	}

	if err := n.rpcServer.RunGateway(); err != nil {
		logging.CLog().WithFields(logrus.Fields{
			"err": err,
		}).Fatal("Failed to start api gateway.")
	}

	n.blockChain.Start()
	n.blockChain.BlockPool().Start()
	n.blockChain.TransactionPool().Start()
	n.eventEmitter.Start()
	n.syncService.Start()

	// start consensus
	chainConf := n.config.Chain
	if chainConf.StartMine {
		n.consensus.Start()
		if chainConf.EnableRemoteSignServer == false {
			passphrase := chainConf.Passphrase
			if len(passphrase) == 0 {
				fmt.Println("***********************************************")
				fmt.Println("miner address:" + n.config.Chain.Miner)
				prompt := console.Stdin
				passphrase, _ = prompt.PromptPassphrase("Enter the miner's passphrase:")
				fmt.Println("***********************************************")
			}
		}
		err := n.consensus.EnableMining(chainConf.Passphrase)
		if err != nil {
			logging.CLog().WithFields(logrus.Fields{
				"err": err,
			}).Fatal("Failed to enable mining.")
		}
	}

	// first sync
	if len(n.Config().Network.Seed) > 0 {
		n.blockChain.StartActiveSync()
	} else {
		if chainConf.StartMine {
			n.Consensus().ResumeMining()
		}
	}

	metricsNebstartGauge.Update(1)

	logging.CLog().Info("Started Deblet.")
}

// Stop stops the services of the neblet.
func (n *Deblet) Stop() {
	n.lock.Lock()
	defer n.lock.Unlock()

	logging.CLog().Info("Stopping Deblet...")

	// try Stop Profiling.
	n.TryStopProfiling()

	if n.config.Chain.StartMine && n.consensus != nil {
		n.consensus.Stop()
		n.consensus = nil
	}

	if n.syncService != nil {
		n.syncService.Stop()
		n.syncService = nil
	}

	if n.eventEmitter != nil {
		n.eventEmitter.Stop()
		n.eventEmitter = nil
	}

	if n.blockChain != nil {
		n.blockChain.TransactionPool().Stop()
		n.blockChain.BlockPool().Stop()
		n.blockChain.Stop()
		n.blockChain = nil
	}

	if n.rpcServer != nil {
		n.rpcServer.Stop()
		n.rpcServer = nil
	}

	if n.netService != nil {
		n.netService.Stop()
		n.netService = nil
	}

	if n.config.Stats.EnableMetrics {
		metrics.Stop()
	}

	n.accountManager = nil

	n.running = false

	logging.CLog().Info("Stopped Deblet.")
}

// SetGenesis set genesis conf
func (n *Deblet) SetGenesis(g *corepb.Genesis) {
	n.genesis = g
}

// Genesis returns genesis conf.
func (n *Deblet) Genesis() *corepb.Genesis {
	return n.genesis
}

// Config returns neblet configuration.
func (n *Deblet) Config() *nebletpb.Config {
	return n.config
}

// Storage returns storage reference.
func (n *Deblet) Storage() storage.Storage {
	return n.storage
}

// BlockChain returns block chain reference.
func (n *Deblet) BlockChain() *core.BlockChain {
	return n.blockChain
}

// EventEmitter returns eventEmitter reference.
func (n *Deblet) EventEmitter() *core.EventEmitter {
	return n.eventEmitter
}

// AccountManager returns account manager reference.
func (n *Deblet) AccountManager() core.AccountManager {
	return n.accountManager
}

// NetService returns p2p manager reference.
func (n *Deblet) NetService() nebnet.Service {
	return n.netService
}

// Consensus returns consensus reference.
func (n *Deblet) Consensus() core.Consensus {
	return n.consensus
}

// SyncService return sync service
func (n *Deblet) SyncService() *nsync.Service {
	return n.syncService
}

// IsActiveSyncing return if the deb is syncing blocks
func (n *Deblet) IsActiveSyncing() bool {
	if n.syncService == nil {
		return false
	}
	return n.syncService.IsActiveSyncing()
}

// Dvm return dvm engine
func (n *Deblet) Dvm() core.DVM {
	return n.dvm
}

// TryStartProfiling try start pprof
func (n *Deblet) TryStartProfiling() {
	if n.config.App == nil {
		logging.CLog().Error("Failed to find app config in config file")
		return
	}
	if n.config.App.Pprof == nil {
		logging.CLog().Error("Failed to find app.pprof config in config file")
		return
	}

	cpuProfile := n.config.App.Pprof.Cpuprofile
	if len(cpuProfile) > 0 {
		f, err := os.Create(cpuProfile)
		if err != nil {
			logging.CLog().WithFields(logrus.Fields{
				"err": err,
			}).Fatalf("Failed to create CPU profile %s", cpuProfile)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			logging.CLog().WithFields(logrus.Fields{
				"err": err,
			}).Fatalf("Failed to start CPU profile")
		}
	}
}

// TryStopProfiling try stop pprof
func (n *Deblet) TryStopProfiling() {
	if n.config.App.Pprof == nil {
		return
	}

	memProfile := n.config.App.Pprof.Memprofile
	if len(memProfile) > 0 {
		f, err := os.Create(memProfile)
		if err != nil {
			logging.CLog().WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Failed to create memory profile %s", memProfile)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			logging.CLog().WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Failed to write memory profile")
		}
		f.Close()
	}

	cpuProfile := n.config.App.Pprof.Cpuprofile
	if len(cpuProfile) > 0 {
		pprof.StopCPUProfile()
	}
}
