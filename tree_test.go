package xroute

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestTree(t *testing.T) {
	hStub := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hIndex := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hFavicon := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleList := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleNear := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleShow := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleShowRelated := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleShowOpts := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleSlug := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hArticleByUser := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hUserList := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hUserShow := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hAdminCatchall := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hAdminAppShow := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hAdminAppShowCatchall := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hUserProfile := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hUserSuper := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hUserAll := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hHubView1 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hHubView2 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hHubView3 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	tr := &node{}

	tr.InsertRoute(true, GET, "/", hIndex)
	tr.InsertRoute(true, GET, "/favicon.ico", hFavicon)

	tr.InsertRoute(true, GET, "/pages/*", hStub)

	tr.InsertRoute(true, GET, "/article", hArticleList)
	tr.InsertRoute(true, GET, "/article/", hArticleList)

	tr.InsertRoute(true, GET, "/article/near", hArticleNear)
	tr.InsertRoute(true, GET, "/article/{id}", hStub)
	tr.InsertRoute(true, GET, "/article/{id}", hArticleShow)
	tr.InsertRoute(true, GET, "/article/{id}", hArticleShow) // duplicate will have no effect
	tr.InsertRoute(true, GET, "/article/@{user}", hArticleByUser)

	tr.InsertRoute(true, GET, "/article/{sup}/{opts}", hArticleShowOpts)
	tr.InsertRoute(true, GET, "/article/{id}/{opts}", hArticleShowOpts) // overwrite above route, latest wins

	tr.InsertRoute(true, GET, "/article/{iffd}/edit", hStub)
	tr.InsertRoute(true, GET, "/article/{id}//related", hArticleShowRelated)
	tr.InsertRoute(true, GET, "/article/slug/{month}/-/{day}/{year}", hArticleSlug)

	tr.InsertRoute(true, GET, "/admin/user", hUserList)
	tr.InsertRoute(true, GET, "/admin/user/", hStub) // will get replaced by next route
	tr.InsertRoute(true, GET, "/admin/user/", hUserList)

	tr.InsertRoute(true, GET, "/admin/user//{id}", hUserShow)
	tr.InsertRoute(true, GET, "/admin/user/{id}", hUserShow)

	tr.InsertRoute(true, GET, "/admin/apps/{id}", hAdminAppShow)
	tr.InsertRoute(true, GET, "/admin/apps/{id}/*ff", hAdminAppShowCatchall) // TODO: ALLOWED...? prob not.. panic..?

	tr.InsertRoute(true, GET, "/admin/*ff", hStub) // catchall segment will get replaced by next route
	tr.InsertRoute(true, GET, "/admin/*", hAdminCatchall)

	tr.InsertRoute(true, GET, "/users/{userID}/profile", hUserProfile)
	tr.InsertRoute(true, GET, "/users/super/*", hUserSuper)
	tr.InsertRoute(true, GET, "/users/*", hUserAll)

	tr.InsertRoute(true, GET, "/hubs/{hubID}/view", hHubView1)
	tr.InsertRoute(true, GET, "/hubs/{hubID}/view/*", hHubView2)
	sr := NewRouter()
	sr.Get("/users", hHubView3)
	tr.InsertRoute(true, GET, "/hubs/{hubID}/*", sr)
	tr.InsertRoute(true, GET, "/hubs/{hubID}/users", hHubView3)

	tests := []struct {
		r string         // input request path
		h ContextHandler // output matched handler
		k []string       // output param keys
		v []string       // output param values
	}{
		{r: "/", h: hIndex, k: []string{}, v: []string{}},
		{r: "/favicon.ico", h: hFavicon, k: []string{}, v: []string{}},

		{r: "/pages", h: nil, k: []string{}, v: []string{}},
		{r: "/pages/", h: hStub, k: []string{"*"}, v: []string{""}},
		{r: "/pages/yes", h: hStub, k: []string{"*"}, v: []string{"yes"}},

		{r: "/article", h: hArticleList, k: []string{}, v: []string{}},
		{r: "/article/", h: hArticleList, k: []string{}, v: []string{}},
		{r: "/article/near", h: hArticleNear, k: []string{}, v: []string{}},
		{r: "/article/neard", h: hArticleShow, k: []string{"id"}, v: []string{"neard"}},
		{r: "/article/123", h: hArticleShow, k: []string{"id"}, v: []string{"123"}},
		{r: "/article/123/456", h: hArticleShowOpts, k: []string{"id", "opts"}, v: []string{"123", "456"}},
		{r: "/article/@peter", h: hArticleByUser, k: []string{"user"}, v: []string{"peter"}},
		{r: "/article/22//related", h: hArticleShowRelated, k: []string{"id"}, v: []string{"22"}},
		{r: "/article/111/edit", h: hStub, k: []string{"iffd"}, v: []string{"111"}},
		{r: "/article/slug/sept/-/4/2015", h: hArticleSlug, k: []string{"month", "day", "year"}, v: []string{"sept", "4", "2015"}},
		{r: "/article/:id", h: hArticleShow, k: []string{"id"}, v: []string{":id"}},

		{r: "/admin/user", h: hUserList, k: []string{}, v: []string{}},
		{r: "/admin/user/", h: hUserList, k: []string{}, v: []string{}},
		{r: "/admin/user/1", h: hUserShow, k: []string{"id"}, v: []string{"1"}},
		{r: "/admin/user//1", h: hUserShow, k: []string{"id"}, v: []string{"1"}},
		{r: "/admin/hi", h: hAdminCatchall, k: []string{"*"}, v: []string{"hi"}},
		{r: "/admin/lots/of/:fun", h: hAdminCatchall, k: []string{"*"}, v: []string{"lots/of/:fun"}},
		{r: "/admin/apps/333", h: hAdminAppShow, k: []string{"id"}, v: []string{"333"}},
		{r: "/admin/apps/333/woot", h: hAdminAppShowCatchall, k: []string{"id", "*"}, v: []string{"333", "woot"}},

		{r: "/hubs/123/view", h: hHubView1, k: []string{"hubID"}, v: []string{"123"}},
		{r: "/hubs/123/view/index.html", h: hHubView2, k: []string{"hubID", "*"}, v: []string{"123", "index.html"}},
		{r: "/hubs/123/users", h: hHubView3, k: []string{"hubID"}, v: []string{"123"}},

		{r: "/users/123/profile", h: hUserProfile, k: []string{"userID"}, v: []string{"123"}},
		{r: "/users/super/123/okay/yes", h: hUserSuper, k: []string{"*"}, v: []string{"123/okay/yes"}},
		{r: "/users/123/okay/yes", h: hUserAll, k: []string{"*"}, v: []string{"123/okay/yes"}},
	}

	// log.Println("~~~~~~~~~")
	// log.Println("~~~~~~~~~")
	// debugPrintTree(0, 0, tr, 0)
	// log.Println("~~~~~~~~~")
	// log.Println("~~~~~~~~~")

	for i, tt := range tests {
		rctx := NewRouteContext()

		_, handlers, _ := tr.FindRoute(rctx, GET, tt.r)

		var handler ContextHandler
		if methodHandler, ok := handlers[GET]; ok {
			handler = methodHandler.handler.Handler(nil)
		}

		paramKeys := rctx.routeParams.Keys
		paramValues := rctx.routeParams.Values

		if fmt.Sprintf("%v", tt.h) != fmt.Sprintf("%v", handler) {
			t.Errorf("input [%d]: find '%s' expecting handler:%v , got:%v", i, tt.r, tt.h, handler)
		}
		if !stringSliceEqual(tt.k, paramKeys) {
			t.Errorf("input [%d]: find '%s' expecting paramKeys:(%d)%v , got:(%d)%v", i, tt.r, len(tt.k), tt.k, len(paramKeys), paramKeys)
		}
		if !stringSliceEqual(tt.v, paramValues) {
			t.Errorf("input [%d]: find '%s' expecting paramValues:(%d)%v , got:(%d)%v", i, tt.r, len(tt.v), tt.v, len(paramValues), paramValues)
		}
	}
}

