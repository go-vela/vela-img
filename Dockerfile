# Copyright (c) 2021 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

#########################################################################
##    docker build --no-cache --target certs -t vela-img:certs .       ##
#########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-img:local .       ##
##########################################################

FROM r.j3ss.co/img:v0.5.11

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-img /bin/vela-img

ENTRYPOINT [ "/bin/vela-img" ]