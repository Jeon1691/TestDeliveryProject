package company

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jeon1691/testdeliveryproject"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type CJ struct{}

func (t CJ) Code() string {
	return "CJ"
}

func (t CJ) Name() string {
	return "CJ대한통운"
}

func (t CJ) TrackingUrl() string {
	return "http://nexs.cjgls.com/web/info.jsp?slipno=%s"
}

func (t CJ) Parse(trackingNumber string) (TestDeliveryGo.Track, *TestDeliveryGo.ApiError) {
	track := TestDeliveryGo.Track{}
	body, err := t.getHtml(trackingNumber)

	if err != nil {
		return track, TestDeliveryGo.NewApiError(TestDeliveryGo.RequestPageError, err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return track, TestDeliveryGo.NewApiError(TestDeliveryGo.ParseError, err.Error())
	}

	track = TestDeliveryGo.Track{
		TrackingNumber: trackingNumber,
		CompanyCode:    t.Code(),
		CompanyName:    t.Name(),
		Sender: strings.TrimSpace(doc.
			Find("table").Eq(2).
			Find("tbody tr").Eq(1).
			Find("td").Eq(0).Text()),
		Receiver: strings.TrimSpace(doc.
			Find("table").Eq(2).
			Find("tbody tr").Eq(1).
			Find("td").Eq(1).Text()),
		Signer: strings.TrimSpace(doc.
			Find("table").Eq(2).
			Find("tbody tr").Eq(2).
			Find("td").Eq(3).Text()),
	}

	history := []TestDeliveryGo.History{}

	numberReg, _ := regexp.Compile("[^0-9-]")

	doc.Find("table").Eq(4).Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		dateText := strings.TrimSpace(s.Find("td").Eq(0).Text()) + " " +
			strings.TrimSpace(s.Find("td").Eq(1).Text())

		if i > 0 && strings.Index(dateText, "Tel:") <= 0 {
			date, err := time.Parse("2006-01-02 15:04:05", dateText)
			if err != nil {
				log.Print(err)
			} else {
				statusText := strings.TrimSpace(s.Find("td").Eq(5).Text())

				if i == 1 {
					track.StatusCode = t.getStatus(statusText)
					track.StatusText = statusText
				}

				tel := numberReg.ReplaceAllString(s.Find("td table tr td").Eq(1).Text(), "")
				if tel == "--" {
					tel = ""
				}
				history = append([]TestDeliveryGo.History{
					TestDeliveryGo.History{
						Date:       date.Add(-time.Hour * 9).Unix(),
						DateText:   date.Format("2006-01-02 15:04"),
						Area:       strings.TrimSpace(s.Find("td table tr td").Eq(0).Text()),
						Tel:        tel,
						StatusCode: t.getStatus(statusText),
						StatusText: statusText,
					},
				}, history...)
			}
		}
	})

	track.History = history

	return track, nil
}

func (t CJ) getHtml(trackingNubmer string) (io.Reader, error) {
	url := fmt.Sprintf(t.TrackingUrl(), trackingNubmer)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	convertBody := transform.NewReader(bytes.NewReader(body), korean.EUCKR.NewDecoder())
	return convertBody, nil
}

func (t CJ) getStatus(statusText string) TestDeliveryGo.TrackingStatus {
	switch statusText {
	case "SM입고":
		return TestDeliveryGo.Ready
	case "집화처리":
		return TestDeliveryGo.PickupCompelete
	case "간선상차":
		return TestDeliveryGo.Loading
	case "간선하차":
		return TestDeliveryGo.Unloading
	case "배달출발":
		return TestDeliveryGo.DeleveryStart
	case "배달완료":
		return TestDeliveryGo.DeleveryComplete
	case "미배달":
		return TestDeliveryGo.DoNOtDelevery
	}
	return TestDeliveryGo.UnknownStatus
}
