/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   auth_test.go                                       :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/25 20:17:37 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/31 00:11:02 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package auth

import (
	"authserver/config"
	"authserver/context"
	"authserver/helpers"
	"authserver/psql"
	"authserver/redis"
	"authserver/user"
	"flag"
	"fmt"
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
	test_table_user = `test_user_infos`
)

func TestMain(m *testing.M) {
	configFile := `../conf.test.json`
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

func fakeHandler(ctx context.Context) {
	called = true
}

func TestValidate(t *testing.T) {
	w := httptest.NewRecorder()

	called = false
	handler := Validate(fakeHandler, All)
	handler(w, nil)
	if called == false {
	}

	called = false
	handler = Validate(fakeHandler, None)
	handler(w, nil)
	if called == true {
	}
	called = false
}

func TestBasicUserRedis(t *testing.T) {
	fakeRedisUser := map[string]interface{}{
		`ID`:        string(fake_user.ID),
		`Password`:  string(helpers.HashPassword([]byte(fake_user.Password))),
		`FirstName`: fake_user.FirstName,
		`LastName`:  fake_user.LastName,
		`Email`:     fake_user.Email,
	}
	rdStatus := redis.Client.HMSet(fake_user.Login, fakeRedisUser)
	if rdStatus.Err() != nil {
		t.Errorf("[ERROR] %s\n", rdStatus.Err().Error())
		t.Errorf("[ERROR] occured setting value to Redis. Check client is valid\n")
	}

	publicUser, ok := basicUserRedis(&fake_user)
	if !ok || publicUser.Email != fake_user.Email {
		t.Errorf("user_mail:%s \t ok:%s \n", publicUser.Email, ok)
	}

	wrongUser := fake_user
	wrongUser.Password = test_wrong_password
	publicUser, ok = basicUserRedis(&wrongUser)
	if ok {
		t.Errorf("user_mail:%s \t ok:%s \n", publicUser.Email, ok)
	}

	redis.Client.Del(test_login)
}

func addUserPSQL(t *testing.T, u *user.User) {
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
		test_table_user,
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

func TestBasicUserPSQL(t *testing.T) {

	addUserPSQL(t, &fake_user)

	publicUser, ok := basicUserPSQL(&fake_user)
	if !ok || publicUser.Email != fake_user.Email {
		t.Errorf("user_mail:%s \t ok:%s \n", publicUser.Email, ok)
	}

	wrongUser := fake_user
	wrongUser.Password = test_wrong_password
	publicUser, ok = basicUserPSQL(&wrongUser)
	if ok {
		t.Errorf("user_mail:%s \t ok:%s \n", publicUser.Email, ok)
	}
}
