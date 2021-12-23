package data

import "fmt"

const (
	// AffiliationOwner represents XMPP muc 'owner' affiliation
	AffiliationOwner = "owner"
	// AffiliationAdmin represents XMPP muc 'admin' affiliation
	AffiliationAdmin = "admin"
	// AffiliationMember represents XMPP muc 'member' affiliation
	AffiliationMember = "member"
	// AffiliationOutcast represents XMPP muc 'outcast' affiliation
	AffiliationOutcast = "outcast"
	// AffiliationNone represents XMPP muc 'none' affiliation
	AffiliationNone = "none"
)

type affiliationNumberType int

const (
	affiliationTypeOutcast affiliationNumberType = iota
	affiliationTypeNone
	affiliationTypeMember
	affiliationTypeAdmin
	affiliationTypeOwner
)

var affiliationLowerThan = [][]bool{
	{false /*none*/, false /*outcast*/, true /*member*/, true /*admin*/, true /*owner*/},    // none
	{false /*none*/, false /*outcast*/, true /*member*/, true /*admin*/, true /*owner*/},    // outcast
	{false /*none*/, false /*outcast*/, false /*member*/, true /*admin*/, true /*owner*/},   // member
	{false /*none*/, false /*outcast*/, false /*member*/, false /*admin*/, true /*owner*/},  // admin
	{false /*none*/, false /*outcast*/, false /*member*/, false /*admin*/, false /*owner*/}, // owner
}

// AffiliationUpdate contains information related to a new and previous affiliation
type AffiliationUpdate struct {
	Nickname string
	Reason   string
	New      Affiliation
	Previous Affiliation
	Actor    *Actor
}

// SelfAffiliationUpdate contains information related to a new and previous affiliation of the self occupant
type SelfAffiliationUpdate struct {
	AffiliationUpdate
}

// AffiliationVisitor implements the Visitor pattern for affiliations
type AffiliationVisitor interface {
	OnNone(*NoneAffiliation)
	OnOutcast(*OutcastAffiliation)
	OnMember(*MemberAffiliation)
	OnAdmin(*AdminAffiliation)
	OnOwner(*OwnerAffiliation)
}

// Affiliation represents an affiliation as specificed by section 5.2 in XEP-0045
type Affiliation interface {
	// IsAdmin will return true if this specific affiliation can modify persistent information
	IsAdmin() bool
	// IsBanned will return true if this specific affiliation means that the jid is banned from the room
	IsBanned() bool
	// IsMember will return true if this specific affiliation means that the jid is a member of the room
	IsMember() bool
	// IsOwner will return true if this specific affiliation means that the jid is an owner of the room
	IsOwner() bool
	// IsNone will return true if if the jid doesn't have affiliation
	IsNone() bool
	// Name returns the string name of the affiliation type
	Name() string
	// IsDifferentFrom returns true if the caller affiliation is different from the affiliation passed as argument
	IsDifferentFrom(Affiliation) bool
	// affiliationTypeAsNumber returns an int value indicating the affiliation number through a AffiliationNumberType
	affiliationTypeAsNumber() affiliationNumberType
	// Visit implements double dispatch an affiliations
	Visit(AffiliationVisitor)
}

// NoneAffiliation is a representation of MUC's "none" affiliation
type NoneAffiliation struct{}

// OutcastAffiliation is a representation of MUC's "banned" affiliation
type OutcastAffiliation struct{}

// MemberAffiliation is a representation of MUC's "member" affiliation
type MemberAffiliation struct{}

// AdminAffiliation is a representation of MUC's "admin" affiliation
type AdminAffiliation struct{}

// OwnerAffiliation is a representation of MUC's "owner" affiliation
type OwnerAffiliation struct{}

// IsAdmin implements Affiliation interface
func (*NoneAffiliation) IsAdmin() bool { return false }

// IsAdmin implements Affiliation interface
func (*OutcastAffiliation) IsAdmin() bool { return false }

// IsAdmin implements Affiliation interface
func (*MemberAffiliation) IsAdmin() bool { return false }

// IsAdmin implements Affiliation interface
func (*AdminAffiliation) IsAdmin() bool { return true }

// IsAdmin implements Affiliation interface
func (*OwnerAffiliation) IsAdmin() bool { return false }

// IsBanned implements Affiliation interface
func (*NoneAffiliation) IsBanned() bool { return false }

// IsBanned implements Affiliation interface
func (*OutcastAffiliation) IsBanned() bool { return true }

// IsBanned implements Affiliation interface
func (*MemberAffiliation) IsBanned() bool { return false }

// IsBanned implements Affiliation interface
func (*AdminAffiliation) IsBanned() bool { return false }

// IsBanned implements Affiliation interface
func (*OwnerAffiliation) IsBanned() bool { return false }

// IsMember implements Affiliation interface
func (*NoneAffiliation) IsMember() bool { return false }

// IsMember implements Affiliation interface
func (*OutcastAffiliation) IsMember() bool { return false }

// IsMember implements Affiliation interface
func (*MemberAffiliation) IsMember() bool { return true }

// IsMember implements Affiliation interface
func (*AdminAffiliation) IsMember() bool { return false }

// IsMember implements Affiliation interface
func (*OwnerAffiliation) IsMember() bool { return false }

// IsModerator implements Affiliation interface
func (*NoneAffiliation) IsModerator() bool { return false }

// IsModerator implements Affiliation interface
func (*OutcastAffiliation) IsModerator() bool { return false }

