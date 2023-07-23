# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information
import os
import sys

sys.path.insert(0, os.path.abspath("../"))  # add project root to system path
sys.path.insert(0, os.path.abspath("./ext"))

project = "Flamethrower"
copyright = "2023, caner"
author = "caner"
release = "0.0.1"

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration
html_theme = "furo"
html_static_path = ["_static"]
html_title = "Flamethrower API Documentation"
html_theme_options = {
    "light_logo": "logo.jpg",
    "dark_logo": "logo.jpg",
}

extensions = [
    "sphinx.ext.autodoc",
    "sphinx.ext.napoleon",
    "sqlalchemy_autodoc",
    "sphinx_tippy",
]

# Napoleon settings
napoleon_google_docstring = True
napoleon_numpy_docstring = True
napoleon_include_init_with_doc = False
napoleon_include_private_with_doc = False
napoleon_include_special_with_doc = False
napoleon_use_admonition_for_examples = False
napoleon_use_admonition_for_notes = False
napoleon_use_admonition_for_references = False
napoleon_use_ivar = False
napoleon_use_param = False
napoleon_use_rtype = False
set_type_checking_flag = True

templates_path = ["_templates"]
exclude_patterns = ["_build", "Thumbs.db", ".DS_Store"]
