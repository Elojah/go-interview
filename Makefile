# **************************************************************************** #
#                                                                              #
#                                                         :::      ::::::::    #
#    Makefile                                           :+:      :+:    :+:    #
#                                                     +:+ +:+         +:+      #
#    By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+         #
#                                                 +#+#+#+#+#+   +#+            #
#    Created: 2017/05/31 09:59:48 by hdezier           #+#    #+#              #
#    Updated: 2017/05/31 10:48:51 by hdezier          ###   ########.fr        #
#                                                                              #
# **************************************************************************** #

ALL:
	docker build -t eg_postgresql .
	docker run --rm -P --name pg_test -p 5432:5432 eg_postgresql&
	go test .
	go test ./auth
	go build
	./authserver&
