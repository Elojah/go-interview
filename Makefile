# **************************************************************************** #
#                                                                              #
#                                                         :::      ::::::::    #
#    Makefile                                           :+:      :+:    :+:    #
#                                                     +:+ +:+         +:+      #
#    By: hdezier <hdezier@student.42.fr>            +#+  +:+       +#+         #
#                                                 +#+#+#+#+#+   +#+            #
#    Created: 2017/05/31 09:59:48 by hdezier           #+#    #+#              #
#    Updated: 2017/05/31 10:12:34 by hdezier          ###   ########.fr        #
#                                                                              #
# **************************************************************************** #

ALL:
	docker build -t eg_postgresql .
	docker run --rm -P --name pg_test eg_postgresql&
	go test .
	go test ./auth
	go build
	./authserver&