// IsModerator implements Affiliation interface
func (*MemberAffiliation) IsModerator() bool { return false }

// IsModerator implements Affiliation interface
func (*AdminAffiliation) IsModerator() bool { return true }

// IsModerator implements Affiliation interface
func (*OwnerAffiliation) IsModerator() bool { return true }

// IsOwner implements Affiliation interface
func (*NoneAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*OutcastAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*MemberAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*AdminAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*OwnerAffiliation) IsOwner() bool { return true }

// IsOutcast implements Affiliation interface
func (*NoneAffiliation) IsOutcast() bool { return false }

// IsOutcast implements Affiliation interface
func (*OutcastAffiliation) IsOutcast() bool { return true }

// IsOutcast implements Affiliation interface
func (*MemberAffiliation) IsOutcast() bool { return false }

// IsOutcast implements Affiliation interface
func (*AdminAffiliation) IsOutcast() bool { return false }

// IsOutcast implements Affiliation interface
func (*OwnerAffiliation) IsOutcast() bool { return false }

// IsNone implements Affiliation interface
func (*NoneAffiliation) IsNone() bool { return true }

// IsNone implements Affiliation interface
func (*OutcastAffiliation) IsNone() bool { return false }

// IsNone implements Affiliation interface
func (*MemberAffiliation) IsNone() bool { return false }

// IsNone implements Affiliation interface
func (*AdminAffiliation) IsNone() bool { return false }

// IsNone implements Affiliation interface
func (*OwnerAffiliation) IsNone() bool { return false }

// Name implements Affiliation interface
func (*NoneAffiliation) Name() string { return AffiliationNone }

// Name implements Affiliation interface
func (*OutcastAffiliation) Name() string { return AffiliationOutcast }

// Name implements Affiliation interface
func (*MemberAffiliation) Name() string { return AffiliationMember }

// Name implements Affiliation interface
func (*AdminAffiliation) Name() string { return AffiliationAdmin }

// Name implements Affiliation interface
func (*OwnerAffiliation) Name() string { return AffiliationOwner }

// IsDifferentFrom implements Affiliation interface
func (a *NoneAffiliation) IsDifferentFrom(a1 Affiliation) bool {
	return areDifferentAffiliations(a, a1)
}

// IsDifferentFrom implements Affiliation interface
func (a *OutcastAffiliation) IsDifferentFrom(a1 Affiliation) bool {
	return areDifferentAffiliations(a, a1)
}

// IsDifferentFrom implements Affiliation interface
func (a *MemberAffiliation) IsDifferentFrom(a1 Affiliation) bool {
	return areDifferentAffiliations(a, a1)
}

// IsDifferentFrom implements Affiliation interface
func (a *AdminAffiliation) IsDifferentFrom(a1 Affiliation) bool {
	return areDifferentAffiliations(a, a1)
}

// IsDifferentFrom implements Affiliation interface
func (a *OwnerAffiliation) IsDifferentFrom(a1 Affiliation) bool {
	return areDifferentAffiliations(a, a1)
}

// affiliationTypeAsNumber implements Affiliation interface
func (*NoneAffiliation) affiliationTypeAsNumber() affiliationNumberType {
	return affiliationTypeNone
}

// affiliationTypeAsNumber implements Affiliation interface
func (*OutcastAffiliation) affiliationTypeAsNumber() affiliationNumberType {
	return affiliationTypeOutcast
}

// affiliationTypeAsNumber implements Affiliation interface
func (*MemberAffiliation) affiliationTypeAsNumber() affiliationNumberType {
	return affiliationTypeMember
}

// affiliationTypeAsNumber implements Affiliation interface
func (*AdminAffiliation) affiliationTypeAsNumber() affiliationNumberType {
	return affiliationTypeAdmin
}

// affiliationTypeAsNumber implements Affiliation interface
func (*OwnerAffiliation) affiliationTypeAsNumber() affiliationNumberType {
	return affiliationTypeOwner
}

// Visit implements the Affiliation interface
func (a *NoneAffiliation) Visit(v AffiliationVisitor) {
	v.OnNone(a)
}

// Visit implements the Affiliation interface
func (a *OutcastAffiliation) Visit(v AffiliationVisitor) {
	v.OnOutcast(a)
}

// Visit implements the Affiliation interface
func (a *MemberAffiliation) Visit(v AffiliationVisitor) {
	v.OnMember(a)
}

// Visit implements the Affiliation interface
func (a *AdminAffiliation) Visit(v AffiliationVisitor) {
	v.OnAdmin(a)
}

// Visit implements the Affiliation interface
func (a *OwnerAffiliation) Visit(v AffiliationVisitor) {
	v.OnOwner(a)
}

// AffiliationFromString returns an Affiliation from the given string, or an error if the string doesn't match a known affiliation type
func AffiliationFromString(s string) (Affiliation, error) {
	switch s {
	case AffiliationNone:
		return &NoneAffiliation{}, nil
	case AffiliationOutcast:
		return &OutcastAffiliation{}, nil
	case AffiliationMember:
		return &MemberAffiliation{}, nil
	case AffiliationAdmin:
		return &AdminAffiliation{}, nil
	case AffiliationOwner:
		return &OwnerAffiliation{}, nil
	default:
		return nil, fmt.Errorf("unknown affiliation string: '%s'", s)
	}
}

func areDifferentAffiliations(a, a1 Affiliation) bool {
	return a.affiliationTypeAsNumber() != a1.affiliationTypeAsNumber()
}
