package mux

import (
	"testing"
)

func helloHandler(ctx *Context) {
	ctx.Write([]byte("Hello world"))
}

func testReverse(t *testing.T, expected string, m *Mux, name string, args ...interface{}) {
	rev, err := m.Reverse(name, args...)
	if expected != "" {
		if err != nil {
			t.Error(err)
		}
	} else {
		if err == nil {
			t.Errorf("Expecting error while reversing %s with arguments %v", name, args)
		}
	}
	if rev != expected {
		t.Errorf("Error reversing %q with arguments %v, expected %q, got %q", name, args, expected, rev)
	} else {
		t.Logf("Reversed %q with %v to %q", name, args, rev)
	}
}

func TestReverse(t *testing.T) {
	m := New()
	m.HandleNamedFunc("^/program/(\\d+)/$", helloHandler, "program")
	m.HandleNamedFunc("^/program/(\\d+)/version/(\\d+)/$", helloHandler, "programversion")
	m.HandleNamedFunc("^/program/(?P<pid>\\d+)/version/(?P<vers>\\d+)/$", helloHandler, "programversionnamed")
	m.HandleNamedFunc("^/program/(\\d+)/(?:version/(\\d+)/)?$", helloHandler, "programoptversion")
	m.HandleNamedFunc("^/program/(\\d+)/(?:version/(\\d+)/)?(?:revision/(\\d+)/)?$", helloHandler, "programrevision")
	m.HandleNamedFunc("^/archive/(\\d+)?$", helloHandler, "archive")
	m.HandleNamedFunc("^/image/(\\w+)\\.(\\w+)$", helloHandler, "image")
	m.HandleNamedFunc("^/image/(\\w+)\\-(\\w+)$", helloHandler, "imagedash")
	m.HandleNamedFunc("^/image/(\\w+)\\\\(\\w+)$", helloHandler, "imageslash")

	testReverse(t, "/program/1/", m, "program", 1)
	testReverse(t, "/program/1/version/2/", m, "programversion", 1, 2)
	testReverse(t, "/program/1/version/2/", m, "programversionnamed", 1, 2)
	testReverse(t, "/program/1/", m, "programoptversion", 1)
	testReverse(t, "/program/1/version/2/", m, "programoptversion", 1, 2)
	testReverse(t, "/program/1/", m, "programrevision", 1)
	testReverse(t, "/program/1/version/2/", m, "programrevision", 1, 2)
	testReverse(t, "/program/1/version/2/revision/3/", m, "programrevision", 1, 2, 3)

	testReverse(t, "/archive/19700101", m, "archive", "19700101")
	testReverse(t, "/archive/", m, "archive")

	// TODO: These don't work
	/*
		m.HandleNamedFunc("^/section/(sub/(\\d+)/subsub(\\d+))?$", helloHandler, "section")
		testReverse(t, "/section/", m, "section")
		testReverse(t, "/section/sub/1/subsub/2", m, "section", 1, 2)
		testReverse(t, "/section/sub/1", m, "section", 1)
	*/

	// Test invalid reverses
	testReverse(t, "", m, "program")
	testReverse(t, "", m, "program", "foo")
	testReverse(t, "", m, "program", 1, 2)
	testReverse(t, "", m, "programrevision", 1, 2, 3, 4)

	// Dot, dash and slash
	testReverse(t, "/image/test.png", m, "image", "test", "png")
	testReverse(t, "/image/test-png", m, "imagedash", "test", "png")
	testReverse(t, "/image/test\\png", m, "imageslash", "test", "png")
}