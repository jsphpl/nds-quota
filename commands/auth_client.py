import logging
import peewee
from cli_app import Command

Log = logging.getLogger(__name__)

SESSION_DURATION = 86400  # 24h in seconds

class AuthClient(Command):
    """Handle auth_client event"""

    @staticmethod
    def register_arguments(parser):
        parser.add_argument('mac_address', type=str)
        parser.add_argument('username', type=str, help='Username passed by NDS')
        parser.add_argument('password', type=str, help='Password passed by NDS')

    def run(self):
        username = self.app.args.username
        password = self.app.args.password
        try:
            user = self.app.models.User.get(username=username)
        except Exception:
            Log.info('Auth attempt with username="%s" failed. User not found.' % username)
            exit(1)

        if (user.password != password):
            Log.info('Auth attempt with username="%s" failed. Wrong password.' % username)
            exit(1)

        self.assign_device_to_user(user)

        Log.info('Auth attempt with username="%s" succeeded.' % username)
        print('%d 0 %d' % (SESSION_DURATION, user.max_bytes))
        exit(0)

