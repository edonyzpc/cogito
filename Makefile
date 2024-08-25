.PHONY: install-deps build-pypi cogito

PYTHONCMD := python3.10

install-deps:
	. venv/bin/activate && pip install .
	. venv/bin/activate && pip install '.[all]'

build-pypi: install-deps
	${PYTHONCMD} -m build

cogito:
	. venv/bin/activate && ${PYTHONCMD} -m main
