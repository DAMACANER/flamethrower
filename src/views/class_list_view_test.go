package views_test

import (
	"flamethrower/src/db"
	"flamethrower/src/views"
	"fmt"
	"log"
	"testing"

	"github.com/blockloop/scan"
	"github.com/stretchr/testify/suite"
)

type ClassListViewTestSuite struct {
	suite.Suite
	repo *db.ClassListRepo
}

func (suite *ClassListViewTestSuite) SetupTest() {
	db.InitDB(*db.TestDBLocation)
	suite.repo = &db.ClassListRepo{BaseRepo: &db.BaseRepo{}}
	suite.repo.DropTableIfExists()
	suite.repo.CreateTable()
	suite.repo.PopulateTable(10)
	err := db.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

func (suite *ClassListViewTestSuite) TearDownTest() {
	suite.repo.DropTableIfExists()

	db.DB.Close()
}

// TestClassListPagination is testing the custom wrappers of repo
//
//	suite.repo.Find().Paginate(currentPageNumber, pageSize).Query()
//
// to see if pagination works correctly by asserting the retrieved  page size with wrappers
// is equal to the default page size.
func (suite *ClassListViewTestSuite) TestClassListPagination() {
	var pageSize uint64 = views.DefaultPageSize
	var currentPageNumber uint64 = views.DefaultStartingPageNumber
	rows, err := suite.repo.Find().Paginate(currentPageNumber, pageSize).Query()
	if err != nil {
		log.Fatal(err)
	}
	var data []db.ClassListColumns
	err = scan.Rows(&data, rows)
	if err != nil {
		log.Fatal(err)
	}
	suite.Equal(int(pageSize), len(data))
}

// TestClassListCount is testing the custom wrappers of repo
//
//	suite.repo.Find().Count().Query()
//
// to see if they match with the raw SQL query.
func (suite *ClassListViewTestSuite) TestClassListCount() {
	dbIndirectQueryCountRows, err := suite.repo.Find().Count().Query()
	if err != nil {
		log.Fatal(err)
	}
	var dbIndirectQueryCount int
	err = scan.Row(&dbIndirectQueryCount, dbIndirectQueryCountRows)
	if err != nil {
		log.Fatal(err)
	}
	dbDirectQueryCountRows, err := db.DB.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", db.ClassListTableName))
	if err != nil {
		log.Fatal(err)
	}
	var dbDirectQueryCount int
	err = scan.Row(&dbDirectQueryCount, dbDirectQueryCountRows)
	if err != nil {
		log.Fatal(err)
	}
	suite.Equal(dbIndirectQueryCount, int(dbDirectQueryCount))
}
func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(ClassListViewTestSuite))
}
