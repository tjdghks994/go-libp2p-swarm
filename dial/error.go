package dial

import (
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// maxDialDialErrors is the maximum number of dial errors we record
const maxDialDialErrors = 16

// DialError is the error type returned when dialing.
type Error struct {
	Peer       peer.ID
	DialErrors []TransportError
	Cause      error
	Skipped    int
}

func (e *Error) recordErr(addr ma.Multiaddr, err error) {
	if len(e.DialErrors) >= maxDialDialErrors {
		e.Skipped++
		return
	}
	e.DialErrors = append(e.DialErrors, TransportError{
		Address: addr,
		Cause:   err,
	})
}

func (e *Error) Error() string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "failed to dial %s:", e.Peer)
	if e.Cause != nil {
		_, _ = fmt.Fprintf(&builder, " %s", e.Cause)
	}
	for _, te := range e.DialErrors {
		_, _ = fmt.Fprintf(&builder, "\n  * [%s] %s", te.Address, te.Cause)
	}
	if e.Skipped > 0 {
		_, _ = fmt.Fprintf(&builder, "\n    ... skipping %d errors ...", e.Skipped)
	}
	return builder.String()
}

// Unwrap implements https://godoc.org/golang.org/x/xerrors#Wrapper.
func (e *Error) Unwrap() error {
	// If we have a context error, that's the "ultimate" error.
	if e.Cause != nil {
		return e.Cause
	}
	return nil
}

var _ error = (*Error)(nil)

// TransportError is the error returned when dialing a specific address.
type TransportError struct {
	Address ma.Multiaddr
	Cause   error
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("failed to dial %s: %s", e.Address, e.Cause)
}

var _ error = (*TransportError)(nil)