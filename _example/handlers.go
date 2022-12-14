/*
Code generated by Near Gothic. DO NOT EDIT.
Generated on 2022/10/19.
Manual edits may be lost when next regenerated.

Near Gothic is (C) K Cartlidge, 2022.
No rights are asserted over generated code.
*/

package main

import (
	"fmt"
	"net/http"
)

func (s *Server) RegisterDefaultHandlers() {
	s.addListAccount()
	s.addListAccountSetting()
	s.addListSetting()
	fmt.Println()
}

// GET handler for Account list
func (s *Server) addListAccount() {
	fmt.Println(fmt.Sprintf("GET    %s/accounts", s.urlPrefix))
	s.router.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		data, err := s.AccountRepo.List()
		s.SendJSONor500(w, 200, data, err)
	}).Methods("GET")
}

// GET handler for Account Setting list
func (s *Server) addListAccountSetting() {
	fmt.Println(fmt.Sprintf("GET    %s/account-settings", s.urlPrefix))
	s.router.HandleFunc("/account-settings", func(w http.ResponseWriter, r *http.Request) {
		data, err := s.AccountSettingRepo.List()
		s.SendJSONor500(w, 200, data, err)
	}).Methods("GET")
}

// GET handler for Setting list
func (s *Server) addListSetting() {
	fmt.Println(fmt.Sprintf("GET    %s/settings", s.urlPrefix))
	s.router.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		data, err := s.SettingRepo.List()
		s.SendJSONor500(w, 200, data, err)
	}).Methods("GET")
}
