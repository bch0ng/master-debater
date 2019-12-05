package db

import (
  "fmt"
  "testing"
  //"regexp"
  "database/sql"
  "github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-sqlmock"
  "github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
 "github.com/stretchr/testify/suite"

)
type Suite struct {
   suite.Suite
   DB   *gorm.DB
   mock sqlmock.Sqlmock
   store users.Store
   user     *users.User
}
func (s *Suite) SetupSuite() {
   var (
      db  *sql.DB
      err error
   )
   //dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/psql_db", "your-password")
   db, s.mock, err = sqlmock.New()
   //db, s.mock, err = sqlmock.NewWithDSN(dsn)
   require.NoError(s.T(), err)

   s.DB, err = gorm.Open("mysql", db)
   require.NoError(s.T(), err)

   s.DB.LogMode(true)
   dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/psql_db", "your-password")
   s.store,err = ConnectToPostgres(dsn)
   require.NoError(s.T(), err)
}
func (s *Suite) AfterTest(_, _ string) {
   require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
func (s *Suite)TestSelectByUsername() {
         res, err := s.store.GetByUserName("leoTran")
         require.NoError(s.T(), err)
         require.Equal(s.T(), res.UserName,"leoTran","Username should be")
}
func (s *Suite)TestInsertDelete() {
        newUser:=users.User{ID:1,
          Email:"testMail0",
          PassHash:[]byte{3,4,1},
          UserName:"l0",
          FirstName:"first",
          LastName:"last",
          PhotoURL:"sssss",
        }
         res, err := s.store.Insert(&newUser)
         res, err2 := s.store.GetByEmail("testMail0")
         require.NoError(s.T(), err)
         require.NoError(s.T(), err2)
         require.Equal(s.T(), "l0", res.UserName,"Fetched profile from a insert should match")
         err3 := s.store.Delete(res.ID)

         require.NoError(s.T(), err3)

}
func (s *Suite)TestInsertUpdateDelete() {
        newUser:=users.User{ID:7,
          Email:"tUser36",
          PassHash:[]byte{3,4,1},
          UserName:"l33333",
          FirstName:"first",
          LastName:"last",
          PhotoURL:"sssss",
        }
        newUpdate:=users.Updates{
          FirstName:"newName",
          LastName:"newLastName",
        }
         res, err := s.store.Insert(&newUser)
         res, err3 := s.store.GetByEmail("tUser36")
         res, err2 := s.store.Update(res.ID,&newUpdate)
         res2, err3 := s.store.GetByEmail("tUser36")
         require.NoError(s.T(), err)
         require.NoError(s.T(), err2)
         require.NoError(s.T(), err3)
         require.Equal(s.T(), "newLastName", res2.LastName,"Updated profile name should match")
         err4 := s.store.Delete(res.ID)
         require.NoError(s.T(), err4)

}
func (s *Suite)TestSelectById() {
  //rowNamesStrArr:=[]string{"id", "email", "passhash", "user_name", "first_name", "last_name", "photo_url"}
  //WithArgs will match given expected args to actual database exec
  // operation arguments. if at least one argument does not match,
  //it will return an error. For specific arguments an sqlmock.Argument
  //interface can be used to match an argument.
  /*s.mock.ExpectQuery(regexp.QuoteMeta(
      `SELECT * FROM contacts WHERE id = $1 limit 1`)).
        WithArgs(1).
        WillReturnRows(sqlmock.NewRows(rowNamesStrArr).
           AddRow(1, "mail", "passhash", "leoTran", "Leo", "Tran", "sss"))
*/
         res, err := s.store.GetByID(1)
         require.NoError(s.T(), err)
         require.Equal(s.T(), res.ID,int64(1),"Id should be one")
}
func (s *Suite)TestSelectByEmail() {
  //rowNamesStrArr:=[]string{"id", "email", "passhash", "user_name", "first_name", "last_name", "photo_url"}
  //WithArgs will match given expected args to actual database exec
  // operation arguments. if at least one argument does not match,
  //it will return an error. For specific arguments an sqlmock.Argument
  //interface can be used to match an argument.
  /*s.mock.ExpectQuery(regexp.QuoteMeta(
      `SELECT * FROM "contacts" WHERE email = $1 limit 1`)).
      WithArgs("mail").
        WillReturnRows(sqlmock.NewRows(rowNamesStrArr).
           AddRow(1, "mail", "passhash", "leoTran", "Leo", "Tran", "sss"))
*/
         res, err := s.store.GetByEmail("mail")
         if err!=nil{
           fmt.Printf("Error!!! getting mail test error.")
         }else{
           fmt.Printf("Results of mail test:%v\n",res)
           fmt.Printf("Email:%v\n",res.Email)
         }
         require.NoError(s.T(), err)
         require.Equal(s.T(), res.Email,"mail","Email should be mail")
}
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSelectSuite(t *testing.T) {
    suite.Run(t, new(Suite))
}
