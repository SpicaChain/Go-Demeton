// Copyright (C) 2017 go-demeton authors
//
// This file is part of the go-demeton library.
//
// the go-demeton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-demeton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-demeton library.  If not, see <http://www.gnu.org/licenses/>.
//

package rpc

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/cyber-demeton/go-demeton/core"
	"github.com/cyber-demeton/go-demeton/crypto/keystore"
	"github.com/cyber-demeton/go-demeton/net"
	"github.com/cyber-demeton/go-demeton/rpc/pb"
	"golang.org/x/net/context"
)

// AdminService implements the RPC admin service interface.
type AdminService struct {
	server GRPCServer
}

// Accounts is the RPC API handler.
func (s *AdminService) Accounts(ctx context.Context, req *rpcpb.NonParamsRequest) (*rpcpb.AccountsResponse, error) {

	deb := s.server.Deblet()
	accs := deb.AccountManager().Accounts()

	resp := new(rpcpb.AccountsResponse)
	addrs := make([]string, len(accs))
	for index, addr := range accs {
		addrs[index] = addr.String()
	}
	resp.Addresses = addrs
	return resp, nil
}

// NewAccount generate a new address with passphrase
func (s *AdminService) NewAccount(ctx context.Context, req *rpcpb.NewAccountRequest) (*rpcpb.NewAccountResponse, error) {

	deb := s.server.Deblet()
	addr, err := deb.AccountManager().NewAccount([]byte(req.Passphrase))
	if err != nil {
		return nil, err
	}
	return &rpcpb.NewAccountResponse{Address: addr.String()}, nil
}

// UnlockAccount unlock address with the passphrase
func (s *AdminService) UnlockAccount(ctx context.Context, req *rpcpb.UnlockAccountRequest) (*rpcpb.UnlockAccountResponse, error) {
	deb := s.server.Deblet()

	addr, err := core.AddressParse(req.Address)
	if err != nil {
		metricsUnlockFailed.Mark(1)
		return nil, err
	}

	duration := time.Duration(req.Duration)
	if duration == 0 {
		duration = keystore.DefaultUnlockDuration
	}
	err = deb.AccountManager().Unlock(addr, []byte(req.Passphrase), duration)
	if err != nil {
		metricsUnlockFailed.Mark(1)
		return nil, err
	}

	metricsUnlockSuccess.Mark(1)
	return &rpcpb.UnlockAccountResponse{Result: true}, nil
}

// LockAccount lock address
func (s *AdminService) LockAccount(ctx context.Context, req *rpcpb.LockAccountRequest) (*rpcpb.LockAccountResponse, error) {
	deb := s.server.Deblet()

	addr, err := core.AddressParse(req.Address)
	if err != nil {
		return nil, err
	}

	err = deb.AccountManager().Lock(addr)
	if err != nil {
		return nil, err
	}

	return &rpcpb.LockAccountResponse{Result: true}, nil
}

// SendTransaction is the RPC API handler.
func (s *AdminService) SendTransaction(ctx context.Context, req *rpcpb.TransactionRequest) (*rpcpb.SendTransactionResponse, error) {
	deb := s.server.Deblet()
	tx, err := parseTransaction(deb, req)
	if err != nil {
		metricsSendTxFailed.Mark(1)
		return nil, err
	}
	if err := deb.AccountManager().SignTransaction(tx.From(), tx); err != nil {
		metricsSendTxFailed.Mark(1)
		return nil, err
	}

	return handleTransactionResponse(deb, tx)
}

// SignHash is the RPC API handler.
func (s *AdminService) SignHash(ctx context.Context, req *rpcpb.SignHashRequest) (*rpcpb.SignHashResponse, error) {
	deb := s.server.Deblet()

	hash := req.Hash
	addr, err := core.AddressParse(req.Address)
	if err != nil {
		return nil, err
	}
	alg := keystore.Algorithm(req.Alg)

	data, err := deb.AccountManager().SignHash(addr, hash, alg)
	if err != nil {
		return nil, err
	}

	return &rpcpb.SignHashResponse{Data: data}, nil
}

// GenerateRandomSeed generate block's rand info
func (s *AdminService) GenerateRandomSeed(ctx context.Context, req *rpcpb.GenerateRandomSeedRequest) (*rpcpb.GenerateRandomSeedResponse, error) {
	deb := s.server.Deblet()
	addr, err := core.AddressParse(req.Address)
	if err != nil {
		return nil, err
	}

	vrfSeed, vrfProof, err := deb.AccountManager().GenerateRandomSeed(addr, req.AncestorHash, req.ParentSeed)
	if err != nil {
		return nil, err
	}

	return &rpcpb.GenerateRandomSeedResponse{
		VrfSeed:  vrfSeed,
		VrfProof: vrfProof,
	}, nil
}