func TestTreeMoar(t *testing.T) {
	hStub := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub1 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub2 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub3 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub4 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub5 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub6 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub7 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub8 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub9 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub10 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub11 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub12 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub13 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub14 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub15 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub16 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	// TODO: panic if we see {id}{x} because we're missing a delimiter, its not possible.
	// also {:id}* is not possible.

	tr := &node{}

	tr.InsertRoute(true, GET, "/articlefun", hStub5)
	tr.InsertRoute(true, GET, "/articles/{id}", hStub)
	tr.InsertRoute(true, DELETE, "/articles/{slug}", hStub8)
	tr.InsertRoute(true, GET, "/articles/search", hStub1)
	tr.InsertRoute(true, GET, "/articles/{id}:delete", hStub8)
	tr.InsertRoute(true, GET, "/articles/{iidd}!sup", hStub4)
	tr.InsertRoute(true, GET, "/articles/{id}:{op}", hStub3)
	tr.InsertRoute(true, GET, "/articles/{id}:{op}", hStub2)                              // this route sets a new handler for the above route
	tr.InsertRoute(true, GET, "/articles/{slug:^[a-z]+}/posts", hStub)                    // up to tail '/' will only match if contents match the rex
	tr.InsertRoute(true, GET, "/articles/{id}/posts/{pid}", hStub6)                       // /articles/123/posts/1
	tr.InsertRoute(true, GET, "/articles/{id}/posts/{month}/{day}/{year}/{slug}", hStub7) // /articles/123/posts/09/04/1984/juice
	tr.InsertRoute(true, GET, "/articles/{id}.json", hStub10)
	tr.InsertRoute(true, GET, "/articles/{id}/data.json", hStub11)
	tr.InsertRoute(true, GET, "/articles/files/{file}.{ext}", hStub12)
	tr.InsertRoute(true, PUT, "/articles/me", hStub13)

	// TODO: make a separate test case for this one..
	// tr.InsertRoute(true,GET, "/articles/{id}/{id}", hStub1)                              // panic expected, we're duplicating param keys

	tr.InsertRoute(true, GET, "/pages/*ff", hStub) // TODO: panic, allow it..?
	tr.InsertRoute(true, GET, "/pages/*", hStub9)

	tr.InsertRoute(true, GET, "/users/{id}", hStub14)
	tr.InsertRoute(true, GET, "/users/{id}/settings/{key}", hStub15)
	tr.InsertRoute(true, GET, "/users/{id}/settings/*", hStub16)

	tests := []struct {
		m MethodType     // input request http method
		r string         // input request path
		h ContextHandler // output matched handler
		k []string       // output param keys
		v []string       // output param values
	}{
		{m: GET, r: "/articles/search", h: hStub1, k: []string{}, v: []string{}},
		{m: GET, r: "/articlefun", h: hStub5, k: []string{}, v: []string{}},
		{m: GET, r: "/articles/123", h: hStub, k: []string{"id"}, v: []string{"123"}},
		{m: DELETE, r: "/articles/123mm", h: hStub8, k: []string{"slug"}, v: []string{"123mm"}},
		{m: GET, r: "/articles/789:delete", h: hStub8, k: []string{"id"}, v: []string{"789"}},
		{m: GET, r: "/articles/789!sup", h: hStub4, k: []string{"iidd"}, v: []string{"789"}},
		{m: GET, r: "/articles/123:sync", h: hStub2, k: []string{"id", "op"}, v: []string{"123", "sync"}},
		{m: GET, r: "/articles/456/posts/1", h: hStub6, k: []string{"id", "pid"}, v: []string{"456", "1"}},
		{m: GET, r: "/articles/456/posts/09/04/1984/juice", h: hStub7, k: []string{"id", "month", "day", "year", "slug"}, v: []string{"456", "09", "04", "1984", "juice"}},
		{m: GET, r: "/articles/456.json", h: hStub10, k: []string{"id"}, v: []string{"456"}},
		{m: GET, r: "/articles/456/data.json", h: hStub11, k: []string{"id"}, v: []string{"456"}},

		{m: GET, r: "/articles/files/file.zip", h: hStub12, k: []string{"file", "ext"}, v: []string{"file", "zip"}},
		{m: GET, r: "/articles/files/photos.tar.gz", h: hStub12, k: []string{"file", "ext"}, v: []string{"photos", "tar.gz"}},
		{m: GET, r: "/articles/files/photos.tar.gz", h: hStub12, k: []string{"file", "ext"}, v: []string{"photos", "tar.gz"}},

		{m: PUT, r: "/articles/me", h: hStub13, k: []string{}, v: []string{}},
		{m: GET, r: "/articles/me", h: hStub, k: []string{"id"}, v: []string{"me"}},
		{m: GET, r: "/pages", h: nil, k: []string{}, v: []string{}},
		{m: GET, r: "/pages/", h: hStub9, k: []string{"*"}, v: []string{""}},
		{m: GET, r: "/pages/yes", h: hStub9, k: []string{"*"}, v: []string{"yes"}},

		{m: GET, r: "/users/1", h: hStub14, k: []string{"id"}, v: []string{"1"}},
		{m: GET, r: "/users/", h: nil, k: []string{}, v: []string{}},
		{m: GET, r: "/users/2/settings/password", h: hStub15, k: []string{"id", "key"}, v: []string{"2", "password"}},
		{m: GET, r: "/users/2/settings/", h: hStub16, k: []string{"id", "*"}, v: []string{"2", ""}},
	}

	// log.Println("~~~~~~~~~")
	// log.Println("~~~~~~~~~")
	// debugPrintTree(0, 0, tr, 0)
	// log.Println("~~~~~~~~~")
	// log.Println("~~~~~~~~~")

	for i, tt := range tests {
		rctx := NewRouteContext()

		_, handlers, _ := tr.FindRoute(rctx, tt.m, tt.r)

		var handler ContextHandler
		if methodHandler, ok := handlers[tt.m]; ok {
			handler = methodHandler.handler.Handler(nil)
		}

		paramKeys := rctx.routeParams.Keys
		paramValues := rctx.routeParams.Values

		if fmt.Sprintf("%v", tt.h) != fmt.Sprintf("%v", handler) {
			t.Errorf("input [%d]: find '%s' expecting handler:%v , got:%v", i, tt.r, tt.h, handler)
		}
		if !stringSliceEqual(tt.k, paramKeys) {
			t.Errorf("input [%d]: find '%s' expecting paramKeys:(%d)%v , got:(%d)%v", i, tt.r, len(tt.k), tt.k, len(paramKeys), paramKeys)
		}
		if !stringSliceEqual(tt.v, paramValues) {
			t.Errorf("input [%d]: find '%s' expecting paramValues:(%d)%v , got:(%d)%v", i, tt.r, len(tt.v), tt.v, len(paramValues), paramValues)
		}
	}
}

