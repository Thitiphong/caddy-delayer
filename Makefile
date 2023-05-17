CADDY_VERSION := v2.6.4
.EXPORT_ALL_VARIABLES:

build:
	xcaddy build $(CADDY_VERSION) --with github.com/Thitiphong/caddy-delayer=.

run:
	xcaddy run

list-modules:
	xcaddy list-modules