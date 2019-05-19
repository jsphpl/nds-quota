import logging
import os

script_location = os.path.dirname(os.path.realpath(__file__))
DEFAULT_FILENAME = os.path.join(script_location, 'log.txt')
DEFAULT_LEVEL = logging.INFO

def setupLogger(filename=DEFAULT_FILENAME, level=DEFAULT_LEVEL):
    logger = logging.getLogger()
    logger.setLevel(level)

    handler = logging.FileHandler(filename)
    handler.setFormatter(logging.Formatter('%(asctime)s [%(levelname)s] %(message)s'))
    logger.addHandler(handler)
