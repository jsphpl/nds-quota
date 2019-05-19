import logging

DEFAULT_FILENAME = './log.txt'
DEFAULT_LEVEL = logging.INFO

def setupLogger(filename=DEFAULT_FILENAME, level=DEFAULT_LEVEL):
    logger = logging.getLogger()
    logger.setLevel(level)

    handler = logging.FileHandler(filename)
    handler.setFormatter(logging.Formatter('%(asctime)s [%(levelname)s] %(message)s'))
    logger.addHandler(handler)
