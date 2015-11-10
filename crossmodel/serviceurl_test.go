// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package crossmodel_test

import (
	"fmt"
	"regexp"
	"strings"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/crossmodel"
)

type ServiceURLSuite struct{}

var _ = gc.Suite(&ServiceURLSuite{})

var urlTests = []struct {
	s, err string
	exact  string
	url    *crossmodel.ServiceURL
}{{
	s:   "local:/u/user/name",
	url: &crossmodel.ServiceURL{"local", "user", "name"},
}, {
	s:   "nonlocal:/u/user/name",
	err: "service URL has invalid schema: $URL",
}, {
	s:     "/u/user/name",
	url:   &crossmodel.ServiceURL{"local", "user", "name"},
	exact: "local:/u/user/name",
}, {
	s:     "u/user/name",
	url:   &crossmodel.ServiceURL{"local", "user", "name"},
	exact: "local:/u/user/name",
}, {
	s:   "local:service",
	err: `service URL has invalid form, missing "/u/<user>": $URL`,
}, {
	s:   "local:user/service",
	err: `service URL has invalid form, missing "/u/<user>": $URL`,
}, {
	s:   "local:/u/user",
	err: `service URL has invalid form, missing service name: $URL`,
}, {
	s:   "service",
	err: `service URL has invalid form, missing "/u/<user>": $URL`,
}, {
	s:   "/user/service",
	err: `service URL has invalid form, missing "/u/<user>": $URL`,
}, {
	s:   "/u/user",
	err: `service URL has invalid form, missing service name: $URL`,
}, {
	s:   "local:/u/user/service.bad",
	err: `service name "service.bad" not valid`,
}, {
	s:   "local:/u/user[bad/service",
	err: `user name "user\[bad" not valid`,
}, {
	s:   ":foo",
	err: `cannot parse service URL: $URL`,
}}

func (s *ServiceURLSuite) TestParseURL(c *gc.C) {
	for i, t := range urlTests {
		c.Logf("test %d: %q", i, t.s)
		url, err := crossmodel.ParseServiceURL(t.s)

		match := t.s
		if t.exact != "" {
			match = t.exact
		}
		if t.url != nil {
			c.Assert(err, gc.IsNil)
			c.Assert(url, gc.DeepEquals, t.url)
			c.Check(url.String(), gc.Equals, match)
		}
		if t.err != "" {
			t.err = strings.Replace(t.err, "$URL", regexp.QuoteMeta(fmt.Sprintf("%q", t.s)), -1)
			c.Assert(err, gc.ErrorMatches, t.err)
			c.Assert(url, gc.IsNil)
		}
	}
}