func TestTreeRegexp(t *testing.T) {
	hStub1 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub2 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub3 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub4 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub5 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub6 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub7 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	tr := &node{}
	tr.InsertRoute(true, GET, "/articles/{rid:^[0-9]{5,6}}", hStub7)
	tr.InsertRoute(true, GET, "/articles/{zid:^0[0-9]+}", hStub3)
	tr.InsertRoute(true, GET, "/articles/{name:^@[a-z]+}/posts", hStub4)
	tr.InsertRoute(true, GET, "/articles/{op:^[0-9]+}/run", hStub5)
	tr.InsertRoute(true, GET, "/articles/{id:^[0-9]+}", hStub1)
	tr.InsertRoute(true, GET, "/articles/{id:^[1-9]+}-{aux}", hStub6)
	tr.InsertRoute(true, GET, "/articles/{slug}", hStub2)

	// log.Println("~~~~~~~~~")
	// log.Println("~~~~~~~~~")
	// debugPrintTree(0, 0, tr, 0)
	// log.Println("~~~~~~~~~")
	// log.Println("~~~~~~~~~")

	tests := []struct {
		r string         // input request path
		h ContextHandler // output matched handler
		k []string       // output param keys
		v []string       // output param values
	}{
		{r: "/articles", h: nil, k: []string{}, v: []string{}},
		{r: "/articles/12345", h: hStub7, k: []string{"rid"}, v: []string{"12345"}},
		{r: "/articles/123", h: hStub1, k: []string{"id"}, v: []string{"123"}},
		{r: "/articles/how-to-build-a-router", h: hStub2, k: []string{"slug"}, v: []string{"how-to-build-a-router"}},
		{r: "/articles/0456", h: hStub3, k: []string{"zid"}, v: []string{"0456"}},
		{r: "/articles/@pk/posts", h: hStub4, k: []string{"name"}, v: []string{"@pk"}},
		{r: "/articles/1/run", h: hStub5, k: []string{"op"}, v: []string{"1"}},
		{r: "/articles/1122", h: hStub1, k: []string{"id"}, v: []string{"1122"}},
		{r: "/articles/1122-yes", h: hStub6, k: []string{"id", "aux"}, v: []string{"1122", "yes"}},
	}

	for i, tt := range tests {
		rctx := NewRouteContext()

		_, handlers, _ := tr.FindRoute(rctx, GET, tt.r)

		var handler ContextHandler
		if methodHandler, ok := handlers[GET]; ok {
			handler = methodHandler.handler.Handler(nil)
		}

		paramKeys := rctx.routeParams.Keys
		paramValues := rctx.routeParams.Values

		if fmt.Sprintf("%v", tt.h) != fmt.Sprintf("%v", handler) {
			t.Errorf("input [%d]: find '%s' expecting handler:%v , got:%v", i, tt.r, tt.h, handler)
		}
		if !stringSliceEqual(tt.k, paramKeys) {
			t.Errorf("input [%d]: find '%s' expecting paramKeys:(%d)%v , got:(%d)%v", i, tt.r, len(tt.k), tt.k, len(paramKeys), paramKeys)
		}
		if !stringSliceEqual(tt.v, paramValues) {
			t.Errorf("input [%d]: find '%s' expecting paramValues:(%d)%v , got:(%d)%v", i, tt.r, len(tt.v), tt.v, len(paramValues), paramValues)
		}
	}
}

