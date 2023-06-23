format:
	autoflake -i **/*.py
	yapf -i **/*.py

clean:
	rm -rf build
	rm -rf dist
	rm -rf *.egg-info

build:
	python -m build --wheel

publish:
	twine upload dist/*