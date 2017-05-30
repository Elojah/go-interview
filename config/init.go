/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   init.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/25 17:42:15 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/30 22:03:29 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Port   string
		Pepper string
	}
	Psql struct {
		Host     string
		Name     string
		Password string
		Port     string
		User     string
	}
	Redis struct {
		Host     string
		Password string
		Port     string
	}
}

var (
	Conf Config
)

func Init(filepath string) (err error) {
	Conf, err = initConf(filepath)
	return err
}

func initConf(filepath string) (conf Config, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	parser := json.NewDecoder(f)
	err = parser.Decode(&conf)
	return
}
