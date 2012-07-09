package mstate_test

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
	state "launchpad.net/juju-core/mstate"
	"sort"
)

// ConnSuite facilitates access to the underlying MongoDB. It is embedded
// in other suites, like StateSuite.
type ConnSuite struct {
	MgoSuite
	session  *mgo.Session
	charms   *mgo.Collection
	machines *mgo.Collection
	services *mgo.Collection
	units    *mgo.Collection
}

func (cs *ConnSuite) SetUpTest(c *C) {
	cs.MgoSuite.SetUpTest(c)
	session, err := mgo.Dial(mgoaddr)
	c.Assert(err, IsNil)
	cs.session = session
	cs.charms = session.DB("juju").C("charms")
	cs.machines = session.DB("juju").C("machines")
	cs.services = session.DB("juju").C("services")
	cs.units = session.DB("juju").C("units")
}

func (cs *ConnSuite) TearDownTest(c *C) {
	cs.session.Close()
	cs.MgoSuite.TearDownTest(c)
}

func (s *ConnSuite) AllMachines(c *C) []string {
	docs := []state.MachineDoc{}
	err := s.machines.Find(bson.D{{"life", state.Alive}}).All(&docs)
	c.Assert(err, IsNil)
	names := []string{}
	for _, v := range docs {
		names = append(names, v.String())
	}
	sort.Strings(names)
	return names
}

type UtilSuite struct {
	MgoSuite
	ConnSuite
	State *state.State
}

func (s *UtilSuite) SetUpTest(c *C) {
	s.ConnSuite.SetUpTest(c)
	st, err := state.Dial(mgoaddr)
	c.Assert(err, IsNil)
	s.State = st
}

func (s *UtilSuite) TearDownTest(c *C) {
	s.State.Close()
	s.ConnSuite.TearDownTest(c)
}
