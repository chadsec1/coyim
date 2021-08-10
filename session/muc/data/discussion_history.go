package data

import (
	"sync"
	"time"
)

// DelayedMessage contains the information of a delayed message
type DelayedMessage struct {
	Nickname  string
	Message   string
	Timestamp time.Time
}

func newDelayedMessage(nickname, message string, timestamp time.Time) *DelayedMessage {
	return &DelayedMessage{
		Nickname:  nickname,
		Message:   message,
		Timestamp: serverTimeInLocal(timestamp),
	}
}

// DelayedMessages contains the delayed messages for specific date
type DelayedMessages struct {
	date     time.Time
	messages []*DelayedMessage
	lock     sync.RWMutex
}

func newDelayedMessages(date time.Time) *DelayedMessages {
	return &DelayedMessages{
		date: serverTimeInLocal(date),
	}
}

// GetDate returns the delayed messages group's date
func (dm *DelayedMessages) GetDate() time.Time {
	return dm.date
}

// GetMessages returns a list of delayed messages
func (dm *DelayedMessages) GetMessages() []*DelayedMessage {
	dm.lock.RLock()
	defer dm.lock.RUnlock()

	result := []*DelayedMessage{}
	result = append(result, dm.messages...)

	return result
}

func (dm *DelayedMessages) add(nickname, message string, timestamp time.Time) {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	shouldAddDelayedMessage := true
	if len(dm.messages) > 0 {
		lastMessage := dm.messages[len(dm.messages)-1]
		shouldAddDelayedMessage = lastMessage.Timestamp.Before(timestamp)
	}

	if shouldAddDelayedMessage {
		dm.messages = append(dm.messages, newDelayedMessage(nickname, message, timestamp))
	}
}

// DiscussionHistory contains the rooms's discussion history
type DiscussionHistory struct {
	history []*DelayedMessages
	lock    sync.RWMutex
}

// NewDiscussionHistory creates a new discussion history instance
func NewDiscussionHistory() *DiscussionHistory {
	return &DiscussionHistory{}
}

// GetHistory returns the discussion history
func (dh *DiscussionHistory) GetHistory() []*DelayedMessages {
	dh.lock.RLock()
	defer dh.lock.RUnlock()

	return dh.history
}

// AddMessage add a new delayed message to the history
func (dh *DiscussionHistory) AddMessage(nickname, message string, timestamp time.Time) {
	t := serverTimeInLocal(timestamp)

	for _, dm := range dh.GetHistory() {
		if sameDate(dm.date, t) {
			dm.add(nickname, message, t)
			return
		}
	}

	dm := dh.addNewMessagesGroup(t)
	dm.add(nickname, message, t)
}

func (dh *DiscussionHistory) addNewMessagesGroup(date time.Time) *DelayedMessages {
	dh.lock.Lock()
	defer dh.lock.Unlock()

	dm := newDelayedMessages(date)
	dh.history = append(dh.history, dm)

	return dm
}

func sameDate(d1, d2 time.Time) bool {
	t1y, t1m, t1d := d1.In(time.Local).Date()
	t2y, t2m, t2d := d2.In(time.Local).Date()

	return t1d == t2d && t1m == t2m && t1y == t2y
}

func serverTimeInLocal(t time.Time) time.Time {
	return t.In(time.Local)
}
