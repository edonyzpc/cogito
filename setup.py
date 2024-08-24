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
        "httpx",
        "humanize",
        "langchain",
        "langchain-core",
        "polars",
        "pydantic",
        "pyperclip",
        "rich",
        "shortuuid",
        "textual",
        "tiktoken",
        "toolong",
        "pyyaml",
        "setuptools",
    ],
    extras_require={
        "openai": ["langchain-openai"],
        "google": ["langchain-google-genai"],
        "sap": ["generative-ai-hub-sdk"],
        "anthropic": ["langchain-anthropic"],
        "qwen": [
            "langchain-community",
            "dashscope",
        ],
        "all": [
            "langchain-openai",
            "langchain-google-genai",
            "generative-ai-hub-sdk",
            "langchain-anthropic",
            "langchain-community",
            "dashscope",
        ],
    },
)
