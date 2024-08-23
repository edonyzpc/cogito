.PHONY: install-deps build-pypi cogito

PYTHONCMD := python3.10

install-deps:
	pip install .
	pip install '.[all]'

build-pypi: install-deps
	${PYTHONCMD} -m build

cogito:
	. venv/bin/activate && ${PYTHONCMD} -m main