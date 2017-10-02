//go:generate stringer -type=RecordType

package recordType

type RecordType int

const (
	A RecordType = iota
	AAAA
	ANY
	CAA
	CNAME
	DNSKEY
	DS
	MX
	NS
	NSEC
	NSEC3
	OPENPGPKEY
	PTR
	RRSIG
	SOA
	SRV
	TLSA
	TXT
)
