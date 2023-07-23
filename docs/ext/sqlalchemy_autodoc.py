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

    def __get_column_index(self, columns, target_column_name) -> int:  # type: ignore  # noqa: F821
        for index, column in enumerate(columns):
            if column.name == target_column_name:
                return index
        raise ValueError("index does not exist")

    def add_content(
        self,
        more_content: StringList | None,
        no_docstring: bool = False,
    ) -> None:
        sourcename = self.get_sourcename()
        table_name = f"{self.object.__name__} Columns"
        self.add_line(f"..  list-table:: {table_name}", sourcename)
        anch_table = f"{self.object.__name__}"
        self.add_line(f"   :name: {anch_table}", sourcename)
        self.add_line("", sourcename)
        self.add_line("   * - Column", sourcename)
        self.add_line("     - Type", sourcename)
        self.add_line("     - Key", sourcename)
        self.add_line("     - Unique", sourcename)

        mapper = class_mapper(self.object)
        for column in mapper.columns:
            anch = f"{self.object.__name__}-{column.name}"
            anch = anch.lower()
            self.add_line("   * - .. raw:: html", sourcename)
            self.add_line("", sourcename)
            self.add_line(
                f'         <a id="{anch}" href="#{anch}">{column.name}</a>', sourcename
            )
            self.add_line("     - {}".format(column.type), sourcename)
            if column.primary_key:
                self.add_line("     - Primary", sourcename)
            else:
                self.add_line("     - ", sourcename)
            if column.unique:
                self.add_line("     - True", sourcename)
            else:
                self.add_line("     - ", sourcename)


def setup(app: Sphinx):
    app.add_autodocumenter(SqlAlchemyDocumenter)
