package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/etc/magicBuilder"
	"github.com/quan-to/remote-signer/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
	"testing"
)

const testKeyFingerprint = "0016A9CA870AFA59"
const testKeyPassword = "I think you will never guess"

const testSignatureData = "huebr for the win!"
const testSignatureSignature = `-----BEGIN PGP SIGNATURE-----

wsFcBAABCgAQBQJcJMPWCRAAFqnKhwr6WQAA3kwQAB6pxQkN+5yMt0LSkpIcjeOS
UPqcMabEQlkD2HQrzisXlUZgllqP4jYAjFLCeErt0uu598LXO6pNTw7MnFQSfgcJ
dJF2S05GwI4k00mMNzCTn7PbJe3d96QwjbTeanoMAjHhypZKi/StbtkFpIa+t9WI
zm+EE5trFdZoE1SMOr5j85afDecl0DsGHEkKdmJ2mLK4ja3uaxsijtLd8d7mdI+Y
LbI8UnpGyWMLkK8FpjBm+BaVeNicUvqkt/LO3LwslbKAViKpdL6Gu5x7x6Q+tAyO
PZ6P6DQKjuGJl8aSv0eoKQ1TQz6vasBZNsYlasU0fM6dXny9XIucUD5sTsUpbMhw
uO/xap6i3mBtFpzSfQCo/23KHeQajXS23Al56iUr85jlSQ9+JvJhZFrU9NQa+ypq
Xi/IxrqTTvttVurXAVME1m06JirpiuD8fDdQTTboekaqLg8rXQ5eKqW0pAMIqHvf
aq97YCqxH4F3T2EE77v6D9iLnbx/+7EGHoCehTMUYiAIAhlo93Xf/hnj40Hl/N18
gYr2Yd/IYVsAoGH6AHrIyUykXgsK6RXiBy0Sa7LN14TMCnQYzG2AUvXCDf184YAQ
1obsUVANy+qxH4lwMbEoznEsAU0ppqLchX1Ixdru5/SEgSV13Qv34rMEHCdVy4Oe
1Jcr1AyB3KmDhw76PaBh
=D//n
-----END PGP SIGNATURE-----`

const testDecryptDataAscii = `-----BEGIN PGP MESSAGE-----
Version: GnuPG v2
Comment: Generated by Quanto Remote Signer

wcFMA6uJF6HKi88OARAADv3z9DhWy8UI/yf09QecNw//3foPoh3I1ZCEaywEXfB6
qvYTIHCdE0aKa/lMVVjv94YGog+EetedZZ6Ow6mBz0EjzQCemjgoHnq+69kf8y9V
6oUwBUUcZz6uxJsUA2y3mfvFXYvA6CWWk9H/RoJzwHO9Px7CNZHIWaAPcPP9bJkf
9VAISjWgCEHnH4O5uavrOqBaHwgtKvb/Ya3Cq/NDlpWcwxOHcBPxyml9Tcs6HA+7
dqnr3qhjeDzGYuyXQSnNI+ut37mCbISC4jniljsSWy5l1YkD/JfpMyMpN83BPbpJ
DP55rMBonzCQ+iV9Wyt0zUUrExuArjqyU/DKx1ZmoKWEv4EU6+BjutVZbc+sQIkk
lxP0E3bMLn73qm4JU6A1A9WqFd+ndP8hxSPb7EnqwvQ6A05NLl44kZ0JfULO4+a7
dhDPPlyGeur09Y0JZiA6k+uF1+dug52E5iW6ohBhki9SNG2Y6m1wez1gLDj7OBbp
upUEH488XKGNGN96DxlQxq2ujfozpMiRXN6IYy8ZqWskVInD5GNRw9n0BNomG1p4
abKUhK+YZR/B8rJCqm3wTGvCc5hrmnEcj944oNaSvWzfbAH81bl/8/as1IL616Hq
rOvKBs2YHJT51yw9U7ShzJTH/6GLuCaViq4d8Txi2a9JEpn3VOXv65ZiyLQSHrzS
4AHkXrdO3DB3CYC/qPoKyU0Z0OH2kOBn4GvhjZ/gK+MVrSDsHhivyODC4PLgauLN
9xqB4GDksZWWXn9uD3YAASYNjI95S+DE4Rd04C/kg0GzuOAMI1yFpV1OK12+heKg
nT5O4fndAA==
=FQ83
-----END PGP MESSAGE-----`

const testDecryptDataOnly = "wcFMA6uJF6HKi88OARAATyFPVauyY3PKircZ3AlTMd2Iy1/FNKVxSKg1jKBhGvPCUdRRMqaJXz4dsEWNZp//QQMN3cd3JqJhw/AEGJJUQglwnXO2bYCXj6/RzsgRKCbj9Ijo1Y33Rbu+3+huWluYEQnWBfkbhnjeIrNRXxGvqQKczXx1aA6D1CvFk8W5LUWmIngxKi+s2TxA/fqfMBETnKa6rVM625by/9Ebo7qKoeetksDYAMvEzaLwKFIQ6O+lQt1YcBnbZ3mkrtSosisRqfmndkffUGsEJJ/g16ZYhwlDUYUGj7O/mRb01edPFLQko0THpAUhT7GH4Cw939W4wqddHSxgz9pEJKt8TsOqry2oiRQ+Qus5ygyMrLp5jH6JrExgGf5dlNUOs6R1JXozWhLXSZo7+kBg4hTRkRmdSu2adNvsO8tF2qjCWd/M0p2HfLEKTdvYFh6+d93wOVDYMvXzUB7NDIGlhi6gXs/D8+Tw+ZkLpRm9iWdLO9YpFquI+964sxAz5E5iEOBipJGTVpyxsU959kQ6hJqT4EWiATYMnqpnG7hGkDfXlKcwCeDBDKsUi3KVLw4PbSnRZQh9JNFv2VhyF2zQKpXI4hyF0QZoef2OT4a7xTOdHkdkAes/fcDhr4dwvQfp0uPOH1C8LViO7bVKBFnj+zTVftI0pVJ8MV/BV0Y5ru1hcRXOp9vS4AHk3apKw1UvqtPqAcgnbNy1euHYYuBt4HXhvbrgWuO+bpowlkmIXeAx4CrgJOKtLADx4PPk454/jrek18yYVGg4AZvEDOBu4b0E4EzkNh+m6OARX0nP/ig4wqHxtuLvNoCX4WOWAA=="

