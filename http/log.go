package http

/*
type Stat struct {
	Method  string
	Path    string
	UA      string
	Code    int
	Elapsed time.Duration
}
*/

/*
func defaultMiddleware(h http.Handler, f func(s Stat)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		c := r.Context()
		c = ctx.WithTag(c, "ua", r.UserAgent())
		c = ctx.WithTag(c, "path", r.URL.Path)
		w2 := &response{w, 0}
		switch w := w.(type) {
		case http.Hijacker:
			h.ServeHTTP(responseHijacker{w2, w}, r)
		default:
			h.ServeHTTP(w2, r)
		}
		f(Stat{
			Path:    r.URL.Path,
			UA:      r.UserAgent(),
			Code:    w2.code,
			Elapsed: time.Since(t0),
		})
	})
}
*/
