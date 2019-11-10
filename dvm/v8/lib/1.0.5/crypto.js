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

'use strict';

const HexStringRegex = /^[0-9a-fA-F]+$/;

var Crypto = function() {
    Object.defineProperty(this, "nativeCrypto", {
        configurable: false,
        enumerable: false,
        get: function(){
            return _native_crypto;
        }
    });
};

Crypto.prototype = {
 
    // case sensitive
    sha256: function(data) {
        if (typeof data !== "string") {
            throw new Error("input must be string");
        }
        // any string
        return this.nativeCrypto.sha256(data);
    },

    // case sensitive
    sha3256: function(data) {
        if (typeof data !== "string") {
            throw new Error("input must be string");
        }
        // any string
        return this.nativeCrypto.sha3256(data);
    },

    // case sensitive
    ripemd160: function(data) {
        if (typeof data !== "string") {
            throw new Error("input must be string");
        }
        // any string
        return this.nativeCrypto.ripemd160(data);
    },

    // case insensitive
    recoverAddress: function(alg, hash, sign) {
        if (!Number.isSafeInteger(alg) || alg < 0) {
            throw new Error("alg must be non-negative integer");
        }

        if (typeof hash !== "string" || !HexStringRegex.test(hash) 
            || typeof sign !== "string" || !HexStringRegex.test(sign)) {
            throw new Error("hash & sign must be hex string");
        }
        // alg: 1
        // hash: sha3256 hex string, 64 chars
        // sign: cipher hex string by private key, 130 chars
        return this.nativeCrypto.recoverAddress(alg, hash, sign);
    },

    // case sensitive
    md5: function(data) {
        if (typeof data !== "string") {
            throw new Error("input must be string");
        }
        // any string
        return this.nativeCrypto.md5(data);
    },

    // case sensitive
    base64: function(data) {
        if (typeof data !== "string") {
            throw new Error("input must be string");
        }
        // any string
        return this.nativeCrypto.base64(data);
    }
};

module.exports = new Crypto();