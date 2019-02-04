package keymagic

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/database"
	"github.com/quan-to/remote-signer/models"
	"io/ioutil"
	"testing"
)

func TestPKSGetKey(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Internal
	c := database.GetConnection()

	z, err := ioutil.ReadFile("../tests/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	gpgKey := models.AsciiArmored2GPGKey(string(z))

	_, _, err = models.AddGPGKey(c, gpgKey)
	if err != nil {
		t.Errorf("Fail to add key to database: %s", err)
		t.FailNow()
	}

	key, _ := PKSGetKey(gpgKey.FullFingerPrint)

	fp, err := remote_signer.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(gpgKey.FullFingerPrint, fp) {
		t.Errorf("Expected %s got %s", gpgKey.FullFingerPrint, fp)
	}

	// Test External
	remote_signer.EnableRethinkSKS = false
	remote_signer.SKSServer = "https://keyserver.ubuntu.com/"

	key, _ = PKSGetKey(remote_signer.ExternalKeyFingerprint)

	fp, err = remote_signer.GetFingerPrintFromKey(key)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(remote_signer.ExternalKeyFingerprint, fp) {
		t.Errorf("Expected %s got %s", remote_signer.ExternalKeyFingerprint, fp)
	}
}

func TestPKSSearchByName(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = PKSSearchByName("", 0, 1)
	}, "SearchByName without RethinkSKS Should panic!")
}

func TestPKSSearchByFingerPrint(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = PKSSearchByFingerPrint("", 0, 1)
	}, "SearchByFingerPrint without RethinkSKS Should panic!")
}

func TestPKSSearchByEmail(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()

	// Test Panics
	remote_signer.EnableRethinkSKS = false
	assertPanic(t, func() {
		_ = PKSSearchByEmail("", 0, 1)
	}, "SearchByEmail without RethinkSKS Should panic!")
}

func TestPKSSearch(t *testing.T) {
	// TODO: Implement method and test
	// For now, should always panic

	assertPanic(t, func() {
		_ = PKSSearch("", 0, 1)
	}, "Search should always panic (NOT IMPLEMENTED)")
}

func TestPKSAdd(t *testing.T) {
	remote_signer.PushVariables()
	defer remote_signer.PopVariables()
	// Test Internal
	z, err := ioutil.ReadFile("../tests/testkey_privateTestKey.gpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := remote_signer.GetFingerPrintFromKey(string(z))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	o := PKSAdd(string(z))

	if o != "OK" {
		t.Errorf("Expected %s got %s", "OK", o)
	}

	p, _ := PKSGetKey(fp)

	if p == "" {
		t.Errorf("Key was not found")
		t.FailNow()
	}

	fp2, err := remote_signer.GetFingerPrintFromKey(string(p))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !remote_signer.CompareFingerPrint(fp, fp2) {
		t.Errorf("FingerPrint does not match. Expected %s got %s", fp, fp2)
	}

	// Test External
	remote_signer.EnableRethinkSKS = false
	// TODO: How to be a good test without stuffying SKS?
}

func assertPanic(t *testing.T, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}
