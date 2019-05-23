from peewee import prefetch
from cli_app import Command

class UpdateQuota(Command):
    """Update used bytes for all users in the database"""

    def run(self):
        users = self.app.models.User.select()
        devices = self.app.models.Device.select()
        users_with_devices = prefetch(users, devices)
        for user in users_with_devices:
            used_kbytes = 0
            for device in user.devices:
                client_data = self.app.ndsctl.json(device.mac_address)
                if client_data is not None and 'downloaded' in client_data:
                    used_kbytes += client_data['downloaded']
            if used_kbytes > 0:
                user.used_bytes = used_kbytes * 1024
                user.save()
