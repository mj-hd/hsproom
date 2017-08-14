all: hsproom ace

hsproom: go_packages
	gom build hsproom.go

go_packages:
	gom install

ace: submodules
	cd ace/; \
	npm install; \
	node Makefile.dryice.js -m -nc; \
	cd ../

submodules: .git
	git submodule update --init

clean:
	rm -f hsproom