func TestTreeRegexMatchWholeParam(t *testing.T) {
	hStub1 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	rctx := NewRouteContext()
	tr := &node{}
	tr.InsertRoute(true, GET, "/{id:[0-9]+}", hStub1)

	tests := []struct {
		url             string
		expectedHandler ContextHandler
	}{
		{url: "/13", expectedHandler: hStub1},
		{url: "/a13", expectedHandler: nil},
		{url: "/13.jpg", expectedHandler: nil},
		{url: "/a13.jpg", expectedHandler: nil},
	}

	for _, tc := range tests {
		_, _, handler := tr.FindRoute(rctx, GET, tc.url, nil)
		if fmt.Sprintf("%v", tc.expectedHandler) != fmt.Sprintf("%v", handler) {
			t.Errorf("expecting handler:%v , got:%v", tc.expectedHandler, handler)
		}
	}
}

func TestTreeFindPattern(t *testing.T) {
	hStub1 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub2 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	hStub3 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	tr := &node{}
	tr.InsertRoute(true, GET, "/pages/*", hStub1)
	tr.InsertRoute(true, GET, "/articles/{id}/*", hStub2)
	tr.InsertRoute(true, GET, "/articles/{slug}/{uid}/*", hStub3)

	if tr.findPattern("/pages") != false {
		t.Errorf("find /pages failed")
	}
	if tr.findPattern("/pages*") != false {
		t.Errorf("find /pages* failed - should be nil")
	}
	if tr.findPattern("/pages/*") == false {
		t.Errorf("find /pages/* failed")
	}
	if tr.findPattern("/articles/{id}/*") == false {
		t.Errorf("find /articles/{id}/* failed")
	}
	if tr.findPattern("/articles/{something}/*") == false {
		t.Errorf("find /articles/{something}/* failed")
	}
	if tr.findPattern("/articles/{slug}/{uid}/*") == false {
		t.Errorf("find /articles/{slug}/{uid}/* failed")
	}
}

