/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   types.go                                           :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/05/27 16:12:39 by hdezier           #+#    #+#             */
/*   Updated: 2017/05/27 16:41:39 by hdezier          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package context

import (
	"authserver/user"
	"net/http"
)

type Context struct {
	W    http.ResponseWriter
	Req  *http.Request
	User user.UserPublic
}
type HandleCtx func(ctx Context)
