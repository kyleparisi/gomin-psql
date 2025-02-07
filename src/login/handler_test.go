package login

import (
	"bytes"
	"encoding/json"
	"github.com/ucarion/urlpath"
	"golang.org/x/crypto/bcrypt"
	"gomin/.gen/tron_local/public/model"
	. "gomin/.gen/tron_local/public/table"
	"gomin/src/framework"
	"gomin/src/tests"
	"gotest.tools/v3/assert"
	"net/http"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	testDB := tests.NewTestDB(t)
	tx := testDB.Tx
	store := framework.NewSessionStore()
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	session, _ := store.Get(req, "user")

	t.Run("no data", func(t *testing.T) {
		login := Login{}
		body, _ := json.Marshal(login)
		res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
		data, _ := res.Data.(LoginError)
		assert.Assert(t, res.StatusCode == 400)
		assert.Assert(t, data.Errors.Email != "")
	})

	t.Run("no password", func(t *testing.T) {
		login := Login{}
		login.Email = "blah"
		body, _ := json.Marshal(login)
		res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
		data, _ := res.Data.(LoginError)
		assert.Assert(t, res.StatusCode == 400)
		assert.Assert(t, data.Errors.Email == "")
		assert.Assert(t, data.Errors.Password != "")
	})

	t.Run("bad email", func(t *testing.T) {
		login := Login{}
		login.Email = "blah"
		login.Password = "blah"
		body, _ := json.Marshal(login)
		res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
		data, _ := res.Data.(LoginError)
		assert.Assert(t, res.StatusCode == 400)
		assert.Assert(t, data.Email != "")
		assert.Assert(t, data.Errors.Email != "")
		assert.Assert(t, data.Errors.Password == "")
	})

	t.Run("flow", func(t *testing.T) {
		hash, hashError := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
		if hashError != nil {
			panic(hashError)
		}
		user := model.AppUser{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Frank",
			Email:     "test@example.com",
			Password:  string(hash[:]),
		}
		insertStmt := AppUser.INSERT(AppUser.CreatedAt, AppUser.UpdatedAt, AppUser.Name, AppUser.Email, AppUser.Password).MODEL(user)
		_, err := insertStmt.Exec(tx)
		if err != nil {
			panic(err)
		}
		login := Login{}
		login.Email = "test@example.com"
		login.Password = "test"
		body, _ := json.Marshal(login)
		res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
		assert.Assert(t, res.Redirect != nil, res.Redirect)
	})

}
