package chat

import (
	"errors"

	"github.com/JaniHarkonen/go-chat-server/internal/utils"
)

type Manager struct {
	snapshot        []*Message
	visibleLength   int
	activeThreshold int
	messageCount    int
	activeUsers     map[*User]int
	usernameTable   map[string]*User
	mutedUsers      map[*User]bool
}

func NewManager(activeThreshold int, visibleLength int) *Manager {
	return &Manager{
		snapshot:        make([]*Message, 0, activeThreshold+1),
		visibleLength:   visibleLength,
		activeThreshold: activeThreshold,
		messageCount:    0,
		activeUsers:     make(map[*User]int),
		usernameTable:   make(map[string]*User),
		mutedUsers:      make(map[*User]bool),
	}
}

func (cm *Manager) RegisterUser(u *User) error {
	if _, ok := cm.usernameTable[*u.name]; !ok {
		cm.usernameTable[*u.name] = u
		return nil
	}

	return errors.New("attempting to register a user whose username already exists in the table")
}

func (cm *Manager) UnregisterUser(u *User) {
	delete(cm.usernameTable, *u.name)
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
			delete(cm.activeUsers, lastUser)
			deactivated = lastUser
		}
		cm.snapshot = cm.snapshot[1:]
	}

	return activated, deactivated
}

func (cm *Manager) MuteUser(u *User) {
	if u == nil {
		return
	}

	cm.mutedUsers[u] = true
}

func (cm *Manager) UnmuteUser(u *User) {
	if u == nil {
		return
	}

	cm.mutedUsers[u] = false
}

func (cm *Manager) IsUserMuted(u *User) bool {
	status, ok := cm.mutedUsers[u]

	if !ok {
		return false
	}

	return status
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

func (cm *Manager) FindUserByName(username string) *User {
	if user, ok := cm.usernameTable[username]; ok {
		return user
	}

	return nil
}
