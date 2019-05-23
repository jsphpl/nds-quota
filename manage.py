#!/usr/bin/env python3

from app import BaseApp
from commands.seed import Seed
from commands.list_users import ListUsers
from commands.update_quota import UpdateQuota
from commands.deauth_exceeded import DeauthExceeded


class Manage(BaseApp):
    """Management script

    License: MIT
    Author: Joseph Paul <joseph@sehrgute.software>
    """

    def register_commands(self):
        self.add_command('seed', Seed)
        self.add_command('list-users', ListUsers)
        self.add_command('update-quota', UpdateQuota)
        self.add_command('deauth-exceeded', DeauthExceeded)

    def register_arguments(self, parser):
        pass


if __name__ == '__main__':
    app = Manage()
    app.run()
