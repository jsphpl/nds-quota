from cli_app import App
import models
from log import setupLogger

class BaseApp(App):
    def __init__(self):
        super(BaseApp, self).__init__()
        setupLogger()
        self.models = models
        self.models.connect()
