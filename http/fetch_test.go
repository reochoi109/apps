package http

import (
	"os"
	"testing"
)

func TestFetchRemoteResource(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	expected := "Hello world"
	data, err := fetchRemoteResource(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(data) {
		t.Errorf("Expected response to be : %s, Got: %s", expected, data)
	}
}

func TestFetchPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	packages, err := fetchPackageData(ts.URL + "/packages")

	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 {
		t.Fatalf("Expected 2 packages, Got back: %d", len(packages))
	}
}

func TestDownloadToFile(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	tmpFile := "test_output.json"
	defer os.Remove(tmpFile)

	if err := downloadToFile(ts.URL+"/packages", tmpFile); err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	if info.Size() == 0 {
		t.Fatal("Downloaded file is empty")
	}
}

func TestRegisterPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	url := ts.URL + "/packages"
	p := pkgData{Name: "mypackage", Version: "0.1"}

	resp, err := registerPackageData(url, p)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Id != "mypackage-0.1" {
		t.Errorf("Expected package id to be mypackage-0.1, Got: %s", resp.Id)
	}
}

func TestRegisterEmptyPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	p := pkgData{}
	_, err := registerPackageData(ts.URL+"/packages", p)

	if err == nil {
		t.Fatal("Expected an error due to invalid package data, but got nil")
	}
}
