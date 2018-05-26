// +build !jstesting

package engo

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSetRoot makes sure set root sets Files.root to the proper value.
func TestFilesSetRoot(t *testing.T) {
	Files.SetRoot("testing")
	if Files.root != "testing" {
		t.Errorf("Root was not set to %v, it was %v instead", "testing", Files.root)
	}
}

type testLoader struct{}

func (l *testLoader) Load(url string, data io.Reader) error {
	return nil
}

func (l *testLoader) Unload(url string) error {
	return nil
}

func (l *testLoader) Resource(url string) (Resource, error) {
	return testResource{url: url}, nil
}

type testResource struct {
	url string
}

func (r testResource) URL() string {
	return r.url
}

func TestFilesRegister(t *testing.T) {
	Files.Register(".test", &testLoader{})
	_, ok := Files.formats[".test"]
	if !ok {
		t.Error("Files.Register failed to register .test")
	}
}

func TestFilesLoad(t *testing.T) {
	Files.Register(".test", &testLoader{})

	content := []byte("testing")
	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		t.Errorf("failed to create temp directory for testing, error: %v", err)
	}
	defer os.RemoveAll(dir)

	Files.SetRoot(dir)

	tmpfn := filepath.Join(dir, "test1.test")

	if err = ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	if err = Files.Load("test1.test"); err != nil {
		t.Errorf("could not load test file %v, error: %v", "test1.test", err)
	}
}

func TestFilesMultipleLoad(t *testing.T) {
	Files.Register(".test", &testLoader{})

	content := []byte("testing")
	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		t.Errorf("failed to create temp directory for testing, error: %v", err)
	}
	defer os.RemoveAll(dir)

	Files.SetRoot(dir)

	tmpfn := filepath.Join(dir, "test1.test")

	if err = ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	tmpfn = filepath.Join(dir, "test2.test")

	if err = ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	tmpfn = filepath.Join(dir, "test3.test")

	if err = ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	if err = Files.Load("test1.test", "test2.test", "test3.test"); err != nil {
		t.Errorf("could not load test file %v, error: %v", "test1.test", err)
	}
}

func TestFilesLoadNotExist(t *testing.T) {
	Files.Register(".test", &testLoader{})

	expected := "unable to open resource:"
	if err := Files.Load("notExist.test"); err == nil {
		t.Error("did not report loading non-existant file as an error")
	} else if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("wrong error returned loading non-existant file. want: %v, got: %v", expected, err.Error())
	}
}

func TestFilesLoadNoFileLoader(t *testing.T) {
	expected := "no `FileLoader` associated with this extension:"
	if err := Files.Load("test.wrongExtension"); err == nil {
		t.Error("did not report loading file without an associated file loader")
	} else if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("wrong error returned loading file without an associated file loader. want: %v, got %v", expected, err.Error())
	}
}

func TestFilesLoadReaderData(t *testing.T) {
	Files.Register(".test", &testLoader{})

	f := bytes.NewReader([]byte("testing"))
	if err := Files.LoadReaderData("readerTest.test", f); err != nil {
		t.Errorf("unable to load from an io.Reader. error: %v", err)
	}
}

func TestFilesLoadReaderDataNoFileLoader(t *testing.T) {
	f := bytes.NewReader([]byte("testing"))
	expected := "no `FileLoader` associated with this extension:"
	if err := Files.LoadReaderData("test.wrongExtension", f); err == nil {
		t.Error("did not report loading reader without an associated file loader")
	} else if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("wrong error returned loading a reader without an associated file loader. want: %v, got: %v", expected, err.Error())
	}
}

func TestFilesUnload(t *testing.T) {
	Files.Register(".test", &testLoader{})
	if err := Files.Unload("test.test"); err != nil {
		t.Errorf("unable to unload a file. error: %v", err)
	}
}

func TestFilesUnloadNoFileLoader(t *testing.T) {
	expected := "no `FileLoader` associated with this extension:"
	if err := Files.Unload("test.wrongExtension"); err == nil {
		t.Error("did not report error unloading without an associated file loader")
	} else if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("wrong error returned unloading without an associated file loader. want: %v, got: %v", expected, err.Error())
	}
}

func TestFilesResource(t *testing.T) {
	Files.Register(".test", &testLoader{})
	if _, err := Files.Resource("test.test"); err != nil {
		t.Errorf("unable to fetch a resource. error: %v", err)
	}
}

func TestFilesResourceNoFileLoader(t *testing.T) {
	expected := "no `FileLoader` associated with this extension:"
	if _, err := Files.Resource("test.wrongExtension"); err == nil {
		t.Error("did not report error retrieving a resource without an associated file loader")
	} else if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("wrong error returned retrieving a resource without an associated file loader. want: %v, got: %v", expected, err.Error())
	}
}
