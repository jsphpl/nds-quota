import subprocess
import json

class NDSCTL():
    def json(self, identifier = ''):
        try:
            return json.loads(self.execute('json', identifier))
        except Exception:
            return None

    def deauth(self, identifier):
        try:
            self.execute('deauth', identifier)
            return True
        except Exception:
            return False

    def execute(self, *args):
        command = subprocess.Popen(
            ['ndsctl', *args],
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT
        )
        stdout, stderr = command.communicate()
        return stdout
