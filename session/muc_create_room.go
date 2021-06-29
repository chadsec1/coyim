package session

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

var (
	// ErrInvalidReserveRoomRequest is an invalid room reservation request error
	ErrInvalidReserveRoomRequest = errors.New("invalid reserve room request")
)

// CreateInstantRoom will create a room "instantly" accepting the default configuration of the room
// For more information see XEP-0045 v1.32.0, section: 10.1.2
func (s *session) CreateInstantRoom(roomID jid.Bare) (<-chan bool, <-chan error) {
	c := s.muc.newCreateMUCInstantRoomContext(roomID)
	return c.createInstantRoom()
}

type createMUCInstantRoomContext struct {
	*createMUCRoomContext
	resultChannel chan bool
}

func (m *mucManager) newCreateMUCInstantRoomContext(roomID jid.Bare) *createMUCInstantRoomContext {
	return &createMUCInstantRoomContext{
		createMUCRoomContext: m.newCreateMUCRoomContext(roomID),
	}
}

func (c *createMUCInstantRoomContext) createInstantRoom() (<-chan bool, <-chan error) {
	c.resultChannel = make(chan bool)
	go c.createRoom(c.sendIQForInstantRoom, func(stanza data.Stanza) error {
		err := c.validateStanzaReceived(stanza)
		if err != nil {
			return err
		}
		c.resultChannel <- true
		return nil
	})
	return c.resultChannel, c.errorChannel
}

func (c *createMUCInstantRoomContext) sendIQForInstantRoom() (<-chan data.Stanza, error) {
	return c.sendInformationQuery("set", c.newRoomConfigurationFormSubmit())
}

// CreateReservedRoom will reserve a room and request the configuration form for it
func (s *session) CreateReservedRoom(roomID jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error) {
	c := s.muc.newCreateMUCReservedRoomContext(roomID)
	return c.createReservedRoom()
}

type createMUCReservedRoomContext struct {
	*createMUCRoomContext
	configFormChannel chan *muc.RoomConfigForm
}

func (m *mucManager) newCreateMUCReservedRoomContext(roomID jid.Bare) *createMUCReservedRoomContext {
	return &createMUCReservedRoomContext{
		createMUCRoomContext: m.newCreateMUCRoomContext(roomID),
		configFormChannel:    make(chan *muc.RoomConfigForm),
	}
}

func (c *createMUCReservedRoomContext) createReservedRoom() (<-chan *muc.RoomConfigForm, <-chan error) {
	c.configFormChannel = make(chan *muc.RoomConfigForm)

	go c.createRoom(c.sendIQForReservedRoom, func(stanza data.Stanza) error {
		form, err := c.getConfigFormFromStanza(stanza)
		if err != nil {
			return err
		}
		c.configFormChannel <- form
		return nil
	})

	return c.configFormChannel, c.errorChannel
}

func (c *createMUCReservedRoomContext) sendIQForReservedRoom() (<-chan data.Stanza, error) {
	return c.sendInformationQuery("get", c.newRoomConfigurationFormRequest())
}

func (c *createMUCReservedRoomContext) getConfigFormFromStanza(stanza data.Stanza) (*muc.RoomConfigForm, error) {
	iq, err := c.getIQFromStanza(stanza)
	if err != nil {
		return nil, err
	}

	cf, err := c.getConfigFormFromIQResponse(iq)
	if err != nil {
		return nil, err
	}

	return muc.NewRoomConfigForm(cf.Form), nil
}

func (c *createMUCReservedRoomContext) getConfigFormFromIQResponse(iq *data.ClientIQ) (*data.MUCRoomConfiguration, error) {
	cf := &data.MUCRoomConfiguration{}
	err := xml.Unmarshal(iq.Query, cf)
	if err != nil {
		return nil, err
	}
	return cf, nil
}

type createMUCRoomContext struct {
	roomID       jid.Bare
	errorChannel chan error
	conn         xi.Conn
	log          coylog.Logger
}

func (m *mucManager) newCreateMUCRoomContext(roomID jid.Bare) *createMUCRoomContext {
	c := &createMUCRoomContext{
		roomID:       roomID,
		errorChannel: make(chan error),
		conn:         m.conn(),
		log:          m.log.WithField("where", "createRoomContext"),
	}

	return c
}

// See XEP-0045 v1.32.0, section: 10.1.1
func (c *createMUCRoomContext) reserveRoom() bool {
	err := c.sendMUCPresence()
	if err != nil {
		c.errorChannel <- err
		return false
	}
	return true
}

func (c *createMUCRoomContext) createRoom(sendIQ func() (<-chan data.Stanza, error), onStanzaReceived func(stanza data.Stanza) error) {
	if !c.reserveRoom() {
		c.error(ErrInvalidReserveRoomRequest, "An error occurred while reserving the room")
		return
	}

	reply, err := sendIQ()
	if err != nil {
		c.error(ErrUnexpectedResponse, "Unexpected information query response")
		return
	}

	stanza, ok := <-reply
	if !ok {
		c.error(ErrInvalidInformationQueryRequest, "Unexpected information query reply")
		return
	}

	err = onStanzaReceived(stanza)
	if err != nil {
		c.error(err, "An error occurred when the stanza was received")
	}
}

func (c *createMUCRoomContext) error(err error, m string) {
	c.logError(err, m)
	c.errorChannel <- err
}

func (c *createMUCRoomContext) logError(err error, m string) {
	c.log.WithError(err).Error(m)
}

func (c *createMUCRoomContext) sendMUCPresence() error {
	err := c.conn.SendMUCPresence(c.roomID.String(), &data.MUC{})
	if err != nil {
		c.logError(err, "An error ocurred while sending a presence for creating an instant room")
		return ErrUnexpectedResponse
	}

	return nil
}

func (c *createMUCRoomContext) sendInformationQuery(tp string, d interface{}) (<-chan data.Stanza, error) {
	reply, _, err := c.conn.SendIQ(c.roomID.String(), tp, d)
	if err != nil {
		c.logError(err, "An error ocurred while sending the information query")
		return nil, err
	}

	return reply, nil
}

func (c *createMUCRoomContext) getIQFromStanza(stanza data.Stanza) (*data.ClientIQ, error) {
	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, ErrUnexpectedResponse
	}

	if iq.Type == "error" {
		return nil, ErrInformationQueryResponse
	}

	return iq, nil
}

func (c *createMUCRoomContext) validateStanzaReceived(stanza data.Stanza) error {
	_, err := c.getIQFromStanza(stanza)
	return err
}

func (c *createMUCRoomContext) newRoomConfigurationFormSubmit() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{
		Form: &data.Form{
			Type: "submit",
		},
	}
}

func (c *createMUCRoomContext) newRoomConfigurationFormRequest() data.MUCRoomConfiguration {
	return data.MUCRoomConfiguration{}
}
