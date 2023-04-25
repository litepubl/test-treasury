package xmldata

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	Url     = "https://www.treasury.gov/ofac/downloads/sdn.xml"
	Timeout = 30 * time.Second
)

type Downloader struct {
	IgnoreEqualHash bool
	hash            [sha256.Size]byte
	xmlLen          int
}

var ErrEqualHash = errors.New("qqual xml hashes")

func (d *Downloader) Download(ctx context.Context) (*SdnList, error) {
	xmlData, err := d.request(ctx)
	if err != nil {
		return nil, err
	}

	if !d.IgnoreEqualHash && d.compareHash(xmlData) {
		return nil, ErrEqualHash
	}

	return d.decode(xmlData)
}

func (d *Downloader) request(ctx context.Context) ([]byte, error) {
	//for test purpose
	//return os.ReadFile("data.xml")

	client := http.Client{
		Timeout: Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*")
	req.Header.Add("User-Agent", "MSIE/15.0")
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (d *Downloader) decode(xmlData []byte) (*SdnList, error) {
	result := &SdnList{}
	err := xml.Unmarshal(xmlData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Downloader) compareHash(xmlData []byte) bool {
	len := len(xmlData)
	hash := sha256.Sum256(xmlData)

	result := (len == d.xmlLen) && (hash == d.hash)

	d.hash = hash
	d.xmlLen = len

	return result
}
