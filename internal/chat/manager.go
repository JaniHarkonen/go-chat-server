package chat

import (
	"fmt"

	"github.com/JaniHarkonen/go-chat-server/internal/utils"
)

type Manager struct {
	snapshot        []*Message
	visibleLength   int
	activeThreshold int
	messageCount    int
	activeUsers     map[*User]int
}

func NewManager(activeThreshold int, visibleLength int) *Manager {
	return &Manager{
		snapshot:        make([]*Message, 0, activeThreshold+1),
		visibleLength:   visibleLength,
		activeThreshold: activeThreshold,
		messageCount:    0,
		activeUsers:     make(map[*User]int),
	}
}

// Posts a new message to the chat and updates the active user list by activating the poster
// if necessary, and deactivating the last poster if they haven't posted a message ever since.
// This function will result in at most one user being activated and one user being deactivated.
// The activated and the deactivated user will be returned.
func (cm *Manager) Post(u *User, msg *string) (activated *User, deactivated *User) {
	activated = nil
	deactivated = nil

	if _, ok := cm.activeUsers[u]; !ok {
		activated = u
	}

	cm.messageCount++
	cm.activeUsers[u] = cm.messageCount
	cm.snapshot = append(cm.snapshot, newMessage(u, cm.messageCount, msg))

	// Update active users by deleting inactive ones
	lastUser := cm.snapshot[0].user
	if cm.messageCount > cm.activeThreshold {
		if last, ok := cm.activeUsers[lastUser]; ok && last <= cm.messageCount-cm.activeThreshold {
			fmt.Println(last, cm.messageCount-cm.activeThreshold)
			delete(cm.activeUsers, lastUser)
			deactivated = lastUser
		}
		cm.snapshot = cm.snapshot[1:]
	}

	return activated, deactivated
}

func (cm *Manager) IsUserActive(u *User) bool {
	_, ok := cm.activeUsers[u]
	return ok
}

func (cm *Manager) VisibleMessages() []*Message {
	return cm.snapshot[utils.MaxInt(0, len(cm.snapshot)-cm.visibleLength):len(cm.snapshot)]
}

func (cm *Manager) ActiveUsers() map[*User]int {
	return cm.activeUsers
}

func (cm *Manager) Snapshot() []*Message {
	return cm.snapshot
}
