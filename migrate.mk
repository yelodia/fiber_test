# migrate options
GOPATH?=$(shell go env GOPATH)
GOBIN:=$(GOPATH)/bin
GOOSE_BASE_CMD=${GOBIN}/goose

migration-up:		## Migrate the DB to the most recent version available *** make -f migrate.mk -C . migration-up
	${GOOSE_BASE_CMD} up

migration-down: 	## Roll back the version by 1 *** make -f migrate.mk -C . migration-down
	${GOOSE_BASE_CMD} down

migration-reset:	## Roll back all migrations *** make -f migrate.mk -C . migration-reset
	${GOOSE_BASE_CMD} reset

.PHONY: migration-up migration-down migration-reset

prerequisites: migration-up migration-down migration-reset

target: prerequisites
