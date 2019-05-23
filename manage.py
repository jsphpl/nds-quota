#!/usr/bin/env python3

from app import BaseApp
from commands.seed import Seed
from commands.list_users import ListUsers
from commands.update_used_bytes import UpdateUsedBytes


class Manage(BaseApp):
    """Management script

    License: MIT
    Author: Joseph Paul <joseph@sehrgute.software>
    """

    def register_commands(self):
        self.add_command('seed', Seed)
        self.add_command('list-users', ListUsers)
        self.add_command('update-used-bytes', UpdateUsedBytes)

    def register_arguments(self, parser):
        pass


if __name__ == '__main__':
    app = Manage()
    app.run()
