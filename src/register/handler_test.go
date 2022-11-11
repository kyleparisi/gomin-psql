package register

import (
	"bytes"
	"encoding/json"
	"github.com/ucarion/urlpath"
	"gomin/src/framework"
	"gotest.tools/v3/assert"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterPostHandler(t *testing.T) {
	db := framework.NewDatabaseConnection()
	defer db.Close()
	store := framework.NewSessionStore()
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	session, _ := store.Get(req, "user")

	register := Register{}
	body, _ := json.Marshal(register)
	res := PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ := res.Data.(RegisterError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Errors.Email != "")
	assert.Assert(t, data.Errors.Password != "")

	register.Email = "test"
	register.Password = "test"
	body, _ = json.Marshal(register)
	res = PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ = res.Data.(RegisterError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Errors.Email != "")
	assert.Assert(t, data.Errors.Password == "")

	register.Email = "test@example.com"
	body, _ = json.Marshal(register)
	res = PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	assert.Assert(t, res.Redirect != nil)
	var count int
	err := db.QueryRow("SELECT count(*) FROM app_user where email = $1", strings.ToLower(register.Email)).Scan(&count)
	if err != nil {
		panic(err)
	}
	assert.Assert(t, count == 1)

	// can't register twice
	res = PostHandler(db, session)(urlpath.Match{}, bytes.NewReader(body))
	assert.Assert(t, res.StatusCode == 400)

	db.Exec("TRUNCATE app_user;")
}
