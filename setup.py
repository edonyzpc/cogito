from setuptools import setup, find_packages

setup(
    name="cogito",
    version="0.0.1",
    packages=find_packages(exclude=["tests"]),
    package_data={"": ["*.tcss"]},
    include_package_data=True,
    description="A LLM-based app serve as a personal AI assistant in Obsidian.",
    author="edonyzpc",
    author_email="edonyzpc@edony.ink",
    url="https://github.com/edonyzpc/cogito",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "License :: OSI Approved :: MIT License",
        "Topic :: Scientific/Engineering :: Artificial Intelligence",
    ],
    entry_points={
        "console_scripts": [
            "cogito = cogito.textual_ui.app:run",
        ],
        "cogito_function": [
            "google_web_search = cogito.runtime.function_calling.functions:google_web_search"
        ],
    },
    keywords="gpt,chat,llm,chatgpt,langchain,openai,qwen,textual,terminal",
    install_requires=[
        "httpx~=0.26.0",
        "humanize~=4.9.0",
        "langchain~=0.1.10",
        "langchain-core~=0.1.28",
        "polars~=0.20.7",
        "pydantic~=2.6.1",
        "pyperclip~=1.8.2",
        "rich~=13.7.0",
        "shortuuid~=1.0.11",
        "textual==0.50.1",
        "tiktoken~=0.5.2",
        "toolong==1.2.0",
        "pyyaml~=6.0.1",
        "setuptools~=69.1.0",
    ],
    extras_require={
        "openai": ["langchain-openai~=0.0.8"],
        "google": ["langchain-google-genai~=0.0.9"],
        "sap": ["generative-ai-hub-sdk~=1.2.2"],
        "anthropic": ["langchain-anthropic~=0.1.1"],
        "qwen": [
            "langchain-community~=0.0.38",
            "dashscope~=1.20.4",
        ],
        "all": [
            "langchain-openai~=0.0.8",
            "langchain-google-genai~=0.0.9",
            "generative-ai-hub-sdk~=1.2.2",
            "langchain-anthropic~=0.1.1",
            "langchain-community~=0.0.38",
            "dashscope~=1.20.4",
        ],
    },
)
