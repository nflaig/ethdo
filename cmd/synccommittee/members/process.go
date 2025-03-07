// Copyright © 2021 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package members

import (
	"context"
	"fmt"
	"strings"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

func process(ctx context.Context, data *dataIn) (*dataOut, error) {
	if data == nil {
		return nil, errors.New("no data")
	}

	epoch, err := calculateEpoch(ctx, data)
	if err != nil {
		return nil, err
	}

	syncCommittee, err := data.eth2Client.(eth2client.SyncCommitteesProvider).SyncCommitteeAtEpoch(ctx, "head", epoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain sync committee information")
	}

	if syncCommittee == nil {
		return nil, errors.New("no sync committee returned")
	}

	results := &dataOut{
		debug:      data.debug,
		quiet:      data.quiet,
		verbose:    data.verbose,
		validators: syncCommittee.Validators,
	}

	return results, nil
}

func calculateEpoch(_ context.Context, data *dataIn) (phase0.Epoch, error) {
	var epoch phase0.Epoch
	if data.epoch != -1 {
		epoch = phase0.Epoch(data.epoch)
	} else {
		switch strings.ToLower(data.period) {
		case "", "current":
			epoch = data.chainTime.CurrentEpoch()
		case "next":
			period := data.chainTime.SlotToSyncCommitteePeriod(data.chainTime.CurrentSlot())
			nextPeriod := period + 1
			epoch = data.chainTime.FirstEpochOfSyncPeriod(nextPeriod)
		default:
			return 0, fmt.Errorf("period %s not known", data.period)
		}
	}

	if data.debug {
		fmt.Printf("epoch is %d\n", epoch)
	}

	return epoch, nil
}
