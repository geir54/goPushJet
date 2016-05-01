// copyright (c) 2016 geir54

package goPushJet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type servResp struct {
	Service Service `json:"service"`
}

type Service struct {
	Created int    `json:"created"`
	Icon    string `json:"icon"`
	Name    string `json:"name"`
	Public  string `json:"public"`
	Secret  string `json:"secret"`
}

type msgResp struct {
	Status string   `json:"status"`
	Error  errorMsg `json:"error"`
}

type errorMsg struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// GetQR - Get QR image
func (serv *Service) GetQR() string {
	return "https://chart.googleapis.com/chart?cht=qr&chl=" + serv.Public + "&choe=UTF-8&chs=200x200"
}

// CreateService - Create new service
func CreateService(name, icon string) (Service, error) {
	resp, err := http.PostForm("https://api.pushjet.io/service",
		url.Values{"name": {name}, "icon": {icon}})

	if err != nil {
		return Service{}, err
	}
	defer resp.Body.Close()

	ser := servResp{}
	err = json.NewDecoder(resp.Body).Decode(&ser)
	if err != nil {
		return Service{}, err
	}

	return ser.Service, nil
}

// SendMessage -
// secret: required stringd2d1820d56b862a6f5b1a69a7af730fa The service secret token
// message: required string Your server is on fire! The notification text
// title: string A custom message title
// level: integer 3 The importance level from 1(low) to 5(high)
// link: string http://i.imgur.com/TerUkQY.gif An optional link
func SendMessage(secret, message, title string, level int, link string) error {
	resp, err := http.PostForm("https://api.pushjet.io/message",
		url.Values{"secret": {secret}, "message": {message}, "title": {title}, "level": {strconv.Itoa(level)}, "link": {link}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	msg := msgResp{}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	if msg.Status != "ok" {
		return errors.New("Did not return status OK")
	}

	return nil
}