// SignTransactionWithPassphrase sign transaction with the from addr passphrase
func (s *AdminService) SignTransactionWithPassphrase(ctx context.Context, req *rpcpb.SignTransactionPassphraseRequest) (*rpcpb.SignTransactionPassphraseResponse, error) {

	deb := s.server.Deblet()
	tx, err := parseTransaction(deb, req.Transaction)
	if err != nil {
		metricsSignTxFailed.Mark(1)
		return nil, err
	}
	if err := deb.AccountManager().SignTransactionWithPassphrase(tx.From(), tx, []byte(req.Passphrase)); err != nil {
		metricsSignTxFailed.Mark(1)
		return nil, err
	}
	pbMsg, err := tx.ToProto()
	if err != nil {
		metricsSignTxFailed.Mark(1)
		return nil, err
	}
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		metricsSignTxFailed.Mark(1)
		return nil, err
	}

	metricsSignTxSuccess.Mark(1)
	return &rpcpb.SignTransactionPassphraseResponse{Data: data}, nil
}

// SendTransactionWithPassphrase send transaction with the from addr passphrase
func (s *AdminService) SendTransactionWithPassphrase(ctx context.Context, req *rpcpb.SendTransactionPassphraseRequest) (*rpcpb.SendTransactionResponse, error) {

	deb := s.server.Deblet()
	tx, err := parseTransaction(deb, req.Transaction)
	if err != nil {
		return nil, err
	}
	if err := deb.AccountManager().SignTransactionWithPassphrase(tx.From(), tx, []byte(req.Passphrase)); err != nil {
		return nil, err
	}

	return handleTransactionResponse(deb, tx)
}

// StartPprof start pprof
func (s *AdminService) StartPprof(ctx context.Context, req *rpcpb.PprofRequest) (*rpcpb.PprofResponse, error) {
	deb := s.server.Deblet()

	if err := deb.StartPprof(req.Listen); err != nil {
		return nil, err
	}
	return &rpcpb.PprofResponse{Result: true}, nil
}

// GetConfig is the RPC API handler.
func (s *AdminService) GetConfig(ctx context.Context, req *rpcpb.NonParamsRequest) (*rpcpb.GetConfigResponse, error) {

	deb := s.server.Deblet()

	resp := &rpcpb.GetConfigResponse{}
	resp.Config = deb.Config()
	resp.Config.Chain.Passphrase = string("")
	return resp, nil
}

// NodeInfo is the RPC API handler
/*
限制来自同一个ip的节点连接请求的数量 （例如来自同一个ip的节点连接不能超过10，如果当前连接的列表中来自同一个ip的节点数量为10，则拒绝所有后面来自该ip的节点的连接请求）
主动发起连接时判断目标节点的ip是否在已连接的列表中，如果已经存在，则不建立该连接
路由同步增加相应的策略
一个桶的地址不能包含两个以上节点相同的 /24 ip地址块
整个路由表不能包含十个以上节点相同的 /24 ip地址块
更改路由同步的算法，路由同步时候不再同步离目标节点最近的那些节点
*/
func (s *AdminService) NodeInfo(ctx context.Context, req *rpcpb.NonParamsRequest) (*rpcpb.NodeInfoResponse, error) {

	deb := s.server.Deblet()

	resp := &rpcpb.NodeInfoResponse{}
	node := deb.NetService().Node()
	resp.Id = node.ID() // FIXME: @leon check eclipse attack
	resp.ChainId = node.Config().ChainID
	resp.BucketSize = int32(node.Config().Bucketsize)
	resp.PeerCount = uint32(node.PeersCount())
	resp.ProtocolVersion = net.NebProtocolID
	resp.Coinbase = deb.Config().Chain.Coinbase

	for k, v := range node.RouteTable().Peers() {
		routeTable := &rpcpb.RouteTable{}
		routeTable.Id = k.Pretty()
		routeTable.Address = make([]string, len(v))

		for i, addr := range v {
			routeTable.Address[i] = addr.String()
		}
		resp.RouteTable = append(resp.RouteTable, routeTable)
	}

	return resp, nil
}
