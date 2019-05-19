import sys
import logging
import random
import string
from cli_app import Command

Log = logging.getLogger(__name__)

DEFAULT_MAX_BYTES = 52428800  # 50MiB

class Seed(Command):
    """Seed the database"""

    @staticmethod
    def register_arguments(parser):
        parser.add_argument('users', type=int)
        parser.add_argument('--max-bytes', type=int, default=DEFAULT_MAX_BYTES)

    def run(self):
        Log.info('Seeding %d user(s)' % self.app.args.users)
        for i in range(self.app.args.users):
            user = self.app.models.User.create(
                username=self.random_string(8),
                password=self.random_string(8),
                max_bytes=self.app.args.max_bytes
            )
            note = 'Created user with username="%s" password="%s" max_bytes=%d' % (user.username, user.password, user.max_bytes)
            print(note)
            Log.info(note)

    def random_string(self, length=8):
        letters = string.ascii_lowercase
        return ''.join(random.choice(letters) for i in range(length))
