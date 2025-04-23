package zammad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// New returns a new Zammad client initialized with an http client. Authentication need to be set seperately. The http
// client uses a timeout of 5 seconds.
func New(URL string) *Client {
	return &Client{Client: &http.Client{Timeout: 5 * time.Second}, Url: URL}
}

// NewRequest constructs a request and converts the payload to JSON.
func (c *Client) NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	if c.FromFunc != nil {
		x := c.FromFunc()
		if x != "" {
			req.Header.Set("From", x)
		}
	}

	return req, nil
}

// send makes a request to the API, the response body will be unmarshaled into v, or if v is an io.Writer, the response
// will be written to it without decoding. This can be helpful when debugging.
//func (c *Client) send(req *http.Request, v interface{}) error {
//	resp, err := c.Client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode < 200 || resp.StatusCode > 299 {
//		errResp := &ErrorResponse{}
//		data, err := io.ReadAll(resp.Body)
//
//		if err == nil && len(data) > 0 {
//			err = json.Unmarshal(data, errResp)
//			if err != nil {
//				return err
//			}
//		}
//
//		return errResp
//	}
//
//	if v == nil {
//		return nil
//	}
//
//	if w, ok := v.(io.Writer); ok {
//		_, err = io.Copy(w, resp.Body)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//
//	return json.NewDecoder(resp.Body).Decode(v)
//}

func (c *Client) send(req *http.Request, v interface{}) error {

	fmt.Println("Final request URL:", req.URL.String())
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Response-Body lesen
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 🪵 Pretty-Print JSON zur Debug-Ausgabe
	fmt.Println("==== JSON Response ====")
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		fmt.Println("❌ Kein gültiges JSON – Rohdaten:")
		fmt.Println(string(data))
	} else {
		fmt.Println(prettyJSON.String())
	}
	fmt.Printf("Statuscode:  %d", resp.StatusCode)
	fmt.Println("========================")

	// Body für spätere Verarbeitung wiederherstellen
	resp.Body = io.NopCloser(bytes.NewBuffer(data))

	// Fehlerstatus behandeln
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{}
		if len(data) > 0 {
			if err := json.Unmarshal(data, errResp); err != nil {
				return err
			}
		}
		return errResp
	}

	// Optionales Schreiben in io.Writer
	if v == nil {
		return nil
	}
	if w, ok := v.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		return err
	}

	// Normale JSON-Dekodierung
	return json.NewDecoder(resp.Body).Decode(v)
}

// sendWithAuth makes a request to the API and apply the proper authentication header automatically.
func (c *Client) sendWithAuth(req *http.Request, v interface{}) error {
	//Detect Authentication Type
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))
	}
	if c.OAuth != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OAuth))
	}

	return c.send(req, v)
}
