package register

import (
	"bytes"
	"encoding/json"
	"github.com/ucarion/urlpath"
	"gomin/src/framework"
	"gomin/src/tests"
	"gotest.tools/v3/assert"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterNoData(t *testing.T) {
	testDB := tests.NewTestDB(t)
	tx := testDB.Tx
	store := framework.NewSessionStore()
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	session, _ := store.Get(req, "user")

	register := Register{}
	body, _ := json.Marshal(register)
	res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ := res.Data.(RegisterError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Errors.Email != "")
	assert.Assert(t, data.Errors.Password != "")
}

func TestRegisterBadEmail(t *testing.T) {
	testDB := tests.NewTestDB(t)
	tx := testDB.Tx
	store := framework.NewSessionStore()
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	session, _ := store.Get(req, "user")

	register := Register{}
	register.Email = "test"
	register.Password = "test"
	body, _ := json.Marshal(register)
	res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
	data, _ := res.Data.(RegisterError)
	assert.Assert(t, res.StatusCode == 400)
	assert.Assert(t, data.Errors.Email != "")
	assert.Assert(t, data.Errors.Password == "")
}

func TestRegisterFlow(t *testing.T) {
	testDB := tests.NewTestDB(t)
	tx := testDB.Tx
	store := framework.NewSessionStore()
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	session, _ := store.Get(req, "user")

	register := Register{}
	register.Email = "test@example.com"
	register.Password = "test"
	body, _ := json.Marshal(register)
	res := PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
	assert.Assert(t, res.Redirect != nil)
	var count int
	err := tx.QueryRow("SELECT count(*) FROM app_user where email = $1", strings.ToLower(register.Email)).Scan(&count)
	if err != nil {
		panic(err)
	}
	assert.Assert(t, count == 1)

	// can't register twice
	res = PostHandler(tx, session)(urlpath.Match{}, bytes.NewReader(body))
	assert.Assert(t, res.StatusCode == 400)
}
