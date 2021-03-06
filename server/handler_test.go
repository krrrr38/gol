package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGolServerGet(t *testing.T) {
	testFile := tempTest("get")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("GET", "/k1/test/test", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	// Check the response body is what we expect.
	expected := `http://test/v1/test/test`
	if rr.Header().Get("Location") != expected {
		t.Errorf("handler returned unexpected location header: got %v", rr.Header().Get("Location"))
	}
}

func TestGolServerGetNotFound(t *testing.T) {
	testFile := tempTest("get_not_found")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("GET", "/k_not_found/test/test", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGolServerDump(t *testing.T) {
	testFile := tempTest("dump")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("GET", "/api/dump", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	expected := `{
	"k1": {
		"Link": "http://test/v1",
		"RegisteredBy": "",
		"CreatedAt": "0001-01-01T00:00:00Z"
	},
	"k2": {
		"Link": "http://v2",
		"RegisteredBy": "",
		"CreatedAt": "0001-01-01T00:00:00Z"
	}
}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body %v", rr.Body.String())
	}

}

func TestGolServerDumpAsHTML(t *testing.T) {
	testFile := tempTest("dumpAsHTML")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	if rr.HeaderMap.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("handler returned unexpected content type %v", rr.HeaderMap.Get("Content-Type"))
	}

}

func TestGolServerPut(t *testing.T) {
	testURL := "http://test/v1"

	testFile := tempTest("put")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("PUT", "/newItem", strings.NewReader(fmt.Sprintf("value=%s", testURL)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	req, err = http.NewRequest("GET", "/newItem", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	// Check the response body is what we expect.
	if rr.Header().Get("Location") != testURL {
		t.Errorf("handler returned unexpected location header: got %v", rr.Header().Get("Location"))
	}
}

func TestGolServerCannotPut(t *testing.T) {
	testFile := tempTest("put")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("PUT", "/k1", strings.NewReader("value=http://test/v1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGolServerReplaceSameElementWhenPut(t *testing.T) {
	testURL := "http://test/v1"

	testFile := tempTest("put")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("PUT", "/k1", strings.NewReader(fmt.Sprintf("value=%s&force=true", testURL)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	req, err = http.NewRequest("GET", "/k1", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	// Check the response body is what we expect.
	if rr.Header().Get("Location") != testURL {
		t.Errorf("handler returned unexpected location header: got %v", rr.Header().Get("Location"))
	}
}

func TestGolServerReplaceSameElementWhenPutWithForceFlagAndRegisteredBy(t *testing.T) {
	testURL := "http://test/v1"
	testRegisteredBy := "foo"

	testFile := tempTest("put")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("PUT", "/registered_test", strings.NewReader(fmt.Sprintf("value=%s&force=true&registeredBy=%s", testURL, testRegisteredBy)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	req, err = http.NewRequest("PUT", "/registered_test", strings.NewReader(fmt.Sprintf("value=%s&force=true&registeredBy=%s", testURL, testRegisteredBy)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	req, err = http.NewRequest("PUT", "/registered_test", strings.NewReader(fmt.Sprintf("value=%s&force=true&registeredBy=%s", testURL, "wrong_registered_by")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGolServerCannotPutWithKeyWhichContainsSlash(t *testing.T) {
	testFile := tempTest("put")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("PUT", "/k1/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGolServerDelete(t *testing.T) {
	testFile := tempTest("delete")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("DELETE", "/k1", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	req, err = http.NewRequest("GET", "/k1/test/test", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGolServerDeleteShouldNotContainSlashInKey(t *testing.T) {
	testFile := tempTest("delete")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("DELETE", "/k1/foo", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGolServerUnSupportedMethod(t *testing.T) {
	testFile := tempTest("unsupported")
	defer os.Remove(testFile)
	initDb(testFile)
	handler := http.HandlerFunc(NewGolServerHandler(testFile))

	req, err := http.NewRequest("OPTIONS", "/k1/test/test", nil)
	if err != nil {
		t.Errorf("create request failed %s", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}
