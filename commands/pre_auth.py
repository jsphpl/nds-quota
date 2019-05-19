import logging
from cli_app import Command

Log = logging.getLogger(__name__)

class PreAuth(Command):
    """Handle pre auth"""

    @staticmethod
    def register_arguments(parser):
        parser.add_argument('query_string', type=str)

    def run(self):
        Log.info('Preauth with query string "%s"' % self.app.args.query_string)
