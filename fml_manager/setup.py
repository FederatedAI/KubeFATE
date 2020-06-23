import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="fml_manager",
    version="0.3.0",
    author="Layne Peng/Jiahao Chen",
    author_email="jiahaoc@vmware.com",
    description="Python client for FATE cluster",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://not-decided.yet",
    packages=setuptools.find_packages(),
    install_requires=[
        'cachetools==3.0.0',
        'requests>=2.21.0'
    ],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
)
