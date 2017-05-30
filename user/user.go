/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   user.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/27 16:23:49 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/28 03:02:56 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package user

type UserPublic struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

type User struct {
	UserPublic
	Login    string
	Password string
}

var (
	DefaultUser = UserPublic{
		ID:        0,
		FirstName: ``,
		LastName:  ``,
		Email:     ``,
	}
)
