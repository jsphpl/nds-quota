from cli_app import App
import models
from log import setupLogger
from ndsctl import NDSCTL

class BaseApp(App):
    def __init__(self):
        super(BaseApp, self).__init__()
        setupLogger()
        self.models = models
        self.models.connect()
        self.ndsctl = NDSCTL()