var sm etc.SMInterface
var gpg etc.PGPInterface
var slog = SLog.Scope("TestRemoteSigner")

var router *mux.Router

func errorDie(err error, t *testing.T) {
	if err != nil {
		fmt.Println("----------------------------------------")
		debug.PrintStack()
		fmt.Println("----------------------------------------")
		t.Error(err)
		t.FailNow()
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func TestMain(m *testing.M) {
	QuantoError.EnableStackTrace()
	SLog.SetTestMode()

	remote_signer.DatabaseName = "qrs_test"
	remote_signer.PrivateKeyFolder = "../tests"
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false
	remote_signer.EnableRethinkSKS = false

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../tests/testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../tests/testprivatekeyPassword.txt"

	sm = magicBuilder.MakeSM()
	gpg = magicBuilder.MakePGP()
	gpg.LoadKeys()

	err := gpg.UnlockKey(testKeyFingerprint, testKeyPassword)

	if err != nil {
		SLog.UnsetTestMode()
		slog.Error(err)
		os.Exit(1)
	}

	router = GenRemoteSignerServerMux(slog, sm, gpg)

	code := m.Run()

	os.Exit(code)
}

// region GPG Endpoint Tests
func TestGenerateKey(t *testing.T) {
	genKeyBody := models.GPGGenerateKeyData{
		Identifier: "Test",
		Password:   "123456",
		Bits:       2048,
	}

	body, err := json.Marshal(genKeyBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/generateKey", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	errorDie(err, t)

	key := string(d)

	fingerPrint, err := remote_signer.GetFingerPrintFromKey(key)

	errorDie(err, t)

	err, _ = gpg.LoadKey(key)

	errorDie(err, t)

	err = gpg.UnlockKey(fingerPrint, genKeyBody.Password)

	errorDie(err, t)
}
func TestDecryptDataOnly(t *testing.T) {

	decryptBody := models.GPGDecryptData{
		DataOnly:         true,
		AsciiArmoredData: testDecryptDataOnly,
	}

	body, err := json.Marshal(decryptBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/decrypt", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	var decryptedData models.GPGDecryptedData

	err = json.Unmarshal(d, &decryptedData)

	errorDie(err, t)

	decryptedBytes, err := base64.StdEncoding.DecodeString(decryptedData.Base64Data)

	errorDie(err, t)

	if string(decryptedBytes) != testSignatureData {
		t.Errorf("Expected \"%s\" got \"%s\"", testSignatureData, string(decryptedBytes))
	}
}
func TestDecrypt(t *testing.T) {
	decryptBody := models.GPGDecryptData{
		DataOnly:         false,
		AsciiArmoredData: testDecryptDataAscii,
	}

	body, err := json.Marshal(decryptBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/decrypt", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	var decryptedData models.GPGDecryptedData

	err = json.Unmarshal(d, &decryptedData)

	errorDie(err, t)

	decryptedBytes, err := base64.StdEncoding.DecodeString(decryptedData.Base64Data)

	errorDie(err, t)

	if string(decryptedBytes) != testSignatureData {
		t.Errorf("Expected \"%s\" got \"%s\"", testSignatureData, string(decryptedBytes))
	}
}
func TestVerifySignature(t *testing.T) {
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
		Signature:  testSignatureSignature,
	}

	body, err := json.Marshal(verifyBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/verifySignature", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}
}
func TestVerifySignatureQuanto(t *testing.T) {
	quantoSignature := remote_signer.GPG2Quanto(testSignatureSignature, testKeyFingerprint, "SHA512")

	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
		Signature:  quantoSignature,
	}

	body, err := json.Marshal(verifyBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/verifySignatureQuanto", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		slog.Debug(errObj.StackTrace)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}
}
func TestSign(t *testing.T) {
	// region Generate Signature
	signBody := models.GPGSignData{
		FingerPrint: testKeyFingerprint,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
	}

	body, err := json.Marshal(signBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/sign", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	slog.Debug("Signature: %s", string(d))
	// endregion
	// region Verify Signature
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
		Signature:  string(d),
	}

	body, err = json.Marshal(verifyBody)

	errorDie(err, t)

	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/verifySignature", r)

	errorDie(err, t)

	res = executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}
	// endregion
}
func TestSignQuanto(t *testing.T) {
	// region Generate Signature
	signBody := models.GPGSignData{
		FingerPrint: testKeyFingerprint,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
	}

	body, err := json.Marshal(signBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/signQuanto", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	slog.Debug("Signature: %s", string(d))
	// endregion
	// region Verify Signature
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
		Signature:  string(d),
	}

	body, err = json.Marshal(verifyBody)

	errorDie(err, t)

	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/verifySignatureQuanto", r)

	errorDie(err, t)

	res = executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		t.Errorf("%s", errObj.String())
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}
	// endregion
}

// endregion