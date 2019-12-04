// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/freman/cymru/ipasn"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	testErrors := []error{
		ipasn.ErrIPIsUnspecified,
		ipasn.ErrIPIsLoopback,
		ipasn.ErrIPIsMulticast,
		ipasn.ErrIPIsPrivate,
		ipasn.ErrNotFound,
	}

	for i, err := range testErrors {
		i, err := i, err
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			t.Parallel()
			terr := fmt.Errorf("Wrapped %w", err)
			require.True(t, errors.Is(terr, err))
		})
	}

	for i, err := range testErrors {
		i, err := i, err
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			t.Parallel()
			terr := fmt.Errorf("Wrapped %w", err)
			var verr ipasn.Error
			require.True(t, errors.As(terr, &verr))
			require.Equal(t, err, verr)
		})
	}
}
