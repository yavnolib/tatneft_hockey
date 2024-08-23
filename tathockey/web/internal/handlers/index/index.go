package index

import (
	"net/http"
	"tat_hockey_pack/internal/service/session"
)

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := session.FromContext(r.Context())
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/feeds", http.StatusFound)
}
