/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   init.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/25 17:31:28 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/30 23:57:27 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package psql

import (
	"authserver/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	Client     *sql.DB
	TableNames = map[string]string{
		`user`: `user_infos`,
	}
)

// Init psql.Client with config.Psql data
// Must be call before ANY attempt to use psql.Client
func Init(conf *config.Config) (err error) {
	Client, err = initPSQL(conf)
	if err != nil {
		return
	}
	return err
}

func initPSQL(conf *config.Config) (*sql.DB, error) {
	// TODO Enable SSL
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		conf.Psql.User,
		conf.Psql.Password,
		conf.Psql.Host,
		conf.Psql.Port,
		conf.Psql.Name,
	)
	return (sql.Open("postgres", url))
}

// Build SQL schemas from scratch
// FTM Schemas are hardcode without any versioning inside the function
func BuildSchemas() (err error) {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
					"id"        serial PRIMARY KEY,
					"login"     text NOT NULL UNIQUE,
					"password"  bytea NOT NULL,
					"firstname" text,
					"lastname"  text,
					"email"     text NOT NULL UNIQUE
				);`, TableNames[`user`])
	_, err = Client.Exec(query)
	return err
}
