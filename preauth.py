#!/usr/bin/env python3

import logging
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
        Log.info('Preauth with query string "%s"' % self.args.query_string)
        print(self.renderer.render('login'))

if __name__ == '__main__':
    app = PreAuth(Renderer())
    app.run()
