package access

import (
	"bytes"
	"time"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/roster"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/otr3"
)

// EventHandler represents the main notifications that the session can emit
// It's really more an observer than an even handler
type EventHandler interface {
	RegisterCallback(title, instructions string, fields []interface{}) error
}

// Connector represents something that connect
type Connector interface {
	Connect()
}

// Session is an interface that defines the functionality of a Session
type Session interface {
	ApprovePresenceSubscription(jid.WithoutResource, string) error
	AutoApprove(string)
	AwaitVersionReply(<-chan data.Stanza, string)
	Close()
	CommandManager() otrclient.CommandManager
	Config() *config.ApplicationConfig
	Conn() xi.Conn
	Connect(string, tls.Verifier) error
	ConversationManager() otrclient.ConversationManager
	CreateSymmetricKeyFor(jid.Any) []byte
	DenyPresenceSubscription(jid.WithoutResource, string) error
	DisplayName() string
	EncryptAndSendTo(jid.Any, string) (int, bool, error)
	GetConfig() *config.Account
	GetInMemoryLog() *bytes.Buffer
	GroupDelimiter() string
	HandleConfirmOrDeny(jid.WithoutResource, bool)
	IsConnected() bool
	IsDisconnected() bool
	ManuallyEndEncryptedChat(jid.Any) error
	PrivateKeys() []otr3.PrivateKey
	R() *roster.List
	ReloadKeys()
	RemoveContact(string)
	RequestPresenceSubscription(jid.WithoutResource, string) error
	Send(jid.Any, string, bool) error
	SendMUCMessage(to, from, body string) error
	SendPing()
	SetCommandManager(otrclient.CommandManager)
	SetConnector(Connector)
	SetLastActionTime(time.Time)
	SetSessionEventHandler(EventHandler)
	SetWantToBeOnline(bool)
	Subscribe(chan<- interface{})
	Timeout(data.Cookie, time.Time)
	SendIQError(*data.ClientIQ, interface{})
	SendIQResult(*data.ClientIQ, interface{})
	PublishEvent(interface{})
	SendFileTo(jid.Any, string, func() bool, func(bool)) *sdata.FileTransferControl
	SendDirTo(jid.Any, string, func() bool, func(bool)) *sdata.FileTransferControl
	StartSMP(jid.WithResource, string, string)
	FinishSMP(jid.WithResource, string)
	AbortSMP(jid.WithResource)
	GetAndWipeSymmetricKeyFor(jid.Any) []byte

	HasRoom(jid.Bare, chan<- *muc.RoomListing) (<-chan bool, <-chan error)
	GetRooms(jid.Domain, string) (<-chan *muc.RoomListing, <-chan *muc.ServiceListing, <-chan error)
	JoinRoom(jid.Bare, string) error
	CreateRoom(jid.Bare) <-chan error
	GetChatServices(jid.Domain) (<-chan jid.Domain, <-chan error, func())
	GetRoomListing(jid.Bare, chan<- *muc.RoomListing)
	LoadRoomInfo(jid.Bare)
	LeaveRoom(room jid.Bare, nickname string) (<-chan bool, <-chan error)
	DestroyRoom(room, alternativeRoom jid.Bare, reason string) (<-chan bool, <-chan error, func())

	Log() coylog.Logger

	NewRoom(jid.Bare) *muc.Room
}

// Factory is a function that can create new Sessions
type Factory func(*config.ApplicationConfig, *config.Account, func(tls.Verifier, tls.Factory) xi.Dialer) Session
