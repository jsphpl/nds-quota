from peewee import prefetch
import logging
from cli_app import Command

Log = logging.getLogger(__name__)

class DeauthExceeded(Command):
    """Deauth all users that have exceeded their quota"""

    def run(self):
        users = self.app.models.User.select().where(self.app.models.User.used_bytes >= self.app.models.User.max_bytes)
        devices = self.app.models.Device.select()
        users_with_devices = prefetch(users, devices)
        for user in users_with_devices:
            for device in user.devices:
                self.app.ndsctl.deauth(device.mac_address)
                device.delete_instance()
            Log.info('Deauthed user #%d with %d >= %d' % (user.id, user.used_bytes, user.max_bytes))

