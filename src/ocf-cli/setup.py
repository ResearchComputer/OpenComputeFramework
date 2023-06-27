from setuptools import setup, find_packages

setup(
    name='ocf-cli',
    author="Xiaozhe Yao",
    author_email="enquiry@autoai.org",
    description="OCF Client",
    version='0.0.1',
    scripts=["ocf_cli/bin/ocf-cli"],
    package_dir={'ocf_cli': 'ocf_cli'},
    packages=find_packages(),
    install_requires=[
        "typer",
        "requests",
        "rich",
        "loguru",
        "huggingface-hub",
        "pynvml"
    ]
)