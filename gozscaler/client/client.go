package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

// Request ... // Needs to review this function
func (c *Client) Request(endpoint, method string, data []byte, contentType string) ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	if contentType == "" {
		contentType = contentTypeJSON
	}

	var req *http.Request
	var err error
	err = c.checkSession()
	if err != nil {
		return nil, err
	}
	req, err = http.NewRequest(method, c.URL+endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api responded with code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("api returned an error status:%d, %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// Create send HTTP Post request
func (c *Client) Create(endpoint string, o interface{}) (interface{}, error) {
	if o == nil {
		return nil, errors.New("tried to create with a nil payload not a Struct")
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("tried to create with a " + t.Kind().String() + " not a Struct")
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}

	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	if err != nil {
		return nil, err
	}
	id := reflect.Indirect(reflect.ValueOf(responseObject)).FieldByName("ID")

	log.Printf("Created Object with ID %v", id)
	return responseObject, nil
}

// Read ...
func (c *Client) Read(endpoint string, o interface{}) error {
	contentType := c.GetContentType()
	resp, err := c.Request(endpoint, "GET", nil, contentType)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, o)
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (c *Client) UpdateWithPut(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PUT", "application/json")
}

// Update ...
func (c *Client) Update(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PATCH", "application/merge-patch+json")
}

// Update ...
func (c *Client) updateGeneric(endpoint string, o interface{}, method, contentType string) (interface{}, error) {
	if o == nil {
		return nil, errors.New("tried to update with a nil payload not a Struct")
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("tried to update with a " + t.Kind().String() + " not a Struct")
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, method, data, contentType)
	if err != nil {
		return nil, err
	}

	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	return responseObject, err
}

// Delete ...
func (c *Client) Delete(endpoint string) error {
	_, err := c.Request(endpoint, "DELETE", nil, "application/json")
	if err != nil {
		return err
	}
	return nil
}
