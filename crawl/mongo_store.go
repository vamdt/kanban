package crawl

import "gopkg.in/mgo.v2"

func NewMongoStore(dsn string) (ms *MongoStore, err error) {
	session, err := mgo.Dial(dsn)
	if err != nil {
		return
	}
	ms = &MongoStore{session: session}
	return
}

type MongoStore struct {
	session *mgo.Session
}

func (p *MongoStore) Close() {
	p.session.Close()
}
