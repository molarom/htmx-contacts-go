package flash

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "flashes"

func store() sessions.Store {
	return sessions.NewCookieStore([]byte("secretkeysflashes"))
}

func Add(w http.ResponseWriter, r *http.Request, value string) {
	session, _ := store().Get(r, sessionName)
	session.AddFlash(value)
	session.Save(r, w)
}

func Messages(w http.ResponseWriter, r *http.Request) []string {
	session, _ := store().Get(r, sessionName)

	m := session.Flashes()
	if len(m) > 0 {
		sessions.Save(r, w)
		sl := make([]string, 0, len(m))

		for _, v := range m {
			sl = append(sl, v.(string))
		}

		return sl
	}

	return nil
}
