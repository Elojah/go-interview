/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   init.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/25 17:30:46 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/25 18:00:21 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package redis

import (
	"authserver/config"
	"fmt"
	"github.com/go-redis/redis"
)

var (
	Client *redis.Client
)

func Init(conf *config.Config) (err error) {
	Client, err = initRedis(conf)
	return err
}

func initRedis(conf *config.Config) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%s`, conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Password, // no Password set
		DB:       0,                   // use default DB
	})
	_, err = client.Ping().Result()
	return
}
