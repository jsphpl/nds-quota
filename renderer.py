import os
from jinja2 import Environment, FileSystemLoader

SCRIPT_LOCATION = os.path.dirname(os.path.realpath(__file__))
TEMPLATE_DIR = os.path.join(SCRIPT_LOCATION, 'templates')

class Renderer():
    def __init__(self):
        self.env = Environment(loader=FileSystemLoader(TEMPLATE_DIR))

    def render(self, name, context = {}):
        template = self.env.get_template(name + '.html.j2')
        return template.render(context)
