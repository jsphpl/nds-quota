from cli_app import Command

class Catchall(Command):
    """Handle all other events"""

    @staticmethod
    def register_arguments(parser):
        parser.add_argument('mac_address', type=str)
        parser.add_argument('incoming_bytes', type=int)
        parser.add_argument('outgoing_bytes', type=int)
        parser.add_argument('session_start', type=int)
        parser.add_argument('session_end', type=int)

    def run(self):
        print(self.app.args)
