// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"bytes"
	"strconv"
	"testing"

	ag_gofuzz "github.com/gagliardetto/gofuzz"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	ag_require "github.com/stretchr/testify/require"
)

func TestEncodeDecode_Transfer(t *testing.T) {
	fu := ag_gofuzz.New().NilChance(0)
	for i := 0; i < 1; i++ {
		t.Run("Transfer"+strconv.Itoa(i), func(t *testing.T) {
			{
				params := new(Transfer)
				fu.Fuzz(params)
				params.AccountMetaSlice = nil
				buf := new(bytes.Buffer)
				err := encodeT(*params, buf)
				ag_require.NoError(t, err)
				//
				got := new(Transfer)
				err = decodeT(got, buf.Bytes())
				got.AccountMetaSlice = nil
				ag_require.NoError(t, err)
				ag_require.Equal(t, params, got)
			}
		})
	}
}

func TestTransfer_AssertEquivalent(t *testing.T) {
	tx1, err := ag_solanago.NewTransaction([]ag_solanago.Instruction{
		NewTransferInstruction(
			1,
			ag_solanago.PK{},
			ag_solanago.PK{},
		).Build(),
	}, ag_solanago.Hash{})
	require.NoError(t, err)

	require.NoError(t, tx1.Message.AssertEquivalent(tx1.Message))

	tx2, err := ag_solanago.NewTransaction([]ag_solanago.Instruction{
		NewTransferInstruction(
			2,
			ag_solanago.PK{},
			ag_solanago.PK{},
		).Build(),
	}, ag_solanago.Hash{}, ag_solanago.TransactionPayer(ag_solanago.PK{}))

	require.Error(t, tx1.Message.AssertEquivalent(tx2.Message))
}
