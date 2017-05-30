/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/25 15:41:38 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/31 00:25:50 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"authserver/auth"
	"authserver/config"
	"authserver/context"
	"authserver/helpers"
	"authserver/psql"
	"authserver/redis"
	"authserver/user"
	"encoding/json"
	"fmt"
	"net/http"
)

// Define all routes we need to serve
func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc(`/api/login`, auth.Validate(postUser, auth.BasicUser))
	mux.HandleFunc(`/user/add`, addUser)
}

func addUserPSQL(u user.User) {
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
	_, _ = psql.Client.Exec(query)
}

// For tests purpose
func addUser(w http.ResponseWriter, req *http.Request) {
	u := user.User{}
	decoder := json.NewDecoder(req.Body)
	_ = decoder.Decode(&u)
	defer req.Body.Close()
	addUserPSQL(u)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created\n"))
}

// Serve `/api/login`
func postUser(ctx context.Context) {
	fmt.Println(`[INFO] Accepted request for: /user`)
	ctx.W.WriteHeader(http.StatusAccepted)
	b, err := json.Marshal(&ctx.User)
	if err != nil {
		fmt.Println(`[ERROR] Json marshal reponse failed`)
	}
	ctx.W.Write(b)
}

func main() {

	// init configuration structure from local conf.json
	configFile := `./conf.json`
	err := config.Init(configFile)
	if err != nil {
		fmt.Printf(
			"[ERROR] occured during configuration initialization with file %s:\n",
			configFile,
		)
		fmt.Println(err.Error())
		return
	}
	conf := config.Conf

	// init database client
	err = psql.Init(&conf)
	if err != nil {
		fmt.Printf(
			"[ERROR] occured during database initialization on %s:%s/%s:\n",
			conf.Psql.Host,
			conf.Psql.Port,
			conf.Psql.Name,
		)
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("[INFO] Connection to DB %s: OK\n", conf.Psql.Name)

	// Build databse schemas
	err = psql.BuildSchemas()
	if err != nil {
		fmt.Printf("[ERROR] occured during database schema building\n")
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("[INFO] Schema building %s: OK\n", conf.Psql.Name)

	// init redis client
	err = redis.Init(&conf)
	if err != nil {
		fmt.Printf(
			"[ERROR] occured during redis initialization with config %s:%s:\n",
			conf.Redis.Host,
			conf.Redis.Port,
		)
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("[INFO] Connection to Redis %s: OK\n", conf.Redis.Host)

	mux := http.NewServeMux()
	AddRoutes(mux)
	s := &http.Server{
		Addr:    fmt.Sprintf(`:%s`, conf.Server.Port),
		Handler: mux,
	}

	fmt.Printf("[INFO] Server will listen on port %s\n", conf.Server.Port)
	defer redis.Client.Close()
	defer psql.Client.Close()
	defer s.Close()

	err = s.ListenAndServe()
	fmt.Printf("[INFO] Server closed: %s\n", err.Error())
}
