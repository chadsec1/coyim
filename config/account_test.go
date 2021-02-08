package config

import (
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

type AccountXMPPSuite struct{}

var _ = Suite(&AccountXMPPSuite{})

func (s *AccountXMPPSuite) Test_Account_Is_recognizesJids(c *C) {
	a := &Account{Account: "hello@bar.com"}
	c.Check(a.Is("foo"), Equals, false)
	c.Check(a.Is("hello@bar.com"), Equals, true)
	c.Check(a.Is("hello@bar.com/foo"), Equals, true)
}

func (s *AccountXMPPSuite) Test_Account_ShouldEncryptTo(c *C) {
	a := &Account{Account: "hello@bar.com", AlwaysEncrypt: false, AlwaysEncryptWith: []string{"one@foo.com", "two@foo.com"}}
	a2 := &Account{Account: "hello@bar.com", AlwaysEncrypt: true, AlwaysEncryptWith: []string{"one@foo.com", "two@foo.com"}, Peers: []*Peer{
		&Peer{UserID: "some@one.com", Nickname: "Hellois", EncryptionSettings: NeverEncrypt},
		&Peer{UserID: "some2@one.com", Nickname: "Hellois", EncryptionSettings: AlwaysEncrypt},
	}}
	c.Check(a.ShouldEncryptTo("foo"), Equals, false)
	c.Check(a.ShouldEncryptTo("hello@bar.com"), Equals, false)
	c.Check(a.ShouldEncryptTo("one@foo.com"), Equals, true)
	c.Check(a.ShouldEncryptTo("two@foo.com"), Equals, true)
	c.Check(a.ShouldEncryptTo("two@foo.com/blarg"), Equals, true)
	c.Check(a2.ShouldEncryptTo("foo"), Equals, true)
	c.Check(a2.ShouldEncryptTo("hello@bar.com"), Equals, true)

	c.Assert(a2.ShouldEncryptTo("some@one.com"), Equals, false)
	c.Assert(a2.ShouldEncryptTo("some2@one.com"), Equals, true)
}

func (s *AccountXMPPSuite) Test_NewAccount_ReturnsNewAccountWithSafeDefaults(c *C) {
	a, err := NewAccount()

	c.Check(err, IsNil)
	c.Check(len(a.PrivateKeys), Equals, 1)
	c.Check(a.AlwaysEncrypt, Equals, true)
	c.Check(a.OTRAutoStartSession, Equals, true)
	c.Check(a.OTRAutoTearDown, Equals, true)
	c.Check(a.Proxies, DeepEquals, []string{"tor-auto://"})
}

func (s *AccountXMPPSuite) Test_SetOTRPoliciesFor_SetupOTRPolicies(c *C) {
	a, _ := NewAccount()
	conv := &otr3.Conversation{}

	expectedConv := &otr3.Conversation{}
	expectedPolicies := expectedConv.Policies
	expectedPolicies.AllowV2()
	expectedPolicies.AllowV3()
	expectedPolicies.SendWhitespaceTag()
	expectedPolicies.WhitespaceStartAKE()
	expectedPolicies.RequireEncryption()
	expectedPolicies.ErrorStartAKE()

	a.SetOTRPoliciesFor("someon@jabber.com", conv)
	c.Check(conv.Policies, Equals, expectedPolicies)
}

func (s *AccountXMPPSuite) Test_SetOTRPoliciesFor_SetupOTRPoliciesWithOptionalEncription(c *C) {
	a, _ := NewAccount()
	a.AlwaysEncrypt = false
	conv := &otr3.Conversation{}

	expectedConv := &otr3.Conversation{}
	expectedPolicies := expectedConv.Policies
	expectedPolicies.AllowV2()
	expectedPolicies.AllowV3()
	expectedPolicies.SendWhitespaceTag()
	expectedPolicies.WhitespaceStartAKE()

	a.SetOTRPoliciesFor("someon@jabber.com", conv)
	c.Check(conv.Policies, Equals, expectedPolicies)
}

func (s *AccountXMPPSuite) Test_EnsurePrivateKey_DoesNotUpdateIfKeyExists(c *C) {
	a, _ := NewAccount()
	changed, err := a.EnsurePrivateKey()

	c.Check(err, IsNil)
	c.Check(changed, Equals, false)
}

func (s *AccountXMPPSuite) Test_EnsurePrivateKey_GeneratePrivateKeyIfMissing(c *C) {
	a := &Account{}
	changed, err := a.EnsurePrivateKey()

	c.Check(err, IsNil)
	c.Check(changed, Equals, true)
	c.Check(len(a.PrivateKeys), Equals, 1)
}

func (s *AccountXMPPSuite) Test_ID_generatesID(c *C) {
	a := &Account{}
	c.Check(a.ID(), Not(HasLen), 0)
}

func (s *AccountXMPPSuite) Test_ID_doesNotChangeID(c *C) {
	a := &Account{
		id: "existing",
	}
	c.Check(a.ID(), Equals, "existing")
}

func (s *AccountXMPPSuite) Test_Account_AllPrivateKeys(c *C) {
	a := &Account{
		id: "existing",
	}
	_, _ = a.EnsurePrivateKey()
	c.Assert(len(a.AllPrivateKeys()), Equals, 1)

	a.DeprecatedPrivateKey = []byte{1, 2, 3}
	c.Assert(len(a.AllPrivateKeys()), Equals, 2)
}

func (s *AccountXMPPSuite) Test_Account_ToggleAlwaysEncrypt(c *C) {
	a := &Account{
		id: "existing",
	}

	a.ToggleAlwaysEncrypt()
	c.Assert(a.AlwaysEncrypt, Equals, true)
}

func (s *AccountXMPPSuite) Test_Account_ToggleConnectAutomatically(c *C) {
	a := &Account{
		id: "existing",
	}

	a.ToggleConnectAutomatically()
	c.Assert(a.ConnectAutomatically, Equals, true)
}

func (s *AccountXMPPSuite) Test_Account_SaveCert(c *C) {
	a := &Account{
		id: "existing",
	}

	a.SaveCert("foo", "bar", []byte{1, 2, 3, 4})
	c.Assert(a.Certificates[0], DeepEquals, &CertificatePin{Subject: "foo", Issuer: "bar", Fingerprint: []byte{0x1, 0x2, 0x3, 0x4}, FingerprintType: "SHA3-256"})
}
