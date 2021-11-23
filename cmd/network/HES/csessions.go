package main

import (
	"net/http"
	"time"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
)

type sessionData struct {
	user *User
	ts   time.Time
}

func NewSessions() *Sessions {
	return &Sessions{
		mpn: make(map[string]*sessionData),
	}
}

func (sessions *Sessions) Set(w http.ResponseWriter, user *User) {
	sessions.mtx.Lock()
	defer sessions.mtx.Unlock()
	for k, v := range sessions.mpn {
		if v.user.Name == user.Name {
			delete(sessions.mpn, k)
			break
		}
	}
	key := cr.RandString(gp.Get("SALT_SIZE").(uint))
	sessions.mpn[key] = &sessionData{
		user: user,
		ts:   time.Now(),
	}
	createCookie(w, key)
}

func (sessions *Sessions) Get(r *http.Request) *User {
	sessions.mtx.Lock()
	defer sessions.mtx.Unlock()
	key := readCookie(r)
	if _, ok := sessions.mpn[key]; !ok {
		return nil
	}
	sessions.mpn[key].ts = time.Now()
	return sessions.mpn[key].user
}

func (sessions *Sessions) Del(w http.ResponseWriter, r *http.Request) {
	sessions.mtx.Lock()
	defer sessions.mtx.Unlock()
	delete(sessions.mpn, readCookie(r))
	deleteCookie(w)
}

func (sessions *Sessions) DelByTime(t time.Duration) {
	sessions.mtx.Lock()
	defer sessions.mtx.Unlock()
	currTime := time.Now()
	for k, v := range sessions.mpn {
		if v.ts.Add(t).Before(currTime) {
			delete(sessions.mpn, k)
		}
	}
}

func createCookie(w http.ResponseWriter, data string) {
	c := http.Cookie{
		Name:   "storage",
		Value:  data,
		MaxAge: 3600,
	}
	http.SetCookie(w, &c)
}

func readCookie(r *http.Request) string {
	c, err := r.Cookie("storage")
	value := ""
	if err == nil {
		value = c.Value
	}
	return value
}

func deleteCookie(w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "storage",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
}
