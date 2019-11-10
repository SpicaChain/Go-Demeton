// Copyright (C) 2018 go-demeton authors
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

package dvm

import "C"

import (
	"crypto/md5"

	"github.com/cyber-demeton/go-demeton/core"
	"github.com/cyber-demeton/go-demeton/crypto/hash"
	"github.com/cyber-demeton/go-demeton/crypto/keystore"
	"github.com/cyber-demeton/go-demeton/util/byteutils"
	"github.com/cyber-demeton/go-demeton/util/logging"
	"github.com/sirupsen/logrus"
)

// Sha256Func ..
//export Sha256Func
func Sha256Func(data *C.char, gasCnt *C.size_t) *C.char {
	s := C.GoString(data)
	*gasCnt = C.size_t(len(s) + CryptoSha256GasBase)

	r := hash.Sha256([]byte(s))
	return C.CString(byteutils.Hex(r))
}

// Sha3256Func ..
//export Sha3256Func
func Sha3256Func(data *C.char, gasCnt *C.size_t) *C.char {
	s := C.GoString(data)
	*gasCnt = C.size_t(len(s) + CryptoSha3256GasBase)

	r := hash.Sha3256([]byte(s))
	return C.CString(byteutils.Hex(r))
}

// Ripemd160Func ..
//export Ripemd160Func
func Ripemd160Func(data *C.char, gasCnt *C.size_t) *C.char {
	s := C.GoString(data)
	*gasCnt = C.size_t(len(s) + CryptoRipemd160GasBase)

	r := hash.Ripemd160([]byte(s))
	return C.CString(byteutils.Hex(r))
}

// RecoverAddressFunc ..
//export RecoverAddressFunc
func RecoverAddressFunc(alg int, data, sign *C.char, gasCnt *C.size_t) *C.char {
	d := C.GoString(data)
	s := C.GoString(sign)

	*gasCnt = C.size_t(CryptoRecoverAddressGasBase)

	plain, err := byteutils.FromHex(d)
	if err != nil {
		logging.VLog().WithFields(logrus.Fields{
			"hash": d,
			"sign": s,
			"alg":  alg,
			"err":  err,
		}).Debug("convert hash to byte array error.")
		return nil
	}
	cipher, err := byteutils.FromHex(s)
	if err != nil {
		logging.VLog().WithFields(logrus.Fields{
			"data": d,
			"sign": s,
			"alg":  alg,
			"err":  err,
		}).Debug("convert sign to byte array error.")
		return nil
	}
	addr, err := core.RecoverSignerFromSignature(keystore.Algorithm(alg), plain, cipher)
	if err != nil {
		logging.VLog().WithFields(logrus.Fields{
			"data": d,
			"sign": s,
			"alg":  alg,
			"err":  err,
		}).Debug("recover address error.")
		return nil
	}

	return C.CString(addr.String())
}

// Md5Func ..
//export Md5Func
func Md5Func(data *C.char, gasCnt *C.size_t) *C.char {
	s := C.GoString(data)
	*gasCnt = C.size_t(len(s) + CryptoMd5GasBase)

	r := md5.Sum([]byte(s))
	return C.CString(byteutils.Hex(r[:]))
}

// Base64Func ..
//export Base64Func
func Base64Func(data *C.char, gasCnt *C.size_t) *C.char {
	s := C.GoString(data)
	*gasCnt = C.size_t(len(s) + CryptoBase64GasBase)

	r := hash.Base64Encode([]byte(s))
	return C.CString(string(r))
}
