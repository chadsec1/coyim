package data

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
)

//Any provides a convenient way to debug any child element
type Any struct {
	XMLName xml.Name
	Body    string `xml:",innerxml"`
}

// Extensions implements generic XEPs.
type Extensions []*Extension

// Extension represents any XML node not included in the Stanza definition
type Extension Any

// GoString implements the GoStringer interface
func (e *Extensions) GoString() string {
	if e == nil {
		return "<nil>"
	}

	result := []string{}

	for _, ee := range *e {
		result = append(result, fmt.Sprintf("%#v", ee))
	}

	return fmt.Sprintf("{%s}", strings.Join(result, ", "))
}

// GoString implements the GoStringer interface
func (e *Extension) GoString() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("<Extension %v body=%q>", e.XMLName, e.Body)
}

// StanzaErrorCondition represents an error condition inside StanzaError
type StanzaErrorCondition struct {
	XMLName xml.Name
	Body    string `xml:",innerxml"`
}

// StanzaError implements RFC 3920, section 9.3.
//TODO RFC 6120 obsoletes RFC 3920, section "8.3.2. Syntax"
type StanzaError struct {
	XMLName xml.Name `xml:"jabber:client error"`

	By   string `xml:"by,attr"`
	Code string `xml:"code,attr"` // NOTE: I could not find this in the spec

	// cancel -- do not retry (the error is unrecoverable)
	// continue -- proceed (the condition was only a warning)
	// modify -- retry after changing the data sent
	// auth -- retry after providing credentials
	// wait -- retry after waiting (the error is temporary)
	Type string `xml:"type,attr"`
	Text string `xml:"urn:ietf:params:xml:ns:xmpp-stanzas text"`

	Condition StanzaErrorCondition `xml:",any"`

	//	ApplicationCondition *Any `xml:",any,omitempty"`

	MUCNotAuthorized        *MUCNotAuthorized
	MUCForbidden            *MUCForbidden
	MUCItemNotFound         *MUCItemNotFound
	MUCNotAllowed           *MUCNotAllowed
	MUCNotAcceptable        *MUCNotAcceptable
	MUCRegistrationRequired *MUCRegistrationRequired
	MUCConflict             *MUCConflict
	MUCServiceUnavailable   *MUCServiceUnavailable
}

func anyNotNil(vals ...interface{}) interface{} {
	for _, v := range vals {
		if !reflect.ValueOf(v).IsNil() {
			return v
		}
	}
	return nil
}

// AnyMUCError will return any of the potential MUC errors, or nil
func (e *StanzaError) AnyMUCError() interface{} {
	return anyNotNil(
		e.MUCNotAuthorized,
		e.MUCForbidden,
		e.MUCItemNotFound,
		e.MUCNotAllowed,
		e.MUCNotAcceptable,
		e.MUCRegistrationRequired,
		e.MUCConflict,
		e.MUCServiceUnavailable,
	)
}

// Subject contains all the subject text of any room when it is received
type Subject struct {
	XMLName xml.Name `xml:"subject"`
	Text    string   `xml:",innerxml"`
}

// ClientMessage implements RFC 3921  B.1  jabber:client
type ClientMessage struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr"`
	ID      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"` // chat, error, groupchat, headline, or normal

	// These should technically be []clientText,
	// but string is much more convenient.
	Subject    *Subject    `xml:",omitempty"`
	Body       string      `xml:"body"`
	Thread     string      `xml:"thread,omitempty"`
	Delay      *Delay      `xml:"delay,omitempty"`
	Encryption *Encryption `xml:"encryption,omitempty"`

	Error *StanzaError `xml:"error"`

	MUCUser *MUCUser `xml:",omitempty"`

	Extensions `xml:",any,omitempty"`
}

// ClientCaps contains information about client capabilities
type ClientCaps struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Ext     string   `xml:"ext,attr"`
	Hash    string   `xml:"hash,attr"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
}
