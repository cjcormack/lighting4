// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package tracks

import (
	"github.com/satori/go.uuid"
)

var subscribers map[uuid.UUID]ValueChangeCallback

type ValuesChange struct {
	TrackChanged bool
	TrackState   *TrackState
}

type ValueChangeCallback func(change ValuesChange)

func init() {
	subscribers = make(map[uuid.UUID]ValueChangeCallback)
}

func Subscribe(callback ValueChangeCallback) func() {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	subscribers[id] = callback

	return func() {
		delete(subscribers, id)
	}
}

func notify(trackChanged bool, trackState *TrackState) {
	for _, callback := range subscribers {
		callback(ValuesChange {
			TrackChanged: trackChanged,
			TrackState:   trackState,
		})
	}
}
