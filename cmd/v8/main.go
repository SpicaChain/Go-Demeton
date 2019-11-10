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

package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/cyber-demeton/go-demeton/core"
	"github.com/cyber-demeton/go-demeton/core/state"
	"github.com/cyber-demeton/go-demeton/dvm"
	"github.com/cyber-demeton/go-demeton/storage"
)

func main() {
	data, _ := ioutil.ReadFile(os.Args[1])

	mem, _ := storage.NewMemoryStorage()
	context, _ := state.NewWorldState(nil, mem)
	contract, _ := context.CreateContractAccount([]byte("account2"), nil, nil)

	ctx, err := dvm.NewContext(core.MockBlock(nil, 1), nil, contract, context)
	if err == nil {
		engine := dvm.NewV8Engine(ctx)
		result, err := engine.RunScriptSource(string(data), 0)

		log.Fatalf("Result is %s. Err is %s", result, err)

		time.Sleep(10 * time.Second)
		engine.Dispose()
	} else {
		os.Exit(1)
	}
}
