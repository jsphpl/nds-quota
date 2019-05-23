from peewee import prefetch
from cli_app import Command

class UpdateUsedBytes(Command):
    """Update used bytes for all users in the database"""

    def run(self):
        users = self.app.models.User.select()
        devices = self.app.models.Device.select()
        users_with_devices = prefetch(users, devices)
        for user in users_with_devices:
            used_bytes = 0
            for device in user.devices:
                client_data = self.app.ndsctl.json(device.mac_address)
                if client_data is not None and 'downloaded' in client_data:
                    used_bytes += client_data['downloaded']
            if used_bytes > 0:
                user.used_bytes = used_bytes
                user.save()
