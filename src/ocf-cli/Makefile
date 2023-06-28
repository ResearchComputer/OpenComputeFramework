install:
	python setup.py install

format:
	autoflake -i **/*.py
	isort -i ocf_cli/**/*.py
	yapf -i **/*.py

clean:
	rm -rf build
	rm -rf dist
	rm -rf ocf_cli.egg-info

build:
	python -m build --wheel

test:
	PYTHONPATH=./ python3 tests/server.py

publish-test:
	twine upload --repository-url https://test.pypi.org/legacy/ dist/*

publish:
	twine upload dist/*

install-local:
	pip install -e .

install-test:
	pip install --index-url https://test.pypi.org/simple/ --extra-index-url https://pypi.org/simple/ ocf_cli -U

