/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   hash.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/27 16:24:45 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/27 17:56:10 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package helpers

import (
	"authserver/config"
	"crypto/sha512"
)

func HashPassword(pwd []byte) []byte {
	hash := sha512.Sum512(append(pwd, []byte(config.Conf.Server.Pepper)...))
	return hash[:]
}
