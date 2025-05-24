package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var pcIds = []string{
	"a11", "a12", "a13", "a14", "a15", "a16", "a17", "a18",
	"a21", "a22", "a23", "a24", "a25", "a26", "a27", "a28",
	"a31", "a32", "a33", "a34", "a35", "a36", "a37", "a38",
	"a41", "a42", "a43", "a44", "a45", "a46", "a47", "a48",
	"a51", "a52", "a53", "a54", "a55", "a56", "a57", "a58",
	"a61", "a62", "a63", "a64", "a65", "a66", "a67", "a68",
	"a71", "a72", "a73", "a74", "a75", "a76", "a77", "a78",
	"b11", "b12", "b13", "b14", "b15", "b16", "b17", "b18",
	"b21", "b22", "b23", "b24", "b25", "b26", "b27", "b28",
	"b31", "b32", "b33", "b34", "b35", "b36", "b37", "b38",
	"b41", "b42", "b43", "b44", "b45", "b46", "b47", "b48",
	"b51", "b52", "b53", "b54", "b55", "b56", "b57", "b58",
	"b61", "b62", "b63", "b64", "b65", "b66", "b67", "b68",
	"b71", "b72", "b73", "b74", "b75", "b76", "b77", "b78",
}

func GetIEDHost(jumpHost SSHHost) (string, error) {
	poweroffPCId := ""
	for _, pcid := range ShuffleSlice(pcIds) {
		req, _ := http.NewRequest(
			"POST",
			"http://termsrv.ied.inf.uec.ac.jp/ajax1.php",
			strings.NewReader(fmt.Sprintf("request=%s", pcid)),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		status, err := execHttpRequestWithFallback(jumpHost, *req)

		if err != nil {
			verboseLogger.Printf("Error checking PC ID %s: %v", pcid, err)
			continue
		}
		if strings.Contains(status, "usable") {
			return fmt.Sprintf("%s.ied.inf.uec.ac.jp", pcid), nil
		}
		if strings.Contains(status, "poweroff") {
			poweroffPCId = pcid
		}
	}

	if poweroffPCId != "" {
		req, _ := http.NewRequest(
			"POST",
			"http://termsrv.ied.inf.uec.ac.jp/ajax2.php",
			strings.NewReader(fmt.Sprintf("request=%s", poweroffPCId)),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		log.Println("Powering on PC ID:", poweroffPCId)

		status, err := execHttpRequestWithFallback(jumpHost, *req)
		if err != nil {
			return "", fmt.Errorf("error powering on PC ID %s: %v", poweroffPCId, err)
		}

		req.URL, _ = url.Parse("http://termsrv.ied.inf.uec.ac.jp/ajax1.php")

		for range 60 {
			if strings.Contains(status, "usable") {
				return fmt.Sprintf("%s.ied.inf.uec.ac.jp", poweroffPCId), nil
			}
			time.Sleep(1 * time.Second)
			status, _ = execHttpRequestWithFallback(jumpHost, *req)
		}
	}

	return "", fmt.Errorf("no usable PC found")
}
