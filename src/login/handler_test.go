package login

import (
	"bytes"
	"encoding/json"
	"github.com/ucarion/urlpath"
	"gomin/src/framework"
	"gomin/src/register"
	"gotest.tools/v3/assert"
	"net/http"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	db := framework.NewDatabaseConnection()
	defer db.Close()
	store := framework.NewSessionStore()
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	session, _ := store.Get(req, "user")

	login := Login{}
	body, _ := json.Marshal(login)
	res := PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ := res.Data.(LoginError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Errors.Email != "")

	login.Email = "blah"
	body, _ = json.Marshal(login)
	res = PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ = res.Data.(LoginError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Errors.Email == "")
	assert.Assert(t, data.Errors.Password != "")

	login.Password = "blah"
	body, _ = json.Marshal(login)
	res = PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ = res.Data.(LoginError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Email != "")
	assert.Assert(t, data.Errors.Email != "")
	assert.Assert(t, data.Errors.Password == "")

	// Register before testing login
	user := register.Register{}
	user.Email = "test@example.com"
	user.Password = "test"
	body, _ = json.Marshal(user)
	res = register.PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	assert.Assert(t, res.Redirect != nil, res.Redirect)

	// Test good login
	login.Email = "test@example.com"
	login.Password = "test"
	body, _ = json.Marshal(login)
	res = PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	assert.Assert(t, res.Redirect != nil, res.Redirect)

	db.Exec("TRUNCATE app_user;")
}
