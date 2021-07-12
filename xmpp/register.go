// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

var (
	// ErrUsernameConflict is an error signaled during account registration, when the username is not available
	ErrUsernameConflict = errors.New("xmpp: the username is not available for registration")
	// ErrMissingRequiredRegistrationInfo is an error signaled during account registration, when some required registration information is missing
	ErrMissingRequiredRegistrationInfo = errors.New("xmpp: missing required registration information")
	// ErrRegistrationFailed is an error signaled during account registration, when account creation failed
	ErrRegistrationFailed = errors.New("xmpp: account creation failed")
	// ErrWrongCaptcha is an error signaled during account registration, when the captcha entered is wrong
	ErrWrongCaptcha = errors.New("xmpp: the captcha entered is wrong")
	// ErrResourceConstraint is an error signaled during account registration, when the configured number of allowable resources is reached
	ErrResourceConstraint = errors.New("xmpp: already reached the configured number of allowable resources")
)

var (
	// ErrNotAllowed is an error signaled during password change, when the server does not allow password change
	ErrNotAllowed = errors.New("xmpp: server does not allow password change")
	// ErrNotAuthorized is an error signaled during password change, when password change is not authorized
	ErrNotAuthorized = errors.New("xmpp: password change not authorized")
	// ErrBadRequest is an error signaled during password change, when the request is malformed
	ErrBadRequest = errors.New("xmpp: malformed password change request")
	// ErrInternalServerError is an error signaled when the server has an internal error
	ErrInternalServerError = errors.New("xmpp: internal server error")
	// ErrUnexpectedRequest is an error signaled during password change, when the user is not registered in the server
	ErrUnexpectedRequest = errors.New("xmpp: user is not registered in the server")
	// ErrChangePasswordFailed is an error signaled during password change, when it fails
	ErrChangePasswordFailed = errors.New("xmpp: password change failed")
)

// XEP-0077
func (d *dialer) negotiateInBandRegistration(c interfaces.Conn) (bool, error) {
	if c.Features().InBandRegistration == nil {
		return false, nil
	}

	user := d.getJIDLocalpart()
	return c.RegisterAccount(user, d.password)
}

func (c *conn) RegisterAccount(user, password string) (bool, error) {
	if c.config.CreateCallback == nil {
		return false, nil
	}

	err := c.createAccount(user, password)
	if err != nil {
		return false, err
	}

	return true, c.closeImmediately()
}

func (c *conn) createAccount(user, password string) error {
	c.ioLock.Lock()
	defer c.ioLock.Unlock()

	c.log.Debug("Attempting to create account")
	fmt.Fprintf(c.out, "<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>")
	var iq data.ClientIQ
	if err := c.in.DecodeElement(&iq, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}

	if iq.Type != "result" {
		return ErrRegistrationFailed
	}

	c.log.Debug("createAccount() - received the registration form")

	var register data.RegisterQuery

	if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&register); err != nil {
		return err
	}

	if len(register.Form.Type) > 0 {
		c.log.Debug("createAccount() - processing form")
		reply, err := processForm(&register.Form, register.Datas, c.config.CreateCallback)
		if err != nil {
			c.log.WithError(err).Debug("createAccount() - couldn't process form")
			return err
		}

		fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'>")

		if err = xml.NewEncoder(c.rawOut).Encode(reply); err != nil {
			return err
		}

		fmt.Fprintf(c.rawOut, "</query></iq>")
		c.log.Debug("createAccount() - have sent the IQ with registration information")
	} else if register.Username != nil && register.Password != nil {
		//TODO: make sure this only happens via SSL
		//TODO: should generate form asking for username and password,
		//and call processForm for consistency

		// Try the old-style registration.
		fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query></iq>", user, password)
	}

	iq2 := &data.ClientIQ{}
	if err := c.in.DecodeElement(iq2, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}

	if iq2.Type == "error" {
		switch {
		case iq2.Error.MUCConflict != nil:
			return ErrUsernameConflict
		case iq2.Error.MUCNotAcceptable != nil:
			return ErrMissingRequiredRegistrationInfo
		case iq2.Error.MUCNotAllowed != nil:
			return ErrWrongCaptcha
		default:
			switch iq2.Error.Condition.XMLName.Local {
			case "bad-request":
				return ErrRegistrationFailed
			case "resource-constraint":
				return ErrResourceConstraint
			default:
				return ErrRegistrationFailed
			}
		}
	}

	c.log.Debug("createAccount() - received a successful response")

	return nil
}

// CancelRegistration cancels the account registration with the server
func (c *conn) CancelRegistration() (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	// https://xmpp.org/extensions/xep-0077.html#usecases-cancel
	registrationCancel := rawXML(`
	<query xmlns='jabber:iq:register'>
		<remove/>
	</query>
	`)

	return c.SendIQ("", "set", registrationCancel)
}

func (c *conn) sendChangePasswordInfo(username, server, password string) (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	// TODO: we might be able to put this in a struct
	changePasswordXML := fmt.Sprintf("<query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query>", username, password)
	return c.SendIQ(server, "set", rawXML(changePasswordXML))
}

// ChangePassword changes the account password registered in the server
// Reference: https://xmpp.org/extensions/xep-0077.html#usecases-changepw
func (c *conn) ChangePassword(username, server, password string) error {
	c.log.WithField("user", username).Debug("Attempting to change account's password")

	reply, _, err := c.sendChangePasswordInfo(username, server, password)
	if err != nil {
		return errors.New("xmpp: failed to send request")
	}

	stanza, ok := <-reply
	if !ok {
		return errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return errors.New("xmpp: failed to parse response")
	}

	if iq.Type == "result" {
		return nil
	}

	// TODO: server can also return a form requiring more information from the user. This should be rendered.
	if iq.Type == "error" {
		if iq.Error.MUCNotAllowed != nil {
			return ErrNotAllowed
		}

		if iq.Error.MUCNotAuthorized != nil {
			return ErrNotAuthorized
		}

		if iq.Error.MUCBadRequest != nil {
			return ErrBadRequest
		}

		if iq.Error.MUCInternalServerError != nil {
			return ErrInternalServerError
		}

		switch iq.Error.Condition.XMLName.Local {
		case "bad-request":
			return ErrBadRequest
		case "unexpected-request":
			return ErrUnexpectedRequest
		default:
			return ErrChangePasswordFailed
		}
	}

	return ErrChangePasswordFailed
}
