package merklization

import "github.com/attestantio/go-eth2-client/spec/phase0"

var CAPELLA_FORK_VERSION = phase0.Version([4]byte{0x03, 0x00, 0x00, 0x00})

var DENEB_FORK_VERSION = phase0.Version([4]byte{0x04, 0x00, 0x00, 0x00})
