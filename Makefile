unit:
	@sh test/unit-test.sh

unit-verbose:
	@sh test/unit-test-verbose.sh

init-git-hooks:
	@git config core.hooksPath .githooks