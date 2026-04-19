#!/bin/bash

# run redis daemon
redis-server /conf/redis.conf --save 20 1

# go module proxy url as argument
/pkgsite -direct_proxy -bypass_license_check -proxy_url $1 -host 0.0.0.0:8888 -static /static -third_party /third_party
