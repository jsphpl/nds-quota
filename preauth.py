#!/usr/bin/env python3

import logging
import urllib.parse
from app import BaseApp
from renderer import Renderer

Log = logging.getLogger(__name__)

class PreAuth(BaseApp):
    def __init__(self, renderer):
        super(PreAuth, self).__init__()
        self.renderer = renderer

    def register_arguments(self, parser):
        parser.add_argument('query_string', type=str)

    def run(self):
        query = self.parse_query()
        Log.info('Preauth with query string "%s"' % query)
        print(self.renderer.render('login', { 'query': query }))

    def parse_query(self):
        try:
            unquoted = urllib.parse.unquote(self.args.query_string[3:])
            splitted = unquoted.split(', ')
            pairs = [i.split('=') for i in splitted]
            return dict(pairs)
        except Exception:
            return {}

if __name__ == '__main__':
    app = PreAuth(Renderer())
    app.run()
