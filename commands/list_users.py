from cli_app import Command

class ListUsers(Command):
    """List all users in the database"""

    def run(self):
        for user in self.app.models.User.select():
            print('%s,%s,%d,%d' % (user.username, user.password, user.max_bytes, user.used_bytes))
