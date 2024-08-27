.PHONY: install-deps build-pypi cogito

PYTHONCMD := python3.10

install-deps:
	. venv/bin/activate && pip install .
	. venv/bin/activate && pip install '.[all]'

build-pypi: install-deps
	${PYTHONCMD} -m build

cogito:
	. venv/bin/activate && ${PYTHONCMD} -m main

clean:
	rm -rf ./bin/conversations/*
	rm -rf ./bin/exports/*
	rm -rf ./bin/logging/*

dev-console:
# ignore the system and event logs
	textual console -x SYSTEM -x EVENT --port 7342

debug-mentis:
	textual run --dev  cogito.textual_ui.app:run --port 7342