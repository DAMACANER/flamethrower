package db

import (
	"flag"
)

var DBLocation = flag.String("db", "../../dnd35.db", "location of the database")

var TestDBLocation = flag.String("test_db", "../../dnd35-test.db", "location of the test database")
