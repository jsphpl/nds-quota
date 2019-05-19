#!/usr/bin/env python3

from app import BaseApp
from commands.auth_client import AuthClient
from commands.catchall import Catchall
from commands.pre_auth import PreAuth


class BinAuth(BaseApp):
    """BinAuth script for implementing per-user quota in Nodogsplash

    License: MIT
    Author: Joseph Paul <joseph@sehrgute.software>
    """

    def register_commands(self):
        self.add_command('pre_auth', PreAuth)
        self.add_command('auth_client', AuthClient)
        self.add_command('client_auth', Catchall)
        self.add_command('client_deauth', Catchall)
        self.add_command('idle_deauth', Catchall)
        self.add_command('timeout_deauth', Catchall)
        self.add_command('ndsctl_auth', Catchall)
        self.add_command('ndsctl_deauth', Catchall)
        self.add_command('shutdown_deauth', Catchall)

    def register_arguments(self, parser):
        pass


if __name__ == '__main__':
    app = BinAuth()
    app.run()
