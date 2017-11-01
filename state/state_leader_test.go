// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	"time" // Only used for time types.

	"github.com/juju/errors"
	gitjujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/txn"

	"github.com/juju/juju/core/globalclock"
	"github.com/juju/juju/core/leadership"
	coretesting "github.com/juju/juju/testing"
)

type LeadershipSuite struct {
	ConnSuite
	checker     leadership.Checker
	claimer     leadership.Claimer
	globalClock globalclock.Updater
}

var _ = gc.Suite(&LeadershipSuite{})

func (s *LeadershipSuite) SetUpTest(c *gc.C) {
	s.ConnSuite.SetUpTest(c)
	err := s.State.SetClockForTesting(s.Clock)
	c.Assert(err, jc.ErrorIsNil)
	s.checker = s.State.LeadershipChecker()
	s.claimer = s.State.LeadershipClaimer()
	s.globalClock, err = s.State.GlobalClockUpdater()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *LeadershipSuite) TestClaimValidatesApplicationname(c *gc.C) {
	err := s.claimer.ClaimLeadership("not/a/service", "u/0", time.Minute)
	c.Check(err, gc.ErrorMatches, `cannot claim lease "not/a/service": not an application name`)
	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *LeadershipSuite) TestClaimValidatesUnitName(c *gc.C) {
	err := s.claimer.ClaimLeadership("application", "not/a/unit", time.Minute)
	c.Check(err, gc.ErrorMatches, `cannot claim lease for holder "not/a/unit": not a unit name`)
	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *LeadershipSuite) TestClaimValidateDuration(c *gc.C) {
	err := s.claimer.ClaimLeadership("application", "u/0", 0)
	c.Check(err, gc.ErrorMatches, `cannot claim lease for 0s?: non-positive`)
	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *LeadershipSuite) TestCheckValidatesApplicationname(c *gc.C) {
	token := s.checker.LeadershipCheck("not/a/service", "u/0")
	err := token.Check(nil)
	c.Check(err, gc.ErrorMatches, `cannot check lease "not/a/service": not an application name`)
	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *LeadershipSuite) TestCheckValidatesUnitName(c *gc.C) {
	token := s.checker.LeadershipCheck("application", "not/a/unit")
	err := token.Check(nil)
	c.Check(err, gc.ErrorMatches, `cannot check holder "not/a/unit": not a unit name`)
	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *LeadershipSuite) TestBlockValidatesApplicationname(c *gc.C) {
	err := s.claimer.BlockUntilLeadershipReleased("not/a/service", nil)
	c.Check(err, gc.ErrorMatches, `cannot wait for lease "not/a/service" expiry: not an application name`)
	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *LeadershipSuite) TestClaimExpire(c *gc.C) {

	// Claim on behalf of one unit.
	err := s.claimer.ClaimLeadership("application", "application/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)

	// Claim on behalf of another.
	err = s.claimer.ClaimLeadership("application", "service/1", time.Minute)
	c.Check(err, gc.Equals, leadership.ErrClaimDenied)

	// Allow the first claim to expire.
	s.expire(c, "application")

	// Reclaim on behalf of another.
	err = s.claimer.ClaimLeadership("application", "service/1", time.Minute)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *LeadershipSuite) TestCheck(c *gc.C) {

	// Create a single token for use by the whole test.
	token := s.checker.LeadershipCheck("application", "application/0")

	// Claim leadership.
	err := s.claimer.ClaimLeadership("application", "application/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)

	// Check token reports current leadership state.
	var ops []txn.Op
	err = token.Check(&ops)
	c.Check(err, jc.ErrorIsNil)
	c.Check(ops, gc.HasLen, 1)

	// Allow leadership to expire.
	s.expire(c, "application")

	// Check leadership still reported accurately.
	var ops2 []txn.Op
	err = token.Check(&ops2)
	c.Check(err, gc.ErrorMatches, `"application/0" is not leader of "application"`)
	c.Check(ops2, gc.IsNil)
}

func (s *LeadershipSuite) TestCloseStateUnblocksClaimer(c *gc.C) {
	err := s.claimer.ClaimLeadership("blah", "blah/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)

	err = s.State.Close()
	c.Assert(err, jc.ErrorIsNil)

	select {
	case err := <-s.expiryChan("blah", nil):
		c.Check(err, gc.ErrorMatches, "lease manager stopped")
	case <-s.Clock.After(coretesting.LongWait):
		c.Fatalf("timed out while waiting for unblock")
	}
}

func (s *LeadershipSuite) TestLeadershipClaimerRestarts(c *gc.C) {
	// SetClockForTesting will restart the workers, and
	// will have replaced them by the time it returns.
	s.State.SetClockForTesting(gitjujutesting.NewClock(time.Time{}))

	err := s.claimer.ClaimLeadership("blah", "blah/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *LeadershipSuite) TestLeadershipCheckerRestarts(c *gc.C) {
	err := s.claimer.ClaimLeadership("application", "application/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)

	// SetClockForTesting will restart the workers, and
	// will have replaced them by the time it returns.
	s.State.SetClockForTesting(gitjujutesting.NewClock(time.Time{}))

	token := s.checker.LeadershipCheck("application", "application/0")
	err = token.Check(nil)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *LeadershipSuite) TestBlockUntilLeadershipReleasedCancel(c *gc.C) {
	err := s.claimer.ClaimLeadership("blah", "blah/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)

	cancel := make(chan struct{})
	close(cancel)

	select {
	case err := <-s.expiryChan("blah", cancel):
		c.Check(err, gc.Equals, leadership.ErrBlockCancelled)
	case <-s.Clock.After(coretesting.LongWait):
		c.Fatalf("timed out while waiting for unblock")
	}
}

func (s *LeadershipSuite) TestApplicationLeaders(c *gc.C) {
	err := s.claimer.ClaimLeadership("blah", "blah/0", time.Minute)
	c.Assert(err, jc.ErrorIsNil)
	err = s.claimer.ClaimLeadership("application", "application/1", time.Minute)
	c.Assert(err, jc.ErrorIsNil)

	leaders, err := s.State.ApplicationLeaders()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(leaders, jc.DeepEquals, map[string]string{
		"application": "application/1",
		"blah":        "blah/0",
	})
}

func (s *LeadershipSuite) expire(c *gc.C, applicationname string) {
	s.Clock.Advance(time.Hour)
	err := s.globalClock.Advance(time.Hour)
	c.Assert(err, jc.ErrorIsNil)
	s.Session.Fsync(false)
	select {
	case err := <-s.expiryChan(applicationname, nil):
		c.Assert(err, jc.ErrorIsNil)
	case <-s.Clock.After(coretesting.LongWait):
		c.Fatalf("never unblocked")
	}
}

func (s *LeadershipSuite) expiryChan(applicationname string, cancel <-chan struct{}) <-chan error {
	expired := make(chan error, 1)
	go func() {
		expired <- s.claimer.BlockUntilLeadershipReleased("blah", cancel)
	}()
	return expired
}
