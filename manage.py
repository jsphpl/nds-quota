#!/usr/bin/env python3

from app import BaseApp
from commands.seed import Seed


class Manage(BaseApp):
    """Management script

    License: MIT
    Author: Joseph Paul <joseph@sehrgute.software>
    """

    def register_commands(self):
        self.add_command('seed', Seed)

    def register_arguments(self, parser):
        pass


if __name__ == '__main__':
    app = Manage()
    app.run()
