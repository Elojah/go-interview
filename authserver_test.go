/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   authserver_test.go                                 :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/30 18:57:47 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/31 10:11:29 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"authserver/auth"
	"authserver/config"
	"authserver/helpers"
	"authserver/psql"
	"authserver/redis"
	"authserver/user"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	called              = false
	test_login          = `test_login`
	test_password       = `test_password`
	test_wrong_password = `test_wrong_password`
	fake_user           = user.User{
		Login:    test_login,
		Password: test_password,
		UserPublic: user.UserPublic{
			ID:        0,
			FirstName: `test_fn`,
			LastName:  `test_ln`,
			Email:     `test@test.test`,
		},
	}
	test_table_user = `user_infos`
)

// Helpers for test purposes
func addUserRedis(t *testing.T, u user.User) {
	mapRedisUser := map[string]interface{}{
		`ID`:        string(u.ID),
		`Password`:  string(helpers.HashPassword([]byte(u.Password))),
		`FirstName`: u.FirstName,
		`LastName`:  u.LastName,
		`Email`:     u.Email,
	}
	rdStatus := redis.Client.HMSet(u.Login, mapRedisUser)
	if rdStatus.Err() != nil {
		t.Errorf(`[ERROR] Failed to add test user in redis`)
	}
}
func addUserFakePSQL(t *testing.T, u user.User) {
	query := fmt.Sprintf(`	INSERT INTO %s (
									"login",
									"password",
									"firstname",
									"lastname",
									"email"
									)
							VALUES (
									'%s',
									E'%x'::bytea,
									'%s',
									'%s',
									'%s'
							);`,
		psql.TableNames[`user`],
		u.Login,
		helpers.HashPassword([]byte(u.Password)),
		u.FirstName,
		u.LastName,
		u.Email,
	)

	_, err := psql.Client.Exec(query)
	if err != nil {
		t.Errorf("[ERROR] occured inserting value to PSQL. Check client is valid\n%s",
			err.Error())
	}
}

func TestMain(m *testing.M) {
	configFile := `./conf.test.json`
	err := config.Init(configFile)
	if err != nil {
		fmt.Printf("Failed to init configuration from %s:\n%s", configFile, err.Error())
		return
	}
	err = redis.Init(&config.Conf)
	if err != nil {
		fmt.Printf("Failed to init redis from %s:\n%s", configFile, err.Error())
		return
	}
	err = psql.Init(&config.Conf)
	if err != nil {
		fmt.Printf("Failed to init psql from %s:\n%s", configFile, err.Error())
		return
	}

	query := fmt.Sprintf(`DROP TABLE %s;`, test_table_user)
	_, err = psql.Client.Exec(query)
	if err != nil {
		fmt.Printf("[ERROR] occured dropping table in PSQL. Check client is valid\n%s",
			err.Error())
	}
	psql.TableNames[`user`] = test_table_user
	err = psql.BuildSchemas()
	if err != nil {
		fmt.Printf("Failed to build psql schemas:\n%s", err.Error())
		return
	}
	fmt.Println(`[TESTS]Start running`)
	flag.Parse()
	resStatus := m.Run()
	fmt.Println(`[TESTS]End`)
	psql.Client.Close()
	os.Exit(resStatus)
}

func TestValidLog(t *testing.T) {
	jsonReq := []byte(fmt.Sprintf(
		`{"login":"%s", "password":"%s"}`,
		fake_user.Login,
		fake_user.Password,
	))
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonReq))

	addUserFakePSQL(t, fake_user)

	mux := http.NewServeMux()
	mux.HandleFunc(`/api/login`, auth.Validate(postUser, auth.BasicUser))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	resp := user.UserPublic{}
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Email != fake_user.Email ||
		resp.FirstName != fake_user.FirstName ||
		resp.LastName != fake_user.LastName {
		t.Errorf("/api/login response is not expected\n")
	}
}

func TestInvalidLog(t *testing.T) {
	jsonReq := []byte(fmt.Sprintf(
		`{"login":"%s", "password":"%s"}`,
		fake_user.Login,
		test_wrong_password,
	))
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonReq))

	mux := http.NewServeMux()
	mux.HandleFunc(`/api/login`, auth.Validate(postUser, auth.BasicUser))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
