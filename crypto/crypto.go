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

package crypto

import (
	"errors"

	"github.com/cyber-demeton/go-demeton/crypto/keystore"
	"github.com/cyber-demeton/go-demeton/crypto/keystore/secp256k1"
)

var (
	// ErrAlgorithmInvalid invalid Algorithm for sign.
	ErrAlgorithmInvalid = errors.New("invalid Algorithm")
)

// NewPrivateKey generate a privatekey with Algorithm
func NewPrivateKey(alg keystore.Algorithm, data []byte) (keystore.PrivateKey, error) {
	switch alg {
	case keystore.SECP256K1:
		var (
			priv *secp256k1.PrivateKey
			err  error
		)
		if len(data) == 0 {
			priv = secp256k1.GeneratePrivateKey()
		} else {
			priv = new(secp256k1.PrivateKey)
			err = priv.Decode(data)
		}
		if err != nil {
			return nil, err
		}
		return priv, nil
	default:
		return nil, ErrAlgorithmInvalid
	}
}

// NewSignature returns a specific signature with the algorithm
func NewSignature(alg keystore.Algorithm) (keystore.Signature, error) {
	switch alg {
	case keystore.SECP256K1:
		return new(secp256k1.Signature), nil
	default:
		return nil, ErrAlgorithmInvalid
	}
}

// CheckAlgorithm check if support the input Algorithm
func CheckAlgorithm(alg keystore.Algorithm) error {
	switch alg {
	case keystore.SECP256K1:
		return nil
	default:
		return ErrAlgorithmInvalid
	}
}
