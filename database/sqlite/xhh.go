package sqlite

import "fmt"

const xhh_cookie_db = "cookies"

func xhh_savecookie(uid int, pkey, token string) (err error) {
	query := fmt.Sprintf("INSERT INTO %s (uid,pkey,token) VALUES (?,?,?) ON CONFLICT(uid) DO UPDATE SET uid = excluded.uid,pkey = excluded.pkey,token=excluded.token", xhh_cookie_db)
	_, err = db.Exec(query, uid, pkey, token)
	return
}

func xhh_getcookie(uid int) (err error, pkey, token string) {
	query := fmt.Sprintf("SELECT pkey,token FROM %s WHERE uid=?", xhh_cookie_db)
	row := db.QueryRow(query, uid)
	err = row.Err()
	if err != nil {
		return
	}
	err = row.Scan(&pkey, &token)
	return
}