func debugPrintTree(parent int, i int, n *node, label byte) bool {
	numEdges := 0
	for _, nds := range n.children {
		numEdges += len(nds)
	}

	// if n.handlers != nil {
	// 	log.Printf("[node %d parent:%d] typ:%d prefix:%s label:%s tail:%s numEdges:%d isLeaf:%v handler:%v pat:%s keys:%v\n", i, parent, n.typ, n.prefix, string(label), string(n.tail), numEdges, n.isLeaf(), n.handlers, n.pattern, n.paramKeys)
	// } else {
	// 	log.Printf("[node %d parent:%d] typ:%d prefix:%s label:%s tail:%s numEdges:%d isLeaf:%v pat:%s keys:%v\n", i, parent, n.typ, n.prefix, string(label), string(n.tail), numEdges, n.isLeaf(), n.pattern, n.paramKeys)
	// }
	if n.endpoints != nil {
		log.Printf("[node %d parent:%d] typ:%d prefix:%s label:%s tail:%s numEdges:%d isLeaf:%v handler:%v\n", i, parent, n.typ, n.prefix, string(label), string(n.tail), numEdges, n.isLeaf(), n.endpoints)
	} else {
		log.Printf("[node %d parent:%d] typ:%d prefix:%s label:%s tail:%s numEdges:%d isLeaf:%v\n", i, parent, n.typ, n.prefix, string(label), string(n.tail), numEdges, n.isLeaf())
	}
	parent = i
	for _, nds := range n.children {
		for _, e := range nds {
			i++
			if debugPrintTree(parent, i, e, e.label) {
				return true
			}
		}
	}
	return false
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if b[i] != a[i] {
			return false
		}
	}
	return true
}

