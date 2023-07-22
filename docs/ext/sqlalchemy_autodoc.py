from typing import Any
from sphinx.ext import autodoc
from sqlalchemy.orm import class_mapper
from sqlalchemy.ext.declarative import declarative_base
from sphinx.application import Sphinx
from docutils.statemachine import StringList


class SqlAlchemyDocumenter(autodoc.ClassDocumenter):
    objtype = "sqlalchemy"
    directivetype = "class"
    priority = 10 + autodoc.ClassDocumenter.priority

    @classmethod
    def can_document_member(
        cls, member: Any, membername: str, isattr: bool, parent: Any
    ) -> bool:
        return isinstance(member, declarative_base().__class__)

    def add_directive_header(self, sig: str) -> None:
        super().add_directive_header(sig)
        self.add_line("   :final:", self.get_sourcename())

    def add_content(
        self,
        more_content: StringList | None,
        no_docstring: bool = False,
    ) -> None:
        sourcename = self.get_sourcename()

        self.add_line(".. list-table::", sourcename)
        self.add_line("   :header-rows: 1", sourcename)
        self.add_line("", sourcename)
        self.add_line("   * - Column", sourcename)
        self.add_line("     - Type", sourcename)

        mapper = class_mapper(self.object)
        for column in mapper.columns:
            self.add_line("", sourcename)
            self.add_line("   * - {}".format(column.name), sourcename)
            self.add_line("     - {}".format(column.type), sourcename)


def setup(app: Sphinx):
    app.add_autodocumenter(SqlAlchemyDocumenter)
