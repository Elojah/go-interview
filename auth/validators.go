/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   validators.go                                      :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/27 15:14:31 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/31 00:22:33 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package auth

import (
	"authserver/context"
	"authserver/helpers"
	"authserver/psql"
	"authserver/redis"
	"authserver/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserValidator func(req *http.Request) (user.UserPublic, bool)

func Validate(handler context.HandleCtx, validator UserValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if publicUser, isValid := validator(req); isValid {
			handler(context.Context{
				W:    w,
				Req:  req,
				User: publicUser,
			})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Please log in to access to this page\n"))
		}
	}
}

func All(req *http.Request) (user.UserPublic, bool) {
	return user.DefaultUser, true
}

func None(req *http.Request) (user.UserPublic, bool) {
	return user.DefaultUser, false
}

func basicUserRedis(credentials *user.User) (user.UserPublic, bool) {

	val, err := redis.Client.HMGet(
		credentials.Login,
		`ID`,
		`FirstName`,
		`LastName`,
		`Email`,
		`Password`,
	).Result()
	if err != nil {
		return user.DefaultUser, false
	}
	if val[4] == nil ||
		(string(helpers.HashPassword([]byte(credentials.Password))) != val[4].(string)) {
		fmt.Printf("[INFO] User(%s)/Pwd not found on redis\n", credentials.Login)
		return user.DefaultUser, false
	}
	id, _ := strconv.Atoi(val[0].(string))
	return user.UserPublic{
		ID:        id,
		FirstName: val[1].(string),
		LastName:  val[2].(string),
		Email:     val[3].(string),
	}, true
}

func basicUserPSQL(credentials *user.User) (publicUser user.UserPublic, ok bool) {

	err := psql.Client.QueryRow(fmt.Sprintf(`	SELECT id, firstname, lastname, email
									FROM %s
									WHERE login = '%s'
									AND password = '%x'
									;`,
		psql.TableNames[`user`],
		credentials.Login,
		helpers.HashPassword([]byte(credentials.Password)))).Scan(
		&publicUser.ID,
		&publicUser.FirstName,
		&publicUser.LastName,
		&publicUser.Email,
	)
	if err != nil {
		fmt.Printf("[INFO] User(%s)/Pwd not found on DB\n", credentials.Login)
		return user.DefaultUser, false
	}
	return publicUser, true
}

func BasicUser(req *http.Request) (publicUser user.UserPublic, isValid bool) {

	if req == nil {
		return user.DefaultUser, false
	}
	credentials := user.User{}
	parser := json.NewDecoder(req.Body)
	err := parser.Decode(&credentials)
	defer req.Body.Close()
	if err != nil {
		return user.DefaultUser, false
	}

	publicUser, isValid = basicUserRedis(&credentials)
	if isValid {
		return
	}
	publicUser, isValid = basicUserPSQL(&credentials)
	if isValid {
		return
	}
	return
}