func BenchmarkTreeGet(b *testing.B) {
	h1 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})
	h2 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	tr := &node{}
	tr.InsertRoute(true, GET, "/", h1)
	tr.InsertRoute(true, GET, "/ping", h2)
	tr.InsertRoute(true, GET, "/pingall", h2)
	tr.InsertRoute(true, GET, "/ping/{id}", h2)
	tr.InsertRoute(true, GET, "/ping/{id}/woop", h2)
	tr.InsertRoute(true, GET, "/ping/{id}/{opt}", h2)
	tr.InsertRoute(true, GET, "/pinggggg", h2)
	tr.InsertRoute(true, GET, "/hello", h1)

	mctx := NewRouteContext()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mctx.Reset()
		tr.FindRoute(mctx, GET, "/ping/123/456")
	}
}

func TestWalker(t *testing.T) {
	r := bigMux()

	// Walk the muxBig router tree.
	if err := Walk(r, func(method string, route string, handler ContextHandler, middlewares ...*Middleware) error {
		t.Logf("%v %v", method, route)
		return nil
	}); err != nil {
		t.Error(err)
	}
}

func TestTreeGetRoute(t *testing.T) {
	hStub1 := &HTTPHandlerFunc{func(w http.ResponseWriter, r *http.Request) {}}
	hStub2 := HttpHandler(func(w http.ResponseWriter, r *http.Request) {})

	tr := &node{}
	tr.InsertRoute(true, GET, "/{id:[0-9]+}", hStub1)
	tr.InsertRoute(true, POST, "/{id:[0-9]+}", hStub2)
}
