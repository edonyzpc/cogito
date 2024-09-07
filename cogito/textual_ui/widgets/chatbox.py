from __future__ import annotations

import re
from langchain.schema import BaseMessage, AIMessage

import pyperclip

from rich.console import Console, ConsoleOptions, RenderableType, RenderResult
from rich.markdown import Markdown, TextElement, MarkdownContext
from rich.panel import Panel
from rich.box import HEAVY
from rich.text import Text
from rich.style import Style
from textual.binding import Binding
from textual.geometry import Size
from textual.widget import Widget
from textual.containers import Container
from markdown_it.token import Token

from cogito.textual_ui.screens.message_info_modal import MessageInfo
from cogito.utils import format_timestamp
from cogito.runtime.langchain.schema import new_message_of_type


class ChatboxContainer(Container): ...


class Chatbox(Widget, can_focus=True):
    BINDINGS = [
        Binding(
            key="ctrl+s",
            action="focus('cl-option-list')",
            description="Focus List",
            key_display="^s",
        ),
        Binding(
            key="i",
            action="focus('chat-input')",
            description="Focus Input",
            key_display="i",
        ),
        Binding(
            key="d", action="details", description="Message details", key_display="d"
        ),
        Binding(key="c", action="copy", description="Copy Message", key_display="c"),
        Binding(key="`", action="copy_code", description="Copy Code"),
    ]

    def __init__(
        self,
        *,
        model_name: str,
        message: BaseMessage | None = None,
        name: str | None = None,
        id: str | None = None,
        classes: str | None = None,
        disabled: bool = False,
    ) -> None:
        super().__init__(
            name=name,
            id=id,
            classes=classes,
            disabled=disabled,
        )
        self._message = message or new_message_of_type(AIMessage)

        self.model_name = model_name
        timestamp = format_timestamp(
            self.message.additional_kwargs.get("timestamp", 0) or 0
        )
        self.tooltip = f"Sent {timestamp}"

    @property
    def is_ai_message(self):
        return self.message and self.message.type.lower().startswith("ai")

    def on_mount(self) -> None:
        if self.message.type.lower().startswith("ai"):
            self.add_class("assistant-message")
        else:
            self.add_class("user-message")

    def get_code_blocks(self, markdown_string):
        pattern = r"```(.*?)\n(.*?)```"
        code_blocks = re.findall(pattern, markdown_string, re.DOTALL)
        return code_blocks

    def action_copy_code(self):
        codeblocks = self.get_code_blocks(self.message.content)
        output = ""
        if codeblocks:
            for lang, code in codeblocks:
                output += f"{code}\n\n"
            pyperclip.copy(output)
            self.notify("Codeblocks have been copied to clipboard", timeout=3)
        else:
            self.notify("There are no codeblocks in the message to copy", timeout=3)

    def action_copy(self) -> None:
        pyperclip.copy(self.message.content)
        self.notify("Message content has been copied to clipboard", timeout=3)

    def action_details(self) -> None:
        self.app.push_screen(
            MessageInfo(message=self.message, model_name=self.model_name)
        )

    @property
    def markdown(self) -> Markdown:
        return Markdown(self.message.content or "", justify="left")

    def render(self) -> RenderableType:
        self.markdown.elements["heading_open"] = MyHeading
        return self.markdown

    def customized_rich_console(
        self, console: Console, options: ConsoleOptions
    ) -> RenderResult:
        text = self.text
        text.justify = "left"
        if self.tag == "h1":
            # Draw a border around h1s
            yield Panel(
                text,
                box=HEAVY,
                style="markdown.h1.border",
            )
        else:
            # Styled text for h2 and beyond
            if self.tag == "h2":
                yield Text("")
            yield text

    def get_content_width(self, container: Size, viewport: Size) -> int:
        # Naive approach. Can sometimes look strange, but works well enough.
        content = self.message.content or ""
        return min(len(content), container.width)

    @property
    def message(self):
        return self._message

    @message.setter
    def message(self, m):
        if not isinstance(m, BaseMessage):
            raise ValueError("Message must be a BaseMessage instance")
        self._message = m


class MyHeading(TextElement):
    """A heading."""

    @classmethod
    def create(cls, markdown: "Markdown", token: Token) -> "MyHeading":
        return cls(token.tag)

    def on_enter(self, context: "MarkdownContext") -> None:
        self.text = Text()
        context.enter_style(self.style_name)

    def __init__(self, tag: str) -> None:
        self.tag = tag
        self.style_name = f"markdown.{tag}"
        super().__init__()

    def __rich_console__(
        self, console: Console, options: ConsoleOptions
    ) -> RenderResult:
        text = self.text
        text.justify = "left"
        if self.tag == "h1":
            # Draw a border around h1s
            yield Panel(
                text,
                box=HEAVY,
                style="markdown.h1.border",
            )
        else:
            # Styled text for h2 and beyond
            if self.tag == "h2":
                yield Text("")

            if self.tag == "h3":
                color = "green"
            elif self.tag == "h4":
                color = "blue"
            else:
                color = "orange"
            customize_style = Style(bold=True, color=color)
            text.style = customize_style
            yield text
