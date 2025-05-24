package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CEDHostStatus struct {
	HostName  string
	Status    string
	UserCount int
}

func GetCEDHost(jumpHost SSHHost) (string, error) {

	req, _ := http.NewRequest(
		"GET",
		"http://jr3.cs.uec.ac.jp/23/materials/monitor.html",
		nil,
	)
	status, err := execHttpRequestWithFallback(jumpHost, *req)

	if err != nil {
		return "", fmt.Errorf("error checking terminal status: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(status))
	if err != nil {
		return "", fmt.Errorf("error parsing HTML document: %v", err)
	}

	var hostStatus []CEDHostStatus
	doc.Find("tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		hostName := s.Find("td:nth-child(1)").Text()
		status := s.Find("td:nth-child(2)").Text()
		userCount, err := strconv.Atoi(s.Find("td:nth-child(3)").Text())
		if err != nil {
			userCount = 1000
		}

		if status == "GOOD" {
			hostStatus = append(hostStatus, CEDHostStatus{
				HostName:  hostName,
				Status:    status,
				UserCount: userCount,
			})
		}
	})
	sort.SliceStable(hostStatus, func(i, j int) bool {
		return hostStatus[i].UserCount < hostStatus[j].UserCount
	})

	if len(hostStatus) == 0 {
		return "", fmt.Errorf("no usable computers found")
	}
	return fmt.Sprintf("%s.ced.cei.uec.ac.jp", hostStatus[0].HostName), nil
}
