from distutils.core import setup
from distutils.extension import Extension
from Cython.Distutils import build_ext

ext_modules = [
    Extension('app', ['app.py']),
    Extension('commands.auth_client', ['commands/auth_client.py']),
    Extension('commands.catchall', ['commands/catchall.py']),
    Extension('commands.seed', ['commands/seed.py']),
    Extension('log', ['log.py']),
    Extension('models', ['models.py']),
]

setup(
    name = 'manage.py',
    cmdclass = {'build_ext': build_ext},
    ext_modules = ext_modules
)
